package conn

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"pool/log"
)

type Conn interface {
	net.Conn
	log.Logger
	ID() string
	SetType(string)
	CloseRead() error
}

type loggedConn struct {
	tcp *net.TCPConn
	net.Conn
	log.Logger
	id  int32
	typ string
}

func Wrap(conn net.Conn, typ string) *loggedConn {
	tcp := conn.(*net.TCPConn)
	wrapped := &loggedConn{tcp,
		conn,
		log.NewPrefixLogger(),
		rand.Int31(), typ}
	wrapped.AddLogPrefix(wrapped.ID())
	return wrapped
}

func DialTimeout(addr, typ string, tlsCfg *tls.Config, timeout time.Duration) (conn *loggedConn, err error) {
	var rawConn net.Conn
	if rawConn, err = net.DialTimeout("tcp", addr, timeout); err != nil {
		return
	}

	conn = Wrap(rawConn, typ)
	conn.Debug("New connection to: %v", rawConn.RemoteAddr())

	if tlsCfg != nil {
		conn.StartTLS(tlsCfg)
	}

	return
}

func Dial(addr, typ string, tlsCfg *tls.Config) (conn *loggedConn, err error) {
	var rawConn net.Conn
	if rawConn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	conn = Wrap(rawConn, typ)
	conn.Debug("New connection to: %v", rawConn.RemoteAddr())

	if tlsCfg != nil {
		conn.StartTLS(tlsCfg)
	}

	return
}

func (c *loggedConn) StartTLS(tlsCfg *tls.Config) {
	c.Conn = tls.Client(c.Conn, tlsCfg)
}

func (c *loggedConn) Close() (err error) {
	if err := c.Conn.Close(); err == nil {
		c.Debug("Closing")
	}
	return
}

func (c *loggedConn) ID() string {
	return fmt.Sprintf("%s:%x", c.typ, c.id)
}

func (c *loggedConn) SetType(typ string) {
	oldID := c.ID()
	c.typ = typ
	c.ClearLogPrefixes()
	c.AddLogPrefix(c.ID())
	c.Debug("Renamed connection %s", oldID)
}

func (c *loggedConn) CloseRead() error {
	// XXX: use CloseRead() in Conn.Join() and in Control.shutdown() for cleaner
	// connection termination. Unfortunately, when I've tried that, I've observed
	// failures where the connection was closed *before* flushing its write buffer,
	// set with SetLinger() set properly (which it is by default).
	return c.tcp.CloseRead()
}

// func Join(c Conn, c2 Conn) (int64, int64) {
// 	var wait sync.WaitGroup

// 	pipe := func(to Conn, from Conn, bytesCopied *int64) {
// 		defer to.Close()
// 		defer wait.Done()

// 		var err error
// 		*bytesCopied, err = io.Copy(to, from)
// 		if err != nil {
// 			from.Warn("Copied %d bytes to %s before failing with error %v", *bytesCopied, to.ID(), err)
// 		} else {
// 			from.Debug("Copied %d bytes to %s", *bytesCopied, to.ID())
// 		}
// 	}

// 	wait.Add(2)
// 	var fromBytes, toBytes int64
// 	go pipe(c, c2, &fromBytes)
// 	go pipe(c2, c, &toBytes)
// 	c.Info("Joined with connection %s", c2.ID())
// 	wait.Wait()
// 	return fromBytes, toBytes
// }

func Join(conn1 Conn, conn2 Conn) (conn1Write int64, conn2Write int64) {
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)

	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer conn2.Close()
		conn2Write = chanToConn(conn2, chan1)
		wait.Done()
	}()

	go func() {
		defer conn1.Close()
		conn1Write = chanToConn(conn1, chan2)
		wait.Done()
	}()
	wait.Wait()
	return

}

func chanToConn(conn Conn, c chan []byte) (wr int64) {
	for b := range c {
		nw, err := conn.Write(b)
		wr += int64(nw)
		if err != nil {
			return
		}
	}
	return
}

// chanFromConn a channel from a Conn object, and sends everything it
//  Read()s from the socket to the channel.
func chanFromConn(conn Conn) chan []byte {
	c := make(chan []byte, 128)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				return
			}

		}()
		defer close(c)

		b := make([]byte, 64*1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				break
			}
		}
	}()

	return c
}
