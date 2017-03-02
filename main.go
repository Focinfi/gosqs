package main

import (
	"time"

	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/service"
)

func main() {
	greeting := "greeting"
	queue := models.NewQueue(external.Root.ID(), greeting)
	if err := service.AddQueue(queue); err != nil {
		panic(err)
	}

	client := &models.Client{
		ID:                 external.TestClient.ID(),
		UserID:             external.Root.ID(),
		QueueName:          greeting,
		Addresses:          []string{":55466"},
		RecentMessageIndex: models.GenIndex0(time.Now().Unix()),
		RecentReceivedAt:   time.Now().Unix(),
	}

	if err := service.AddClient(client); err != nil {
		panic(err)
	}

	service.Start(":5546")
}
