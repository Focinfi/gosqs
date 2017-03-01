package service

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = database{Storage: storage.DefaultStorage}

func (d database) PushMessage(userID int64, queueName, content string) error {
	index := messageIndex()
	msg := &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Content:   content,
		Index:     index,
	}

	return d.Message.Add(msg)
}

func (d database) RegisterClient(userID int64, clientID int64, queueName string) error {
	client := &models.Client{
		ID:        clientID,
		UserID:    userID,
		QueueName: queueName,
	}

	return d.Client.Add(client)
}

// AddQueue adds a queue into root queues
func AddQueue(q *models.Queue) error {
	return db.Queue.Add(q)
}

// Queues returns all queues
func Queues(userID int64) ([]models.Queue, error) {
	return db.Queue.All(userID)
}
