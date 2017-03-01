package storage

import "github.com/Focinfi/sqs/models"

// Message for message storage
type Message struct {
	db KV
}

// One returns a message
func (s *Message) One(userID int64, queueName string, index int64) (*models.Message, error) {
	return nil, nil
}

// Add adds a message
func (s *Message) Add(m *models.Message) error {
	return nil
}

// All returns all messages
func (s *Message) All(userID int64, queueName string, filter interface{}) (map[int64]*models.Message, error) {
	return nil, nil
}

// DefaultMessage default message
var DefaultMessage = &Message{db: defaultKV}
