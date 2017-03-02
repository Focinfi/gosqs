package service

import (
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

func (d *database) ReceivehMessage(userID int64, queueName, content string, index int64) error {
	msg := &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Content:   content,
		Index:     index,
	}

	return d.Message.Add(msg)
}

func (d *database) RegisterClient(c *models.Client) error {
	client, err := d.Client.One(c.UserID, c.ID, c.QueueName)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	// the client had received message in clientControlTimeoutSecond, can not register for this node
	if now-client.RecentReceivedAt < config.Config().ClientControlTimeoutSecond {
		return errors.ClientHasAlreadyRegistered
	}

	c.RecentMessageIndex = client.RecentMessageIndex
	c.RecentReceivedAt = now
	return d.Client.Update(c)
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
