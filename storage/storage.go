package storage

import (
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
)

// Storage defines storage
type Storage struct {
	*Queue
	*Message
	*Client
	*Cache
}

// DefaultStorage default storage
var DefaultStorage = &Storage{}

func init() {
	DefaultStorage.Queue = &Queue{db: defaultKV, store: DefaultStorage}
	DefaultStorage.Message = &Message{db: defaultKV, store: DefaultStorage}
	DefaultStorage.Client = &Client{db: defaultKV, store: DefaultStorage}
	DefaultStorage.Cache = NewCache(DefaultStorage)

	DefaultStorage.Queue.db.Put(models.QueueListKey(external.Root.ID()), "")
}
