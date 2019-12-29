package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
	"github.com/dylanwh/gjallarhorn/useragent"
)

func main() {
	cfg := config.NewClient()
	cfg.CheckArgs()
	msg, err := message.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if msg.PublicIP == nil {
		if cfg.Verbose() {
			log.Fatalf("No IPv6 on interface %s\n", cfg.IfName())
		}
		return
	}
	if msg.KnownIP == nil {
		if cfg.Verbose() {
			log.Fatalf("No IPv6 for hostname %s\n", msg.FullHostname)
		}
		return
	}
	if !msg.KnownIP.Equal(*msg.PublicIP) {
		if cfg.Verbose() {
			log.Fatalf(
				"%s is misconfigured.\n  %s was found\n  %s was expected\n",
				msg.FullHostname,
				msg.KnownIP,
				msg.PublicIP,
			)
		}
		return
	}

	if cfg.Debug() {
		buf, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "json error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(buf))
		return
	}

	ua := useragent.New(cfg)
	resp, err := ua.Send(msg)
	if err != nil {
		fmt.Printf("http error: %s\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("http status: %d\n", resp.StatusCode)
	}
}
