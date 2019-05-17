package client

import (
	"crypto"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"pool/tunnel/conn"
	"pool/tunnel/msg"
	"pool/tunnel/util"
	"pool/tunnel/version"
	"themis/log"

	update "github.com/inconshreveable/go-update"
)

const (
	pingInterval   = 5 * time.Second
	maxPongLatency = 15 * time.Second
)

type Client struct {
	id string
	log.Logger
	serverVersion string
	authToken     string
	serverAddr    string
	tlsConfig     *tls.Config
	connStatus    ConnStatus
	ctlConn       conn.Conn
	lastPong      int64
}

func NewClient(serverAddr string, id string, tls *tls.Config) *Client {
	c := Client{
		id:         id,
		Logger:     log.NewPrefixLogger("client"),
		serverAddr: serverAddr,
		tlsConfig:  tls,
	}
	return &c
}
func (c *Client) GetID() string {
	return c.id
}

func (c *Client) Run() {
	// how long we should wait before we reconnect
	maxWait := 30 * time.Second
	wait := 1 * time.Second

	for {
		// run the control channel
		c.control()

		// control only returns when a failure has occurred, so we're going to try to reconnect
		if c.connStatus == ConnOnline {
			wait = 1 * time.Second
		}

		log.Info("Waiting %d seconds before reconnecting", int(wait.Seconds()))
		time.Sleep(wait)
		// exponentially increase wait time
		wait = 2 * wait
		wait = time.Duration(math.Min(float64(wait), float64(maxWait)))
		c.connStatus = ConnReconnecting
	}
}

func (c *Client) control() {
	defer func() {
		if r := recover(); r != nil {
			c.Error("control recovering from failure %v", r)
		}
	}()

	// establish control channel
	if err := c.controlConnect(); err != nil {
		panic(err)
	}
	defer c.ctlConn.Close()

	if err := c.auth(); err != nil {
		panic(err)
	}

	c.startHeartbeat()

	c.mainLoop()
}

func (c *Client) controlConnect() (err error) {
	c.ctlConn, err = conn.Dial(c.serverAddr, "ctl", c.tlsConfig)
	return
}

func (c *Client) update(resp *msg.AuthResp) {
	fmt.Printf("client update to %s\n", resp.MmVersion)
	client := new(http.Client)
	sha256Bytes := make([]byte, 64)
	sha256Resp, err := client.Do(sha256File(resp))
	if err != nil {
		panic(err)
	}
	defer sha256Resp.Body.Close()
	_, err = sha256Resp.Body.Read(sha256Bytes)
	if err != nil && err != io.EOF {
		panic(err)
	}
	checksum, err := hex.DecodeString(strings.TrimSpace(string(sha256Bytes)))
	if err != nil && err != io.EOF {
		panic(err)
	}

	client = new(http.Client)
	file, err := client.Do(updateFile(resp))
	if err != nil {
		panic(err)
	}
	defer file.Body.Close()
	err = update.Apply(file.Body, update.Options{
		Hash:     crypto.SHA256, // this is the default, you don't need to specify it
		Checksum: checksum,
	})
	if err != nil && err != io.EOF {
		if rerr := update.RollbackError(err); rerr != nil {
			panic(rerr)
		}
		panic(err)
	}
	os.Exit(0)
}

func updateFile(resp *msg.AuthResp) *http.Request {
	url := fmt.Sprintf("https://update.webvpn.net.cn/client/client-%s-%s", resp.Version, resp.MmVersion)
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func sha256File(resp *msg.AuthResp) *http.Request {
	url := fmt.Sprintf("https://update.webvpn.net.cn/client/client-%s-%s.sha256", resp.Version, resp.MmVersion)
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func (c *Client) auth() error {
	// authenticate with the server
	auth := &msg.Auth{
		ClientId:  c.id,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   version.Proto,
		MmVersion: version.MajorMinor(),
		User:      c.authToken,
	}

	if err := msg.WriteMsg(c.ctlConn, auth); err != nil {
		return err
	}

	var authResp msg.AuthResp
	if err := msg.ReadMsgInto(c.ctlConn, &authResp); err != nil {
		return err
	}

	if authResp.Error != "" {
		if authResp.Version != "" && authResp.Version != version.Proto {
			c.update(&authResp)
		}
		return fmt.Errorf("Failed to authenticate to server: %s",
			authResp.Error)
	}

	c.id = authResp.ClientId
	c.serverVersion = authResp.MmVersion
	c.Info("Authenticated with server, client id: %v", c.id)
	return nil
}

func (c *Client) Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := util.MakePanicTrace(r)
				c.Error(err)
			}
		}()

		fn()
	}()
}

func (c *Client) startHeartbeat() {
	c.lastPong = time.Now().UnixNano()
	c.Go(c.heartbeat)
}

// Hearbeating to ensure our connection still live
func (c *Client) heartbeat() {
	lastPongAddr := &c.lastPong
	lastPing := time.Unix(atomic.LoadInt64(lastPongAddr)-1, 0)
	ping := time.NewTicker(pingInterval)
	pongCheck := time.NewTicker(time.Second)

	defer func() {
		c.ctlConn.Close()
		ping.Stop()
		pongCheck.Stop()
	}()

	for {
		select {
		case <-pongCheck.C:
			lastPong := time.Unix(0, atomic.LoadInt64(lastPongAddr))
			needPong := lastPong.Sub(lastPing) < 0
			pongLatency := time.Since(lastPing)

			if needPong && pongLatency > maxPongLatency {
				c.Info("Last ping: %v, Last pong: %v", lastPing, lastPong)
				c.Info("Connection stale, haven't gotten PongMsg in %d seconds", int(pongLatency.Seconds()))
				return
			}

		case <-ping.C:
			err := msg.WriteMsg(c.ctlConn, &msg.Ping{})
			if err != nil {
				c.ctlConn.Debug("Got error %v when writing PingMsg", err)
				return
			}
			lastPing = time.Now()
		}
	}
}

func (c *Client) mainLoop() {
	for {
		rawMsg, err := msg.ReadMsg(c.ctlConn)
		if err != nil {
			panic(err)
		}

		c.ctlConn.Info("RawMsg: %#v", rawMsg)

		switch m := rawMsg.(type) {
		case *msg.ReqProxy:
			c.Go(c.proxy)

		case *msg.Pong:
			atomic.StoreInt64(&c.lastPong, time.Now().UnixNano())

		default:
			c.ctlConn.Warn("Ignoring unknown control message %v ", m)
		}
	}
}

func (c *Client) proxy() {
	var (
		remoteConn conn.Conn
		err        error
	)
	remoteConn, err = conn.Dial(c.serverAddr, "pxy", c.tlsConfig)
	if err != nil {
		c.Error("Failed to establish proxy connection: %v", err)
		return
	}
	defer remoteConn.Close()

	err = msg.WriteMsg(remoteConn, &msg.RegProxy{ClientId: c.id})
	if err != nil {
		remoteConn.Error("Failed to write RegProxy: %v", err)
		return
	}

	var startPxy msg.StartProxy
	if err = msg.ReadMsgInto(remoteConn, &startPxy); err != nil {
		remoteConn.Error("Client failed to read StartProxy: %v", err)
		return
	}

	localConn, err := conn.DialTimeout(startPxy.LocalAddr, "prv", nil, 60*time.Second)
	if err != nil {
		repl := &msg.StartProxyRepl{
			Status: "error",
			Err:    err.Error(),
		}
		msg.WriteMsg(remoteConn, repl)

		remoteConn.Warn("Failed to open private leg %s: %v", startPxy.LocalAddr, err)
		return
	}
	defer localConn.Close()
	repl := &msg.StartProxyRepl{
		Status: "ok",
	}
	err = msg.WriteMsg(remoteConn, repl)
	if err != nil {
		return
	}
	conn.Join(localConn, remoteConn)
}
