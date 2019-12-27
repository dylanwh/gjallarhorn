package config

import (
	"os"

	"github.com/pborman/getopt"
)

type Config struct {
	domain  *string
	monitor *string
	key     *string
	ifname  *string
}

func ParseArgs() *Config {
	var config = &Config{
		domain:  getopt.StringLong("domain", 'd', "", "the base domain used to fully qualify hostnames (required)"),
		monitor: getopt.StringLong("monitor", 'm', "", "url of backend server (required)"),
		key:     getopt.StringLong("key", 'k', os.Getenv("GJALLARHORN_KEY"), "secret key for signature of monitor message."),
		ifname:  getopt.StringLong("ifname", 'i', ""),
	}
	getopt.Parse()
	if *config.domain == "" || *config.monitor == "" || *config.key == "" {
		getopt.Usage()
		os.Exit(1)
	}

	return config
}

func (c *Config) Domain() string {
	return *c.domain
}

func (c *Config) Monitor() string {
	return *c.monitor
}

func (c *Config) Key() []byte {
	return []byte(*c.key)
}

func (c *Config) IfName() string {
	return *c.ifname
}
