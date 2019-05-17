package main

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"pool/client/clientconfig"
)

var networkTemp = `
auto bond0
iface bond0 inet static
    address %s
    netmask %s
    gateway %s
    slaves eth0 eth1 eth2 eth3 eth4 eth5
    bond_mode active-backup
    bond_miimon 100
    bond_downdelay 200
    bond_updelay 200

iface bond0 inet6 dhcp
dns-nameservers %s 223.5.5.5
`

var netConfFile = "/etc/network/interfaces.d/bond0"

func init() {
	if runtime.GOOS == "darwin" {
		netConfFile = "bond0"
	}
}

func setNetwork() {
	ip, _ := clientconfig.Get("network.ip")
	netmask, _ := clientconfig.Get("network.netmask")
	gateway, _ := clientconfig.Get("network.gateway")
	dns, _ := clientconfig.Get("network.dns")
	config := fmt.Sprintf(networkTemp, ip, netmask, gateway, dns)
	ioutil.WriteFile(netConfFile, []byte(config), 0644)
}
