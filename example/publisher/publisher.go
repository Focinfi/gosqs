package main

import (
	"github.com/Focinfi/gosqs/client"
	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/example"
	"github.com/Focinfi/gosqs/log"
)

func main() {
	cli := client.New(config.Config.DefaultMasterAddress, "", "")
	queue, err := cli.Queue(example.Greeting, example.Home)
	if err != nil {
		log.Internal.Fatalln("failed to create a queue, err:", err)
	}

	if err := queue.ApplyNode(); err != nil {
		log.Internal.Fatalln("failed to apply node, err:", err)
	}

	if err := queue.PushMessage("first message"); err != nil {
		log.Internal.Fatalln("failed to push message, err:", err)
	}
}
