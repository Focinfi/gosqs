package admin

import (
	"github.com/Focinfi/sqs/queue"
	"github.com/Focinfi/sqs/storage"
)

type queues struct {
	storage.Storage
}

var defualtQueues = queues{Storage: storage.DefaultStorage}

func (qs *queues) PushMessage(userID int64, name, content string) error {
	msg := queue.Message{
		UserID:    userID,
		QueueName: name,
		Content:   content,
	}

	return defaultStorage.AddMessage(msg)
}

// AddQueue adds a queue into root queues
func AddQueue(name string, configs ...queue.Config) error {
	q := queue.New(configs...)
	return defualtQueues.AddQueue(q)
}
