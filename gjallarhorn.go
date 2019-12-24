package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"log"
	"net"
	"os"
	"strings"
    "context"
)

var domainFlag = getopt.StringLong("domain", 'd', "invalid", "the base domain used to fully qualify hostnames.")

func main() {
    ctx := context.Background()
	getopt.Parse()
	hostname := findPublicHostname(*domainFlag)
	publicIPs := findPublicIPs()
    resolvedIPs, err := net.DefaultResolver.LookupIPAddr(ctx, hostname)
    if err != nil {
        log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()));
        return
    }
    fmt.Printf("publicIPs: %v\nresolvedIPs: %v\n", publicIPs, resolvedIPs)

}

func findPublicHostname(domain string) string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return strings.SplitN(hostname, ".", 2)[0] + "." + domain
}

func findPublicIPs() []net.IP {
	var ips []net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
		return ips
	}

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
			ips = append(ips, ip)
		}
	}

	return ips
}

func skipInterface(iface net.Interface) bool {
	return iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0
}

var _, uniqueLocalNetwork, _ = net.ParseCIDR("fc00::/7")
var _, linkLocalNetwork, _ = net.ParseCIDR("fe80::/10")
var _, legacyNetwork, _ = net.ParseCIDR("0.0.0.0/0")

func skipIP(ip net.IP) bool {
	return uniqueLocalNetwork.Contains(ip) || linkLocalNetwork.Contains(ip) || legacyNetwork.Contains(ip)
}
