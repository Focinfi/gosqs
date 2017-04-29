package storage

import "github.com/Focinfi/gosqs/models"

// Storage defines storage
type Storage struct {
	*Nodes
	*Queue
	*Message
	*Squad
}

// DefaultStorage default storage
var DefaultStorage = &Storage{}

func init() {
	DefaultStorage.Queue = &Queue{db: sqsMetaKV, store: DefaultStorage, inc: sqsMetaIncrementer}
	DefaultStorage.Message = &Message{db: messageKV, store: DefaultStorage}
	DefaultStorage.Squad = &Squad{db: sqsMetaKV, store: DefaultStorage}
	DefaultStorage.Nodes = &Nodes{db: ClusterMetaKV, store: DefaultStorage}

	// TODO: move into db/seeds
	DefaultStorage.Queue.db.Put(models.QueueListKey(1), "[]")
}
