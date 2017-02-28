package admin

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type queues struct {
	storage.Storage
}

var defualtQueues = queues{Storage: storage.DefaultStorage}

func (qs queues) PushMessage(userID int64, queueName, content string) error {
	index := int64(1)
	msg := &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Content:   content,
		Index:     index,
	}

	return defualtQueues.AddMessage(msg)
}

// AddQueue adds a queue into root queues
func AddQueue(q *models.Queue) error {
	return defualtQueues.AddQueue(q)
}
