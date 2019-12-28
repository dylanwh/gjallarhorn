package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/dylanwh/gjallarhorn/config"
)

/*Message is all the information that gjallarhorn sends to gjallarheim. */
type Message struct {
	Hostname           string
	Ifname             string
	PublishedAddress   net.IP
	InterfaceAddresses map[string][]net.IP
	OperatingSystem    string
}

/* Unique Local Addresses prefix is fc00::/7 */
var _, uniqueLocalNetwork, _ = net.ParseCIDR("fc00::/7")

/* Link Local addresses are fe80::/10 */
var _, linkLocalNetwork, _ = net.ParseCIDR("fe80::/10")

/* Consider the entire IPv4 internet to be legacy */
var _, legacyNetwork, _ = net.ParseCIDR("0.0.0.0/0")

/*New constructs a message containing all the information of the system it is running on. */
func New(c *config.Client) *Message {
	hostname := findFullHostname(c.Domain())
	publishedAddr := findPublishedAddress(hostname)
	ifAddrs := findAddresses()

	return &Message{
		Hostname:           hostname,
		Ifname:             c.IfName(),
		PublishedAddress:   publishedAddr,
		InterfaceAddresses: ifAddrs,
		OperatingSystem:    runtime.GOOS,
	}
}

func (m *Message) Sign(k config.Keyer) (string, error) {
	if k == nil {
		return "", errors.New("got nil Keyer")
	}
	h := hmac.New(sha256.New, []byte(k.Key()))
	enc := json.NewEncoder(h)
	if err := enc.Encode(m); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func findPublishedAddress(hostname string) net.IP {

	ips, err := net.LookupIP(hostname)
	if err != nil {
		log.Println(fmt.Errorf("gjallarhorn: %v", err.Error()))
		return nil
	}
	if len(ips) > 0 {
		return ips[0]
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

/*
 * This returns a list of IPv6 addresses that are (probably) routable. outable
 * means they're usable across the public internet and not just a LAN. *
 */
func findAddresses() map[string][]net.IP {
	var ips map[string][]net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println(fmt.Errorf("gjallarhorn: %v", err.Error()))
	}

	for _, iface := range ifaces {
		/* we don't bother with loopback (localhost) or point-to-point (vpn?) interfaces */
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Println(fmt.Errorf("gjallarhorn: %v", err.Error()))
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Println(fmt.Errorf("gjallarhorn: %v", err.Error()))
				continue
			}

			/*
			 * we ignore ULA, link local, and legacy (IPv4) ips. Anything that is
			 * not one of those is probably a routable IPv6 address.
			 */
			if uniqueLocalNetwork.Contains(ip) || linkLocalNetwork.Contains(ip) || legacyNetwork.Contains(ip) {
				continue
			}
			ips[iface.Name] = append(ips[iface.Name], ip)
			ips["ALL"] = append(ips["ALL"], ip)
		}
	}

	return ips
}
