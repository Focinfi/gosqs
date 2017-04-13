package storage

import (
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
)

// Storage defines storage
type Storage struct {
	*Queue
	*Message
	*Squad
}

// DefaultStorage default storage
var DefaultStorage = &Storage{}

func init() {
	DefaultStorage.Queue = &Queue{db: EtcdKV, store: DefaultStorage, inc: etcdIncrementer}
	DefaultStorage.Message = &Message{db: EtcdKV, store: DefaultStorage}
	DefaultStorage.Squad = &Squad{db: EtcdKV, store: DefaultStorage}

	DefaultStorage.Queue.db.Put(models.QueueListKey(external.Root.ID()), "[]")
}
