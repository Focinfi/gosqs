package master

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

// AddQueue adds a queue into root queues
func AddQueue(q *models.Queue) error {
	if err := db.Queue.Add(q); err != nil {
		return err
	}

	return nil
}

// Queues returns all queues
func Queues(userID int64) ([]models.Queue, error) {
	return db.Queue.All(userID)
}
