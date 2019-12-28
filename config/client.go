package config

import (
	"os"

	"github.com/pborman/getopt"
)

type Client struct {
	domain  *string
	monitor *string
	key     *string
	ifname  *string
}

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

func (c *Client) CheckArgs() {
	if *c.domain == "" || *c.monitor == "" || *c.key == "" {
		getopt.Usage()
		os.Exit(1)
	}
}

func (c *Client) Domain() string {
	return *c.domain
}

func (c *Client) Monitor() string {
	return *c.monitor
}

func (c *Client) Key() string {
	return *c.key
}

func (c *Client) IfName() string {
	return *c.ifname
}
