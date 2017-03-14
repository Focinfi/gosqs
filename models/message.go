package models

import "fmt"

// Message contains info
type Message struct {
	UserID    int64
	QueueName string
	Content   string
	Index     int64
}

const messageKeyPrefix = "sqs.message"

// MessageListKey for message list storage key
func MessageListKey(userID int64, queueName string, gorupID int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, gorupID)
}

// MessageKey for message storage key
func MessageKey(userID int64, queueName string, index int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, index)
}
