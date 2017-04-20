package main

import (
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/master"
	"github.com/Focinfi/sqs/models"
)

func main() {
	queue := models.NewQueue(1, example.Greeting)
	if err := master.AddQueue(queue); err != nil {
		panic(err)
	}

	master.NewService(config.Config.DefaultMasterAddress).Start()
}
