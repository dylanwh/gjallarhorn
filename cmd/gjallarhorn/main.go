package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
	"github.com/dylanwh/gjallarhorn/useragent"
)

func main() {
	cfg := config.NewClient()
	cfg.CheckArgs()
	msg := message.New(cfg)
	ua := useragent.New(cfg)
	resp, err := ua.Send(msg)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(buf))

}
