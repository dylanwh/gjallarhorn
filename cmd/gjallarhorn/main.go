package main

import (
	"encoding/json"
	"fmt"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
	"github.com/dylanwh/gjallarhorn/useragent"
)

func main() {
	cfg := config.NewClient()
	cfg.CheckArgs()
	msg := message.New(cfg)
	if cfg.Debug() {
		buf, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Printf("%s\n", string(buf))
		return
	}
	ua := useragent.New(cfg)
	_, err := ua.Send(msg)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
