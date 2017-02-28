package storage

// Message for message storage
type Message struct {
	data interface{}
}

// DefaultMessage default message
var DefaultMessage = &Message{}
