package main

import (
	"fmt"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
	"github.com/dylanwh/gjallarhorn/useragent"
)

func main() {
	config := config.ParseArgs()
	msg := message.New(config)
	ua := useragent.New()
	resp, err := ua.Post(config.Monitor(), "application/json", msg.Reader())
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(resp)

}
