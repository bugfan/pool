package main

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/abiosoft/ishell"
	"themis/client/clientconfig"
)

func main() {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("\nWebVPN Tunnel Shell")

	// register a function for "greet" command.
	registerIPCofnig := func(name, conf string) {
		shell.Register(name, func(args ...string) (string, error) {
			if len(args) != 1 {
				return fmt.Sprintf("Usage: %s xx.xx.xx.xx", name), errors.New("Argument count not match")
			}
			clientconfig.Set(fmt.Sprintf(conf, name), args[0])
			return "OK", nil
		})
	}
	registerIPCofnig("ip", "network.ip")
	registerIPCofnig("netmask", "network.netmask")
	registerIPCofnig("gateway", "network.gateway")
	registerIPCofnig("dns", "network.dns")
	registerIPCofnig("server", "server.server")

	shell.Register("show_conf", func(args ...string) (string, error) {
		return clientconfig.String(), nil
	})

	shell.Register("client_id", func(args ...string) (string, error) {
		clientID, _ := clientconfig.Get("server.client_id")
		return clientID, nil
	})

	shell.Register("reboot", func(args ...string) (string, error) {
		exec.Command("reboot").Start()
		return "OK", nil
	})

	shell.Register("save", func(args ...string) (string, error) {
		clientconfig.Save()
		setNetwork()
		return "OK", nil
	})
	// start shell
	shell.Start()
}
