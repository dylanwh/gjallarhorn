package useragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
)

/*UserAgent ...*/
type UserAgent struct {
	config *config.Client
	client *http.Client
}

/*New ...*/
func New(c *config.Client) *UserAgent {
	return &UserAgent{config: c, client: &http.Client{Transport: defaultTransport}}
}

/*Monitor ...*/
func (ua *UserAgent) Monitor() string {
	return ua.config.Monitor()
}

/*Send ...*/
func (ua *UserAgent) Send(msg *message.Message) (*http.Response, error) {
	json, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(json)
	req, err := http.NewRequest("POST", ua.Monitor(), body)
	fmt.Println(req)
	if err != nil {
		return nil, fmt.Errorf("unable to create POST request: %w", err)
	}

	if msg.PublicIP != nil {
		req.Header.Set("Source-IP", msg.PublicIP.String())
	}

	sig, err := msg.Sign(ua.config)
	if err != nil {
		return nil, fmt.Errorf("unable to sign message: %w", err)
	}
	req.Header.Set("Signature", sig)

	resp, err := ua.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("unable to send POST request: %w", err)
	}

	return resp, nil
}
