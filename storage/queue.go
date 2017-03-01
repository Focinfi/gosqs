package storage

import "github.com/Focinfi/sqs/models"

// Queue stores data
type Queue struct {
	db KV
}

// All returns queue map for userID
func (s *Queue) All(userID int64) (map[string]*models.Queue, error) {
	return nil, nil
}

// One returns a queue for the userID with the queueName
func (s *Queue) One(userID int64, queueName string) (*models.Queue, error) {
	return nil, nil
}

// Add add q for userID
func (s *Queue) Add(q *models.Queue) error {
	return nil
}

// Remove removes the queue with the name
func (s *Queue) Remove(userID int64, queueName string) error {
	return nil
}

// DefaultQueue default
var DefaultQueue = &Queue{
	db: defaultKV,
}
