package master

import (
	"github.com/Focinfi/gosqs/models"
	"github.com/Focinfi/gosqs/storage"
)

type database struct {
	*storage.Storage
}

var db = &database{Storage: storage.DefaultStorage}

func (d *database) fetchNodes() ([]string, error) {
	return d.Storage.Nodes.All(models.NodesKey)
}

func (d *database) updateNodes(nodes []string) error {
	return d.Storage.Nodes.UpdateAll(models.NodesKey, nodes)
}

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
