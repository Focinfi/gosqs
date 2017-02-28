package main

import "github.com/Focinfi/sqs/admin"
import "github.com/Focinfi/sqs/models"
import "github.com/Focinfi/sqs/external"

func main() {
	queue := models.NewQueue(external.Root.ID(), "greeting")

	if err := admin.AddQueue(queue); err != nil {
		panic(err)
	}
	admin.Start(":5546")
}
