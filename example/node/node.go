package main

import (
	"time"

	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/service"
)

func main() {
	queue := models.NewQueue(external.Root.ID(), example.Greeting)
	if err := service.AddQueue(queue); err != nil {
		panic(err)
	}

	now := time.Now().Unix()

	clients := []*models.Client{
		{
			ID:                 external.TestClient.ID(),
			UserID:             external.Root.ID(),
			QueueName:          example.Greeting,
			Addresses:          []string{"http://localhost:55466/greeting/1"},
			RecentMessageIndex: models.GenIndex0(now),
			RecentPushedAt:     now,
			RecentReceivedAt:   now,
		},
		{
			ID:                 external.TestClient.ID() + 1,
			UserID:             external.Root.ID(),
			QueueName:          example.Greeting,
			Addresses:          []string{"http://localhost:55466/greeting/2"},
			RecentMessageIndex: models.GenIndex0(now),
			RecentPushedAt:     now,
			RecentReceivedAt:   now,
		},
		{
			ID:                 external.TestClient.ID() + 2,
			UserID:             external.Root.ID(),
			QueueName:          example.Greeting,
			Addresses:          []string{"http://localhost:55466/greeting/3"},
			RecentMessageIndex: models.GenIndex0(now),
			RecentPushedAt:     now,
			RecentReceivedAt:   now,
		},
	}

	for _, client := range clients {
		if err := service.AddClient(client); err != nil {
			panic(err)
		}
	}
	service.Start(":5546")
}
