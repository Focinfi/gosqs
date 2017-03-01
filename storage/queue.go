package storage

import (
	"encoding/json"
	"fmt"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

const queueKeyPrefix = "sqs.queue."

func queueListKey(userID int64) string {
	return fmt.Sprintf("%s%d", queueKeyPrefix, userID)
}

func queueKey(userID int64, queueName string) string {
	return fmt.Sprintf("%s%d.%s", queueKeyPrefix, userID, queueName)
}

// Queue stores data
type Queue struct {
	db KV
}

// All returns queue map for userID
func (s *Queue) All(userID int64) ([]models.Queue, error) {
	all := []models.Queue{}
	key := queueListKey(userID)

	val, ok := s.db.Get(key)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	if val == "" {
		return all, nil
	}

	if err := json.Unmarshal([]byte(val), &all); err != nil {
		return nil, errors.DataBroken(key)
	}

	return all, nil
}

// One returns a queue for the userID with the queueName
func (s *Queue) One(userID int64, queueName string) (*models.Queue, error) {
	all, err := s.All(userID)
	if err != nil {
		return nil, err
	}

	var theQueue *models.Queue
	for _, queue := range all {
		if queue.Name == queueName {
			*theQueue = queue
		}
	}

	if theQueue == nil {
		return nil, errors.QueueNotFound
	}

	return theQueue, nil
}

// Add add q for userID
func (s *Queue) Add(q *models.Queue) error {
	all, err := s.All(q.UserID)
	if err != nil {
		return err
	}

	for _, queue := range all {
		if queue.Name == q.Name {
			return errors.DuplicateQueue
		}
	}

	all = append(all, *q)
	data, err := json.Marshal(all)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(queueListKey(q.UserID), string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// Remove removes the queue with the name
func (s *Queue) Remove(userID int64, queueName string) error {
	all, err := s.All(userID)
	if err != nil {
		return err
	}

	index := -1
	for i, queue := range all {
		if queue.Name == queueName {
			index = i
		}
	}

	if index < 0 {
		return errors.QueueNotFound
	}

	all = append(all[:index], all[index+1:]...)
	data, err := json.Marshal(all)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(queueListKey(userID), string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// DefaultQueue default
var DefaultQueue = &Queue{
	db: defaultKV,
}
