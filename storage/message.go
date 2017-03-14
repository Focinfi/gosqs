package storage

import (
	"encoding/json"

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
	all, ok := s.db.Get(key)
	if !ok {
		all = "[]"
	}

	list := []int64{}
	if err := json.Unmarshal([]byte(all), &list); err != nil {
		return nil, errors.DataBroken(key, err)
	}

	return list, nil
}

// One returns a message string
func (s *Message) One(userID int64, queueName string, index int64) (string, bool) {
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
		mVal, ok := s.One(userID, queueName, nextIdx)
		log.Biz.Infoln("mVal", mVal, ok)

		if ok {
			if mVal == "" {
				return nil, errors.DataLost(models.MessageKey(userID, queueName, index))
			}
			message = mVal
			break
		}

		nextIdx++
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
	if _, ok := s.One(m.UserID, m.QueueName, m.Index); ok {
		return errors.DuplicateMessage
	}

	err := s.db.Put(models.MessageKey(m.UserID, m.QueueName, m.Index), m.Content)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}
