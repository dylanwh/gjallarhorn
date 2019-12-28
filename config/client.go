package config

import (
	"os"

	"github.com/pborman/getopt"
)

/*Client is the client configuration */
type Client struct {
	domain  *string
	monitor *string
	key     *string
	ifname  *string
}

/*NewClient ...*/
func NewClient() *Client {
	c := &Client{
		domain:  getopt.StringLong("domain", 'd', "", "the base domain used to fully qualify hostnames (required)"),
		monitor: getopt.StringLong("monitor", 'm', "", "url of backend server (required)"),
		key:     keyflag(),
		ifname:  getopt.StringLong("ifname", 'i', ""),
	}
	getopt.Parse()
	return c
}

/*CheckArgs ... */
func (c *Client) CheckArgs() {
	if *c.domain == "" || *c.monitor == "" || *c.key == "" {
		getopt.Usage()
		os.Exit(1)
	}
}

/*Domain ...*/
func (c *Client) Domain() string {
	return *c.domain
}

/*Monitor ...*/
func (c *Client) Monitor() string {
	return *c.monitor
}

/*Key ...*/
func (c *Client) Key() string {
	return *c.key
}

/*IfName ...*/
func (c *Client) IfName() string {
	return *c.ifname
}
