package storage

import (
	"log"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/models"
)

// Storage stores data
type Storage struct {
	queueMap map[int64]map[string]*models.Queue
	Status   interface{}
}

// Queuer defines what a queue can do
type Queuer interface {
	UserID() int64
	Name() string
	Messages() map[int64]Messager
}

// Messager defines what message can do
type Messager interface {
	UserID() int64
	QueueName() string
	Index() int64
	Content() string
	Recievers() []int64
	AddReceiver(clientID int64)
}

// Queues returns queue map for userID
func (s *Storage) Queues(userID int64) map[string]*models.Queue {
	c := map[string]*models.Queue{}
	for k := range s.queueMap[userID] {
		c[k] = s.queueMap[userID][k]
	}

	return c
}

// Queue returns a queue for the userID with the queueName
func (s *Storage) Queue(userID int64, queueName string) (*models.Queue, error) {
	queues, hasUser := s.queueMap[userID]
	if !hasUser {
		return nil, errors.ErrUserNotFound
	}

	queue, hasQueue := queues[queueName]
	if !hasQueue {
		return nil, errors.ErrQueueNotFound
	}

	return queue, nil
}

// AddQueue add q for userID
func (s *Storage) AddQueue(q *models.Queue) error {
	if _, err := s.Queue(q.UserID, q.Name); err != errors.ErrQueueNotFound {
		return err
	}

	queues := s.queueMap[q.UserID]

	if queues == nil {
		queues = map[string]*models.Queue{}
	}

	s.queueMap[q.UserID][q.Name] = q
	log.Println(s.queueMap)
	return nil
}

// RemoveQueue removes the queue with the name
func (s *Storage) RemoveQueue(userID int64, queueName string) error {
	if _, err := s.Queue(userID, queueName); err != nil {
		return err
	}

	delete(s.queueMap[userID], queueName)
	return nil
}

// Message returns a message
func (s *Storage) Message(userID int64, queueName string, index int64) (*models.Message, error) {
	queue, err := s.Queue(userID, queueName)
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
func (s *Storage) AddMessage(m *models.Message) error {
	if _, err := s.Message(m.UserID, m.QueueName, m.Index); err != errors.ErrMessageNotFound {
		return err
	}

	s.queueMap[m.UserID][m.QueueName].Messages[m.Index] = m
	log.Println(s.queueMap[m.UserID][m.QueueName])
	return nil
}

// Messages returns message array for queueName of userID
func (s *Storage) Messages(userID int64, queueName string) (map[int64]*models.Message, error) {
	queue, err := s.Queue(userID, queueName)
	if err != nil {
		return nil, err
	}

	return queue.Messages, nil
}

// AddReceiver adds new receiver for one message
func (s *Storage) AddReceiver(userID int64, queueName string, messageIndex int64, clientID int64) error {
	message, err := s.Message(userID, queueName, messageIndex)
	if err != nil {
		return err
	}

	message.AddReceiver(clientID)
	return nil
}

// Receivers returns receivers of one message
func (s *Storage) Receivers(userID int64, queueName string, messageIndex int64) ([]int64, error) {
	message, err := s.Message(userID, queueName, messageIndex)
	if err != nil {
		return nil, err
	}

	return message.Recievers, nil
}

// DefaultStorage default storage
var DefaultStorage = Storage{
	queueMap: map[int64]map[string]*models.Queue{
		external.Root.ID(): map[string]*models.Queue{},
	},
}
