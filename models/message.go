package models

import "fmt"

// Message contains info
type Message struct {
	UserID    int64  `json:"-"`
	QueueName string `json:"-"`
	Content   string `json:"content"`
	Index     int64  `json:"message_id"`
}

const messageKeyPrefix = "sqs.message"

// MessageKey for message storage key
func MessageKey(userID int64, queueName string, index int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, index)
}
