package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"pool/client/clientconfig"
	"pool/tunnel/client"
	"pool/tunnel/version"

	_ "net/http/pprof"
)

func main() {

	if os.Getenv("GO_PROF") == "yes" {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}
	opts, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.version {
		fmt.Print(version.Full())
		os.Exit(0)
	}

	data, err := Asset("cert/tunnel.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tlsConfig := client.LoadTLSConfigFromBytes(data)
	if tlsConfig != nil {
		tlsConfig.ClientSessionCache = tls.NewLRUClientSessionCache(5000)
	}
	c := client.NewClient(opts.serverAddr, opts.clientID, tlsConfig)
	go func() {
		for {
			id := c.GetID()
			if id != "" {
				clientconfig.Set("server.client_id", id)
				clientconfig.Save()
				return
			} else {
				runtime.Gosched()
			}
		}
	}()
	c.Run()
}
