package useragent

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type transport struct {
	transport http.Transport
}

var defaultTransport http.RoundTripper = &transport{
	transport: http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           newDialer(nil).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Gjallarhorn")

	sourceIP := req.Header.Get("Source-IP")
	if sourceIP != "" {
		ip := net.ParseIP(sourceIP)
		if ip != nil {
			t.transport.DialContext = newDialer(ip).DialContext
		} else {
			return nil, fmt.Errorf("%s is not a valid IP address", sourceIP)
		}
	}

	return t.transport.RoundTrip(req)
}

func newDialer(ip net.IP) *net.Dialer {
	d := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	if ip != nil {
		d.LocalAddr = &net.TCPAddr{IP: ip}
	}
	return d
}
