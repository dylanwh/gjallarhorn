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
	Hostname        string
	FullHostname    string
	KnownIP         *net.IP
	Ifname          string
	PublicIP        *net.IP
	OperatingSystem string
}

var _, globalUnicastNetwork, _ = net.ParseCIDR("2000::/3")

/*New constructs a message containing all the information of the system it is running on. */
func New(c *config.Client) (*Message, error) {
	hostname, fullHostname, err := findHostname(c.Domain())
	if err != nil {
		return nil, err
	}

	iface, err := findInterface(c.IfName())
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Hostname:        hostname,
		FullHostname:    fullHostname,
		KnownIP:         findKnownIP(fullHostname),
		PublicIP:        findPublicIP(iface),
		Ifname:          c.IfName(),
		OperatingSystem: runtime.GOOS,
	}

	return msg, nil
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

func findHostname(domain string) (hostname string, fullHostname string, err error) {
	hostname, err = os.Hostname()
	if err != nil {
		return
	}
	if hostname == "" {
		err = fmt.Errorf("system hostname is blank?")
		return
	}

	if strings.HasSuffix(hostname, domain) {
		fullHostname = hostname
	} else {
		if strings.Contains(hostname, ".") {
			parts := strings.SplitN(hostname, ".", 2)
			fullHostname = parts[0] + domain
		} else {
			fullHostname = hostname + domain
		}
	}

	return
}

func findKnownIP(fullHostname string) *net.IP {
	resolver := &net.Resolver{PreferGo: true}
	ips, err := resolver.LookupIP(fullHostname)
	if err != nil {
		log.Println(fmt.Errorf("Unable to lookup ip for %s: %w", fullHostname, err))
		return nil
	}

	for _, ip := range ips {
		if globalUnicastNetwork.Contains(ip) {
			return &ip
		}
	}
	return nil
}

/*
 * This returns a list of IPv6 addresses that are (probably) routable. outable
 * means they're usable across the public internet and not just a LAN. *
 */
func findInterface(ifname string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == ifname {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("interface not found: %s", ifname)
}

func findPublicIP(iface *net.Interface) *net.IP {
	addrs, err := iface.Addrs()
	if err != nil {
		log.Println(fmt.Errorf("Unable to get ips from interface %s: %w", iface.Name, err))
		return nil
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			log.Println(fmt.Errorf("Unable to parse CIDR from %s: %v", addr.String(), err))
			continue
		}
		if globalUnicastNetwork.Contains(ip) {
			return &ip
		}
	}

	return nil
}
