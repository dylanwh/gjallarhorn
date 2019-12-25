package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/pborman/getopt/v2"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/* Unique Local Addresses prefix is fc00::/7 */
var _, uniqueLocalNetwork, _ = net.ParseCIDR("fc00::/7")

/* Link Local addresses are fe80::/10 */
var _, linkLocalNetwork, _ = net.ParseCIDR("fe80::/10")

/* Consider the entire IPv4 internet to be legacy */
var _, legacyNetwork, _ = net.ParseCIDR("0.0.0.0/0")

var domainFlag = getopt.StringLong("domain", 'd', "", "the base domain used to fully qualify hostnames (required)")
var monitorFlag = getopt.StringLong("monitor", 'm', "", "url of backend server (required)")
var secretKey = []byte(os.Getenv("GJALLARHORN_KEY"))

func main() {
	getopt.Parse()
	if *domainFlag == "" || *monitorFlag == "" {
		getopt.Usage()
		os.Exit(1)
	}
	if len(secretKey) == 0 {
		fmt.Printf("GJALLARHORN_KEY is not set\n")
		os.Exit(1)
	}

	hostname := findFullHostname()
	ips := findAddresses()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				LocalAddr: &net.TCPAddr{IP: ips[0]},
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	values := url.Values{}
	authen := hmac.New(sha256.New, secretKey)
	authen.Write([]byte(hostname))
	values.Set("hostname", hostname)
	for _, ip := range ips {
		ip := ip.String()
		values.Add("ip", ip)
		authen.Write([]byte(ip))
	}
	values.Set("authen", base64.RawStdEncoding.EncodeToString(authen.Sum(nil)))

	client.PostForm(*monitorFlag, values)
}

func findFullHostname() string {
	domain := *domainFlag
	if domain[0] != '.' {
		domain = "." + domain
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return strings.SplitN(hostname, ".", 2)[0] + domain
}

/*
 * This returns a list of IPv6 addresses that are (probably) routable. outable
 * means they're usable across the public internet and not just a LAN. *
 */
func findAddresses() []net.IP {
	var ips []net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("gjallarhorn: %v\n", err.Error()))
		return ips
	}

	for _, iface := range ifaces {
		/* we don't bother with loopback (localhost) or point-to-point (vpn?) interfaces */
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

			/*
			 * we ignore ULA, link local, and legacy (IPv4) ips. Anything that is
			 * not one of those is probably a routable IPv6 address.
			 */
			if uniqueLocalNetwork.Contains(ip) || linkLocalNetwork.Contains(ip) || legacyNetwork.Contains(ip) {
				continue
			}
			ips = append(ips, ip)
		}
	}

	return ips
}
