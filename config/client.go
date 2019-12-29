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
	debug   *bool
	verbose *bool
}

/*NewClient ...*/
func NewClient() *Client {
	c := &Client{
		domain:  getopt.StringLong("domain", 'd', "", "the base domain used to fully qualify hostnames"),
		monitor: getopt.StringLong("monitor", 'm', "", "url of backend server"),
		key:     keyflag(),
		ifname:  getopt.StringLong("ifname", 'i', "", "name of the interface to look at"),
		debug:   getopt.BoolLong("debug", 'D', "debug mode"),
		verbose: getopt.BoolLong("verbose", 'v', "verbose mode"),
	}
	getopt.Parse()
	return c
}

/*CheckArgs ... */
func (c *Client) CheckArgs() {
	if *c.debug {
		var verbose = true
		c.verbose = &verbose
	}
	if *c.domain == "" || *c.ifname == "" || !*c.debug && (*c.monitor == "" || *c.key == "") {
		getopt.Usage()
		os.Exit(1)
	}
}

/*Domain ...*/
func (c *Client) Domain() string {
	d := *c.domain
	if len(d) > 0 && d[0] != '.' {
		return "." + d
	}
	return d
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

func (c *Client) Debug() bool {
	return *c.debug
}

func (c *Client) Verbose() bool {
	return *c.verbose
}
