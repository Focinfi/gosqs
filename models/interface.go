package models

// Consumer defines a consumer
type Consumer interface {
	Client() (client *Client)
	SetClient(client *Client)
	IncPriority(p int)
}

// PriorityList defines as priority list
type PriorityList interface {
	Push(item Consumer)
	Pop() (item Consumer)
}

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Delete(key string) error
}
