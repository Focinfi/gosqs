package models

import "fmt"

// Queue contains name-message map
type Queue struct {
	UserID int64
	Name   string
}

// NewQueue returns a new queue
func NewQueue(userID int64, name string) *Queue {
	return &Queue{
		UserID: userID,
		Name:   name,
	}
}

const queueKeyPrefix = "sqs.queue"

// QueueListKey for queue list storage key
func QueueListKey(userID int64) string {
	return fmt.Sprintf("%s.%d", queueKeyPrefix, userID)
}

// QueueRecentMessageIDKey for queue storage key, the value is the almost recent message group id
func QueueRecentMessageIDKey(userID int64, queueName string) string {
	return fmt.Sprintf("%s.%d.%s", queueKeyPrefix, userID, queueName)
}

// QueueMaxIDKey for record the max id has been distributed
func QueueMaxIDKey(userID int64, queueName string) string {
	return fmt.Sprintf("%s.maxId.%d.%s", queueKeyPrefix, userID, queueName)
}
