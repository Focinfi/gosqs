package service

import (
	"fmt"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

func (d *database) ReceiveMessage(userID int64, queueName, content string, index int64) error {
	msg := &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Content:   content,
		Index:     index,
	}

	if err := d.Message.Add(msg); err != nil {
		fmt.Println("Added Message, err: ", err)
		return err
	}

	// try to update recent message index in background
	time.AfterFunc(time.Second, func() {
		err := d.Queue.UpdateRecentMessageID(userID, queueName, index)
		if err != nil {
			log.DB.Error(err)
		}
	})

	return nil
}

func (d *database) RegisterClient(c *models.Client) (err error) {
	client, err := d.Client.One(c.UserID, c.ID, c.QueueName)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	// the client had received message in clientControlTimeoutSecond, can not register for this node
	if c.Publisher != client.Publisher &&
		now-client.RecentPushedAt < config.Config().ClientControlTimeoutSecond {
		return errors.ClientHasAlreadyRegistered
	}

	c.RecentMessageIndex = client.RecentMessageIndex
	c.RecentPushedAt = now
	c.RecentReceivedAt = client.RecentReceivedAt
	if c.RecentReceivedAt == 0 {
		c.RecentReceivedAt = now
	}
	log.Biz.Printf("RegisterClient: %v", c)

	return d.Client.Update(c)
}

// AddQueue adds a queue into root queues
func AddQueue(q *models.Queue) error {
	if err := db.Queue.Add(q); err != nil {
		return err
	}

	return nil
}

// AddClient adds cleint
func AddClient(client *models.Client) error {
	maxID, err := db.Queue.MessageMaxID(client.UserID, client.QueueName)
	if err != nil {
		return err
	}

	client.RecentMessageIndex = maxID
	return db.Client.Add(client)
}

// Queues returns all queues
func Queues(userID int64) ([]models.Queue, error) {
	return db.Queue.All(userID)
}
