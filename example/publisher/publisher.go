package main

import (
	"github.com/Focinfi/sqs/client"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/log"
)

func main() {
	cli := client.New(config.Config().DefaultMasterAddress, "", "")
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
