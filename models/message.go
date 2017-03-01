package models

// Message contains info
type Message struct {
	UserID    int64
	QueueName string
	Content   string
	Index     int64
}

// MessageIndex for one entry of message index
type MessageIndex struct {
	Timestamp int64
	Indexes   []int64
}
