package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pborman/getopt/v2"
)

var _, uniqueLocalNetwork, _ = net.ParseCIDR("fc00::/7")
var _, linkLocalNetwork, _ = net.ParseCIDR("fe80::/10")
var _, legacyNetwork, _ = net.ParseCIDR("0.0.0.0/0")

var domainFlag = getopt.StringLong("domain", 'd', "", "the base domain used to fully qualify hostnames.")

func main() {
	getopt.Parse()
	if *domainFlag == "" {
		getopt.Usage()
		return
	}

	ctx := context.Background()
	hostname := findFullHostname(*domainFlag)
	ip := findPublishedAddress(ctx, hostname)
	ips := findAddresses()
	fmt.Printf("hostname = %v (%v)\naddresses = %v\n", hostname, ip, ips)
}

func findPublishedAddress(ctx context.Context, hostname string) net.IP {
	addrs, _ := net.DefaultResolver.LookupIPAddr(ctx, hostname)
	for _, addr := range addrs {
		if !legacyNetwork.Contains(addr.IP) {
			return addr.IP
		}
	}
	return nil
}

func findFullHostname(domain string) string {
	if domain[0] != '.' {
		domain = "." + domain
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return strings.SplitN(hostname, ".", 2)[0] + domain
}

func findAddresses() []net.IP {
	var ips []net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
		return ips
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0 {
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
			if uniqueLocalNetwork.Contains(ip) || linkLocalNetwork.Contains(ip) || legacyNetwork.Contains(ip) {
				continue
			}
			ips = append(ips, ip)
		}
	}

	return ips
}
