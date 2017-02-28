package storage

import (
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
)

// Queue stores data
type Queue struct {
	data   map[int64]map[string]*models.Queue
	Status interface{}
}

// All returns queue map for userID
func (s *Queue) All(userID int64) map[string]*models.Queue {
	c := map[string]*models.Queue{}
	for k := range s.data[userID] {
		c[k] = s.data[userID][k]
	}

	return c
}

// One returns a queue for the userID with the queueName
func (s *Queue) One(userID int64, queueName string) (*models.Queue, error) {
	queues, hasUser := s.data[userID]
	if !hasUser {
		return nil, errors.ErrUserNotFound
	}

	queue, hasQueue := queues[queueName]
	if !hasQueue {
		return nil, errors.ErrQueueNotFound
	}

	return queue, nil
}

// Add add q for userID
func (s *Queue) Add(q *models.Queue) error {
	if _, err := s.One(q.UserID, q.Name); err != errors.ErrQueueNotFound {
		return err
	}

	queues := s.data[q.UserID]

	if queues == nil {
		queues = map[string]*models.Queue{}
	}

	s.data[q.UserID][q.Name] = q
	return nil
}

// Remove removes the queue with the name
func (s *Queue) Remove(userID int64, queueName string) error {
	if _, err := s.One(userID, queueName); err != nil {
		return err
	}

	delete(s.data[userID], queueName)
	return nil
}

// Message returns a message
func (s *Queue) Message(userID int64, queueName string, index int64) (*models.Message, error) {
	queue, err := s.One(userID, queueName)
	if err != nil {
		return nil, err
	}

	message, ok := queue.Messages[index]
	if !ok {
		return nil, errors.ErrMessageNotFound
	}

	return message, nil
}

// AddMessage adds m
func (s *Queue) AddMessage(m *models.Message) error {
	if _, err := s.Message(m.UserID, m.QueueName, m.Index); err != errors.ErrMessageNotFound {
		return err
	}

	s.data[m.UserID][m.QueueName].Messages[m.Index] = m
	return nil
}

// Messages returns message array for queueName of userID
func (s *Queue) Messages(userID int64, queueName string) (map[int64]*models.Message, error) {
	queue, err := s.One(userID, queueName)
	if err != nil {
		return nil, err
	}

	return queue.Messages, nil
}

// AddReceiver adds new receiver for one message
func (s *Queue) AddReceiver(userID int64, queueName string, messageIndex int64, clientID int64) error {
	message, err := s.Message(userID, queueName, messageIndex)
	if err != nil {
		return err
	}

	message.AddReceiver(clientID)
	return nil
}

// Receivers returns receivers of one message
func (s *Queue) Receivers(userID int64, queueName string, messageIndex int64) ([]int64, error) {
	message, err := s.Message(userID, queueName, messageIndex)
	if err != nil {
		return nil, err
	}

	return message.Recievers, nil
}

// DefaultQueue default
var DefaultQueue = &Queue{
	data: map[int64]map[string]*models.Queue{
		external.Root.ID(): map[string]*models.Queue{},
	},
}
