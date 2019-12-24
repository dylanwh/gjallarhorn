package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var _, uniqueLocalNetwork, _ = net.ParseCIDR("fc00::/7")
var _, linkLocalNetwork, _ = net.ParseCIDR("fe80::/10")
var _, legacyNetwork, _ = net.ParseCIDR("0.0.0.0/0")

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
	}
	fmt.Printf("Hostname is %s\n", hostname)

	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
				continue
			}
			if skipIP(ip) {
				continue
			}
			fmt.Printf("ip = %v\n", ip)
		}
	}
}

func skipInterface(iface net.Interface) bool {
	return iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0
}

func skipIP(ip net.IP) bool {
	return uniqueLocalNetwork.Contains(ip) || linkLocalNetwork.Contains(ip) || legacyNetwork.Contains(ip)
}
