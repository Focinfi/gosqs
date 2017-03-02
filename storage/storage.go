package storage

import "github.com/Focinfi/sqs/external"

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
	DefaultStorage.Cache = &Cache{KV: defaultKV, store: DefaultStorage}

	DefaultStorage.Queue.db.Put(queueListKey(external.Root.ID()), "")
}
