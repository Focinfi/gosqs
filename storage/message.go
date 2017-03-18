package storage

import (
	"encoding/json"

	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
)

// Message for message storage
type Message struct {
	store *Storage
	db    models.KV
}

// All returns all messages index list
func (s *Message) All(userID int64, queueName string, gorupID int64, filters ...interface{}) ([]int64, error) {
	key := models.MessageListKey(userID, queueName, gorupID)
	all, getErr := s.db.Get(key)
	if getErr == errors.DBNotFound {
		all = "[]"
	}

	if getErr != nil {
		return nil, getErr
	}

	list := []int64{}
	if err := json.Unmarshal([]byte(all), &list); err != nil {
		return nil, errors.DataBroken(key, err)
	}

	return list, nil
}

// One returns a message string
func (s *Message) One(userID int64, queueName string, index int64) (string, error) {
	key := models.MessageKey(userID, queueName, index)

	return s.db.Get(key)
}

// Next for next message of current Message
func (s *Message) Next(userID int64, queueName string, index int64, maxIdx int64) (*models.Message, error) {
	log.Biz.Infoln("NEXT: ", userID, queueName, index, maxIdx)
	nextIdx := index + 1
	upperIdx := maxIdx + int64(config.Config().MaxTryMessageCount)
	var message string

	for nextIdx <= upperIdx {
		mVal, getErr := s.One(userID, queueName, nextIdx)
		log.Biz.Infoln("mVal", mVal, getErr)

		if getErr == errors.DBNotFound {
			nextIdx++
			continue
		}

		if getErr != nil {
			return nil, getErr
		}

		if mVal == "" {
			return nil, errors.DataLost(models.MessageKey(userID, queueName, index))
		}

		message = mVal
		break
	}

	// next is nil
	if nextIdx > upperIdx {
		return nil, nil
	}

	return &models.Message{
		UserID:    userID,
		QueueName: queueName,
		Index:     nextIdx,
		Content:   message,
	}, nil
}

// Add adds a message
func (s *Message) Add(m *models.Message) error {
	_, getErr := s.One(m.UserID, m.QueueName, m.Index)
	log.Biz.Printf("Get Error: %v, time: %v\n", getErr, time.Now())

	if getErr == errors.DBNotFound {
		err := s.db.Put(models.MessageKey(m.UserID, m.QueueName, m.Index), m.Content)
		if err != nil {
			return errors.NewInternalErrorf(err.Error())
		}

		return nil
	}

	if getErr != nil {
		return getErr
	}

	return errors.DuplicateMessage
}
