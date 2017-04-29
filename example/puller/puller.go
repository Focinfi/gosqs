package main

import (
	"time"

	"github.com/Focinfi/gosqs/client"
	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/example"
	"github.com/Focinfi/gosqs/log"
)

func main() {
	cli := client.New(config.Config.DefaultMasterAddress, "", "")
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

		err := queue.PullMessages(func(messages []client.Message) error {
			log.Biz.Infoln(messages)
			return nil
		})

		if err != nil {
			log.Biz.Error(err)
			continue
		}
	}
}
