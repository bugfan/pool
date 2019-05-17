package main

import (
	"flag"
	"fmt"
	"os"

	"pool/client/clientconfig"
)

const usage string = "Usage: %s [OPTIONS]\n"

type Options struct {
	certPath   string
	clientID   string
	serverAddr string
	logto      string
	loglevel   string
	version    bool
}

func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	logto := flag.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings")

	loglevel := flag.String(
		"level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")

	version := flag.Bool(
		"v",
		false,
		"version")

	flag.Parse()

	serverAddr, _ := clientconfig.Get("server.server")
	clientID, _ := clientconfig.Get("server.client_id")
	return &Options{
		certPath:   "cert",
		logto:      *logto,
		loglevel:   *loglevel,
		serverAddr: serverAddr,
		clientID:   clientID,
		version:    *version,
	}, nil
}
