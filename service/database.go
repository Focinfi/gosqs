package service

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = database{Storage: storage.DefaultStorage}

func (d database) PushMessage(userID int64, queueName, content string, index int64) error {
	msg := &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Content:   content,
		Index:     index,
	}

	return d.Message.Add(msg)
}

func (d database) RegisterClient(userID int64, clientID int64, queueName string) error {
	_, err := d.Client.One(userID, clientID, queueName)
	if err == nil {
		return err
	}

	// save into cache
	return nil
}

// AddQueue adds a queue into root queues
func AddQueue(q *models.Queue) error {
	return db.Queue.Add(q)
}

// AddClient adds cleint
func AddClient(client *models.Client) error {
	return db.Client.Add(client)
}

// Queues returns all queues
func Queues(userID int64) ([]models.Queue, error) {
	return db.Queue.All(userID)
}
