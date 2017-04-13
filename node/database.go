package node

import (
	"fmt"

	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

func (d *database) PushMessage(userID int64, queueName, content string, index int64) error {
	log.Internal.Infoln("[PushMessage]", index)
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
	// time.AfterFunc(time.Second, func() {
	// 	err := d.Queue.UpdateRecentMessageID(userID, queueName, index)
	// 	if err != nil {
	// 		log.DB.Error(err)
	// 	}
	// })

	return nil
}

func (d *database) updateSquadReceivedMessageID(userID int64, queueName, squadName string, messageID int64) error {
	squad, err := d.Squad.One(userID, queueName, squadName)
	if err != nil {
		return err
	}

	if squad.ReceivedMessageID >= messageID {
		return nil
	}

	squad.ReceivedMessageID = messageID
	return d.Squad.Update(squad)
}
