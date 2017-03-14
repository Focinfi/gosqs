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
	DefaultStorage.Queue = &Queue{db: etcdKV, store: DefaultStorage, inc: etcdIncrementer}
	DefaultStorage.Message = &Message{db: memcachedKV, store: DefaultStorage}
	DefaultStorage.Client = &Client{db: etcdKV, store: DefaultStorage}
	DefaultStorage.Cache = &Cache{pl: redisPriorityList, watcher: etcdWatcher, store: DefaultStorage}

	DefaultStorage.Queue.db.Put(models.QueueListKey(external.Root.ID()), "[]")
}
