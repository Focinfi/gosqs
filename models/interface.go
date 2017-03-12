package models

// Consumer defines a consumer
type Consumer interface {
	Client() (client *Client)
	SetClient(client *Client)
	Priority() (p int)
	IncPriority(p int)
}

// PriorityList defines as priority list
type PriorityList interface {
	Push(item Consumer) error
	Pop() (item Consumer, err error)
}

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Delete(key string) error
}

// Watcher defines a wacther
type Watcher interface {
	Watch(key string) (value <-chan string)
	Close() error
}
