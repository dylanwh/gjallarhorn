package useragent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
)

type UserAgent struct {
	config *config.Client
	client *http.Client
}

func New(c *config.Client) *UserAgent {
	return &UserAgent{config: c, client: &http.Client{Transport: DefaultTransport}}
}

func (ua *UserAgent) Monitor() string {
	return ua.config.Monitor()
}

func (ua *UserAgent) Send(msg *message.Message) (*http.Response, error) {
	json, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(json)
	req, err := http.NewRequest("POST", ua.Monitor(), body)
	if err != nil {
		return nil, err
	}
	if msg.PublishedAddress != nil {
		req.Header.Set("Source-IP", msg.PublishedAddress.String())
	} else if len(msg.InterfaceAddresses) > 0 {
		req.Header.Set("Source-IP", msg.InterfaceAddresses[0].String())
	}

	sig, err := msg.Sign(ua.config)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Signature", sig)

	resp, err := ua.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
