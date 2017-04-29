package node

import (
	"github.com/Focinfi/gosqs/models"
	"github.com/Focinfi/gosqs/storage"
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
		return err
	}
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
