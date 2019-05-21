// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log

import (
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/alecthomas/log4go"
)

const (
	LOCAL0 = 16
	LOCAL1 = 17
	LOCAL2 = 18
	LOCAL3 = 19
	LOCAL4 = 20
	LOCAL5 = 21
	LOCAL6 = 22
	LOCAL7 = 23
)

// This log writer sends output to a socket
type SysLogWriter chan *log.LogRecord

// This is the SocketLogWriter's output method
func (w SysLogWriter) LogWrite(rec *log.LogRecord) {
	w <- rec
}

func (w SysLogWriter) Close() {
	close(w)
}

func connectSyslogDaemon() (sock net.Conn, err error) {
	network := "unix"
	raddr := "/dev/log"
	sock, err = net.Dial(network, raddr)
	if err != nil {
		err = fmt.Errorf("cannot connect to Syslog Daemon: %s", err)
	} else {
		fmt.Fprintf(os.Stderr, "syslog uses %s:%s\n", network, raddr)
		return
	}
	return
}

func NewSysLogWriter(facility int) (w SysLogWriter) {
	offset := facility * 8
	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot obtain hostname: %s\n", err.Error())
		host = "unknown"
	}
	sock, err := connectSyslogDaemon()
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSysLogWriter: %s\n", err.Error())
		return
	}
	w = SysLogWriter(make(chan *log.LogRecord, log.LogBufferLength))
	go func() {
		defer func() {
			if sock != nil {
				sock.Close()
			}
		}()
		var timestr string
		var timestrAt int64
		for rec := range w {
			if rec.Created.Unix() != timestrAt {
				timestrAt = rec.Created.Unix()
				timestr = time.Unix(timestrAt, 0).Local().Format(time.RFC3339)
			}
			fmt.Fprintf(sock, "<%d>%s %s %s: %s\n", offset+int(rec.Level), timestr, host, "WebVPN-Daemon", rec.Message)
		}
	}()
	return
}
