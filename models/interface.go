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
	Add(item Consumer) error
	Push(item Consumer) error
	Pop() (item Consumer, err error)
	Remove(item Consumer) error
}

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
}

// Incrementer defines increment a integer value in a transaction
type Incrementer interface {
	// Increment try to increment the current value of the key.
	// It will returns non-nil error if current value exists but not a number string or
	// failed to increment by resource race.
	// result is the result value of the key if the increment succeed.
	Increment(key string, number int) (result int64, err error)
}

// Watcher defines a watcher
type Watcher interface {
	Watch(key string) (value <-chan string)
	Close() error
}
