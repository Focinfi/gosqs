package main

import (
	"time"

	"github.com/Focinfi/sqs/client"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/log"
)

func main() {
	cli := client.New(config.Config().DefaultMasterAddress, "", "")
	queue, err := cli.Queue(example.Greeting, example.Home)
	if err != nil {
		panic(err)
	}
	if err := queue.ApplyNode(); err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C

		err := queue.PullMessage(func(messages []client.Message) error {
			log.Biz.Infoln(messages)
			return nil
		})

		if err != nil {
			log.Biz.Error(err)
			continue
		}
	}
}
