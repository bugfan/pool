package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"pool/tunnel/conn"
	"pool/tunnel/msg"
	"pool/tunnel/util"
	"pool/tunnel/version"
	"themis/log"
)

const (
	pingTimeoutInterval    = 30 * time.Second
	connReapInterval       = 10 * time.Second
	controlWriteTimeout    = 10 * time.Second
	proxyStaleDuration     = 60 * time.Second
	resendReqProxyInterval = 1 * time.Second
	proxyMaxPoolSize       = 1024
	proxyMaxRetry          = 20
)

type Control struct {
	// logger
	log.Logger
	// auth message
	auth *msg.Auth

	// actual connection
	conn conn.Conn

	// put a message in this channel to send it over
	// conn to the client
	out chan (msg.Message)

	// read from this channel to get the next message sent
	// to us over conn by the client
	in chan (msg.Message)

	// the last time we received a ping from the client - for heartbeats
	lastPing time.Time

	// proxy connections
	proxies chan conn.Conn

	// proxy request chain
	proxyRequests chan chan<- conn.Conn

	// identifier
	id string

	// synchronizer for controlled shutdown of writer()
	writerShutdown *util.Shutdown

	// synchronizer for controlled shutdown of reader()
	readerShutdown *util.Shutdown

	// synchronizer for controlled shutdown of manager()
	managerShutdown *util.Shutdown

	// synchronizer for controller shutdown of entire Control
	shutdown *util.Shutdown

	averageConn int
}

func NewControl(ctlConn conn.Conn, authMsg *msg.Auth) {
	var err error

	// create the object
	c := &Control{
		Logger:          log.NewPrefixLogger("Control"),
		auth:            authMsg,
		conn:            ctlConn,
		out:             make(chan msg.Message, proxyMaxPoolSize),
		in:              make(chan msg.Message, proxyMaxPoolSize),
		proxies:         make(chan conn.Conn, proxyMaxPoolSize),
		proxyRequests:   make(chan chan<- conn.Conn, proxyMaxPoolSize),
		lastPing:        time.Now(),
		averageConn:     1,
		writerShutdown:  util.NewShutdown(),
		readerShutdown:  util.NewShutdown(),
		managerShutdown: util.NewShutdown(),
		shutdown:        util.NewShutdown(),
	}

	failAuth := func(e error) {
		_ = msg.WriteMsg(ctlConn, &msg.AuthResp{
			Version:   version.Proto,
			MmVersion: version.MajorMinor(),
			Error:     e.Error()},
		)
		ctlConn.Close()
	}

	// register the clientid
	c.id = authMsg.ClientId
	if c.id == "" {
		// it's a new session, assign an ID
		if c.id, err = util.SecureRandId(16); err != nil {
			failAuth(err)
			return
		}
	}

	// set logging prefix
	ctlConn.SetType("ctl")
	ctlConn.AddLogPrefix(c.id)

	if authMsg.Version != version.Proto {
		failAuth(fmt.Errorf("Incompatible versions. Server %s, client %s.", version.MajorMinor(), authMsg.Version))
		return
	}

	// register the control
	if replaced := controlRegistry.Add(c.id, c); replaced != nil {
		replaced.shutdown.WaitComplete()
	}

	// start the writer first so that the following messages get sent
	go c.writer()

	// Respond to authentication
	c.out <- &msg.AuthResp{
		Version:   version.Proto,
		MmVersion: version.MajorMinor(),
		ClientId:  c.id,
	}

	// As a performance optimization, ask for a proxy connection up front
	c.out <- &msg.ReqProxy{}

	// manage the connection
	go c.manager()
	go c.reader()
	go c.stopper()
	go c.connManager()
}

func (c *Control) manager() {
	// don't crash on panics
	defer func() {
		if err := recover(); err != nil {
			c.conn.Info("Control::manager failed with error %v: %s", err, debug.Stack())
		}
	}()

	// kill everything if the control manager stops
	defer c.shutdown.Begin()

	// notify that manager() has shutdown
	defer c.managerShutdown.Complete()

	// reaping timer for detecting heartbeat failure
	reap := time.NewTicker(connReapInterval)
	defer reap.Stop()

	for {
		select {
		case <-reap.C:
			if time.Since(c.lastPing) > pingTimeoutInterval {
				c.conn.Info("Lost heartbeat")
				c.shutdown.Begin()
			}

		case mRaw, ok := <-c.in:
			// c.in closes to indicate shutdown
			if !ok {
				return
			}

			switch mRaw.(type) {
			case *msg.Ping:
				c.lastPing = time.Now()
				c.out <- &msg.Pong{}
			}
		}
	}
}

func (c *Control) writer() {
	defer func() {
		if err := recover(); err != nil {
			c.conn.Info("Control::writer failed with error %v: %s", err, debug.Stack())
		}
	}()

	// kill everything if the writer() stops
	defer c.shutdown.Begin()

	// notify that we've flushed all messages
	defer c.writerShutdown.Complete()

	// write messages to the control channel
	for m := range c.out {
		c.conn.SetWriteDeadline(time.Now().Add(controlWriteTimeout))
		if err := msg.WriteMsg(c.conn, m); err != nil {
			panic(err)
		}
	}
}

func (c *Control) reader() {
	defer func() {
		if err := recover(); err != nil {
			c.conn.Warn("Control::reader failed with error %v: %s", err, debug.Stack())
		}
	}()

	// kill everything if the reader stops
	defer c.shutdown.Begin()

	// notify that we're done
	defer c.readerShutdown.Complete()

	// read messages from the control channel
	for {
		if msg, err := msg.ReadMsg(c.conn); err != nil {
			if err == io.EOF {
				c.conn.Info("EOF")
				return
			} else {
				panic(err)
			}
		} else {
			// this can also panic during shutdown
			c.in <- msg
		}
	}
}

func (c *Control) stopper() {
	defer func() {
		if r := recover(); r != nil {
			c.conn.Error("Failed to shut down control: %v", r)
		}
	}()

	// wait until we're instructed to shutdown
	c.shutdown.WaitBegin()

	// remove ourself from the control registry
	controlRegistry.Del(c.id)

	// shutdown manager() so that we have no more work to do
	close(c.in)
	c.managerShutdown.WaitComplete()

	// shutdown writer()
	close(c.out)
	c.writerShutdown.WaitComplete()

	// close connection fully
	c.conn.Close()

	// shutdown all of the proxy connections
	close(c.proxies)
	for p := range c.proxies {
		p.Close()
	}

	close(c.proxyRequests)

	c.shutdown.Complete()
	c.conn.Info("Shutdown complete")
}

func (c *Control) RegisterProxy(conn conn.Conn) {
	conn.AddLogPrefix(c.id)

	conn.SetDeadline(time.Now().Add(proxyStaleDuration))
	select {
	case c.proxies <- conn:
		conn.Debug("Registered")
	default:
		conn.Warn("Proxies buffer is full, discarding.")
		conn.Close()
	}
}

func (c *Control) connManager() {
	defer c.shutdown.Begin()
	var counter uint64 = 0

	clean := func() {
		for len(c.proxies) > c.averageConn*2 {
			select {
			case p, ok := <-c.proxies:
				if !ok {
					return
				} else {
					p.Close()
				}
			default:
				return
			}
		}
	}

	timeout := time.After(1 * time.Minute)

	for {

		select {
		case req, ok := <-c.proxyRequests:
			if !ok {
				return
			} else {
				atomic.AddUint64(&counter, 1)
				c.getConnFor(req)
			}
		case <-timeout:
			count := atomic.LoadUint64(&counter)
			c.averageConn = int(count)
			atomic.StoreUint64(&counter, 0)
			timeout = time.After(1 * time.Minute)
			go clean()
		}
	}

}

func (c *Control) getConnFor(req chan<- conn.Conn) {

	// get a proxy connection from the pool
	select {
	case proxyConn, ok := <-c.proxies:
		if !ok {
			close(req)
			return
		}
		req <- proxyConn
		if len(c.proxies) < c.averageConn {
			go c.requestProxy()
		}

	default:
		// no proxy available in the pool, ask for one over the control channel
		c.conn.Debug("No proxy in pool, requesting proxy from control . . .")
		if err := c.requestProxy(); err != nil {
			close(req)
			return
		}
		// request new proxy for next conn if proxie pool is empty
		defer c.requestProxy()

		resendReq := time.After(resendReqProxyInterval)
		for {
			select {
			case proxyConn, ok := <-c.proxies:
				if !ok {
					close(req)
					return
				}
				req <- proxyConn
				return
			case <-resendReq:
				if err := c.requestProxy(); err != nil {
					close(req)
					return
				}
				resendReq = time.After(resendReqProxyInterval)
			case <-time.After(pingTimeoutInterval):
				close(req)
				return
			}

		}
	}
	return
}

// Remove a proxy connection from the pool and return it
// If not proxy connections are in the pool, request one
// and wait until it is available
// Returns an error if we couldn't get a proxy because it took too long
// or the tunnel is closing
func (c *Control) GetProxy() (proxyConn conn.Conn, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("tunnel closed")
		}
	}()
	req := make(chan conn.Conn, 1)
	c.proxyRequests <- req
	conn, ok := <-req
	if ok {
		return conn, nil
	}
	return conn, errors.New("errror when get proxy")
}

func (c *Control) requestProxy() error {
	return util.PanicToError(func() { c.out <- &msg.ReqProxy{} })

}

// Called when this control is replaced by another control
// this can happen if the network drops out and the client reconnects
// before the old tunnel has lost its heartbeat
func (c *Control) Replaced(replacement *Control) {
	c.conn.Info("Replaced by control: %s", replacement.conn.ID())

	// set the control id to empty string so that when stopper()
	// calls registry.Del it won't delete the replacement
	c.id = ""

	// tell the old one to shutdown
	c.shutdown.Begin()
}

func (c *Control) Dial(network, address string) (net.Conn, error) {
	var proxyConn conn.Conn
	var err error

	for i := 0; i < proxyMaxRetry; i++ {
		// get a proxy connection
		if proxyConn, err = c.GetProxy(); err != nil {
			c.Warn("Failed to get proxy connection: %v", err)
			return nil, err
		}
		c.Debug("Got proxy connection %s", proxyConn.ID())
		proxyConn.AddLogPrefix(proxyConn.ID())

		// tell the client we're going to start using this proxy connection
		startPxyMsg := &msg.StartProxy{
			LocalAddr: address,
		}

		if err = msg.WriteMsg(proxyConn, startPxyMsg); err != nil {
			proxyConn.Warn("Failed to write StartProxyMessage: %v, attempt %d", err, i)
			proxyConn.Close()
		} else {
			proxyConn.SetDeadline(time.Time{})
			// success
			var repl msg.StartProxyRepl
			err = msg.ReadMsgInto(proxyConn, &repl)
			if err == nil {
				if repl.Status == "ok" {
					break
				} else {
					err = errors.New(repl.Err)
					proxyConn.Close()
				}
			} else {
				proxyConn.Close()
			}
		}
	}

	if err != nil {
		// give up
		c.Error("Too many failures starting proxy connection")
		return nil, err
	}
	return proxyConn, nil
}

func GetControl(name string) *Control {
	return controlRegistry.Get(name)
}

func GetHTTPTransport(tunnel, ip, client_certficate string) (*http.Transport, error) {
	var transport *http.Transport
	var cert []tls.Certificate

	if client_certficate != "" {
		bytes := []byte(client_certficate)
		keypair, err := tls.X509KeyPair(bytes, bytes)
		if err != nil {
			return nil, err
		}
		cert = []tls.Certificate{keypair}
	}
	tlsCfg := &tls.Config{
		Certificates:       cert,
		InsecureSkipVerify: true,
	}

	if tunnel != "" {
		control := GetControl(tunnel)
		if control == nil {
			return nil, fmt.Errorf("No such tunnel %s", tunnel)
		}
		if ip == "" {
			transport = &http.Transport{Dial: control.Dial, TLSClientConfig: tlsCfg}
		} else {
			dialer := func(network, address string) (net.Conn, error) {
				port := strings.Split(address, ":")[1]
				newAddress := fmt.Sprintf("%s:%s", ip, port)
				return control.Dial(network, newAddress)
			}
			transport = &http.Transport{Dial: dialer, TLSClientConfig: tlsCfg}
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: tlsCfg,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
		}
	}
	transport.DisableKeepAlives = true
	transport.DisableCompression = true
	transport.ResponseHeaderTimeout = 1 * time.Minute

	// transport.Proxy = http.ProxyFromEnvironment
	transport.MaxIdleConns = 100
	transport.IdleConnTimeout = 90 * time.Second
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ExpectContinueTimeout = 1 * time.Second
	return transport, nil
}
