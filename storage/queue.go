package storage

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
)

// Queue stores data
type Queue struct {
	store *Storage
	db    models.KV
	inc   models.Incrementer
}

// All returns queue map for userID
func (s *Queue) All(userID int64) ([]models.Queue, error) {
	all := []models.Queue{}
	key := models.QueueListKey(userID)

	val, ok := s.db.Get(key)
	if !ok {
		return nil, errors.UserNotFound
	}

	if val == "" {
		return all, nil
	}

	if err := json.Unmarshal([]byte(val), &all); err != nil {
		return nil, errors.DataBroken(key, err)
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

	newAll := append(all, *q)
	data, err := json.Marshal(newAll)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(models.QueueListKey(q.UserID), string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// UpdateRecentMessageGroupID updates the groupID of one queue
func (s *Queue) UpdateRecentMessageGroupID(userID int64, queueName string, groupID int64) error {
	k := models.QueueKey(userID, queueName)
	groupIDVal := fmt.Sprintf("%d", groupID)

	oldVal, ok := s.db.Get(k)
	if !ok {
		return s.db.Put(k, groupIDVal)
	}

	id, err := strconv.ParseInt(oldVal, 10, 64)
	if err != nil {
		log.DB.Error(errors.DataBroken(k, err))
		return s.db.Put(k, groupIDVal)
	}

	if groupID > id {
		return s.db.Put(k, groupIDVal)
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

	err = s.db.Put(models.QueueListKey(userID), string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// ApplyMessageIDRange try to apply message id range
func (s *Queue) ApplyMessageIDRange(userID int64, queueName string, size int) (int64, error) {
	key := models.QueueMaxIDKey(userID, queueName)
	return s.inc.Increment(key, size)
}
