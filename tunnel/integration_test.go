package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"pool/log"
	"pool/tunnel/client"
	"pool/tunnel/server"
)

const (
	controlName = "idxxxxx"
	maxWait     = 10
	echoAddr    = "127.0.0.1:3333"
)

func TestIntegration(t *testing.T) {
	var control *server.Control
	for count := 0; count < maxWait; count++ {
		control = server.GetControl(controlName)
		if control != nil {
			break
		}
		time.Sleep(time.Second)
	}

	if control == nil {
		t.Fatalf("Can't get control retryed %d times", maxWait)
	}

	conn, err := control.Dial("tcp", echoAddr)
	if err != nil {
		t.Fatal("connect echo server fail")
	}

	writeTest(conn, "hello", t)
}

func writeTest(c net.Conn, s string, t *testing.T) {
	msg := []byte(s)
	c.Write(msg)
	buf := make([]byte, 1024)
	length, err := c.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("read %d bytes from proxy before fail: %s", length, err)
	}

	buf = buf[:length]
	if string(msg) != string(buf) {
		t.Fatalf("read msg unmatch expect %s, actual %s", msg, buf)
	}
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		teardown()
	}
	os.Exit(ret)
}

func setup() {
	log.LogTo("stdout", "")
	go func() {
		server.TunnelListener("127.0.0.1:4433", server.LoadTLSConfig("../cert"))
	}()
	go func() {
		client.NewClient("127.0.0.1:4433", controlName, client.LoadTLSConfig("../cert")).Run()
	}()

	go echoServer()
}

func teardown() {

}

func echoServer() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", echoAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + echoAddr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	length, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Write the message in the connection channel.
	conn.Write(buf[:length])
	// Close the connection when you're done with it.
}
