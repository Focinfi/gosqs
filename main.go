package main

import (
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/service"
)

func main() {
	queue := models.NewQueue(external.Root.ID(), "greeting")

	if err := service.AddQueue(queue); err != nil {
		panic(err)
	}
	service.Start(":5546")
}
