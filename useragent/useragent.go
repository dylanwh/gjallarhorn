package useragent

import (
	"io"
	"net"
	"net/http"
)

type UserAgent struct {
	SourceIP *net.IP
	client   *http.Client
}

func New() *UserAgent {
	return &UserAgent{SourceIP: nil, client: &http.Client{Transport: DefaultTransport}}
}

func (ua *UserAgent) Do(req *http.Request) (*http.Response, error) {

	if ua.SourceIP != nil {
		req.Header.Set("Source-IP", ua.SourceIP.String())
	}
	return ua.client.Do(req)
}

func (ua *UserAgent) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ua.Do(req)
}

func (ua *UserAgent) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return ua.Do(req)
}
