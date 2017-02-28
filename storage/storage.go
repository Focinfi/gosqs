package storage

// Storage defines storage
type Storage struct {
	*Queue
	*Message
	*Client
}

// DefaultStorage default storage
var DefaultStorage = &Storage{
	Queue:   DefaultQueue,
	Message: DefaultMessage,
	Client:  DefaultClient,
}
