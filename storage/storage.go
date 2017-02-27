package storage

import (
	"github.com/Focinfi/sqs/queue"
)

// Storage stores data
type Storage struct {
	Status interface{}
}

// Queues returns queue map for userID
func (s *Storage) Queues(userID int64) map[string]queue.Queue {
	return nil
}

// AddQueue add q for userID
func (s *Storage) AddQueue(q queue.Queue) error {
	return nil
}

// AddMessage adds m
func (s *Storage) AddMessage(m queue.Message) error {
	return nil
}

// Messages returns message array for queueName of userID
func (s *Storage) Messages(userID int64, queueName string) []queue.Message {
	return []queue.Message{}
}

// NextMessageIndex returns next message index
func (s *Storage) NextMessageIndex(userID int64, queueName string) int64 {
	return 0
}

// AddReceiver adds new receiver for one message
func (s *Storage) AddReceiver(userID int64, queueName string, nessageIndex int) error {
	return nil
}

// Receivers returns receivers of one message
func (s *Storage) Receivers(userID int64, queueName string, nessageIndex int) []int64 {
	return []int64{}
}

// DefaultStorage default storage
var DefaultStorage = Storage{}
