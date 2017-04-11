package node

import (
	"fmt"
	"time"

	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

func (d *database) PushMessage(userID int64, queueName, content string, index int64) error {
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

func (d *database) updateSquadReceivedMessageID(userID int64, queueName, squadName string, index int64) error {
	//key := models.SquadKey(userID, queueName, squadName)

	return nil
}
