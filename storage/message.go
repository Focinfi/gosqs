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
	if getErr == errors.DataNotFound {
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

func (s *Message) Nextn(userID int64, queueName string, currentID int64, n int) ([]models.Message, error) {
	nextIdx := currentID + 1
	upperID := nextIdx + int64(config.Config().MaxTryMessageCount)
	messages := make([]models.Message, 0, n)

	for nextIdx <= upperID {
		if len(messages) >= n {
			return messages, nil
		}

		msgContent, err := s.One(userID, queueName, nextIdx)
		log.Biz.Infof("message[%d]='%s', err: %v", currentID, msgContent, err)

		if err == errors.DataNotFound {
			nextIdx++
			continue
		}

		if err != nil {
			return nil, err
		}

		if msgContent == "" {
			return nil, errors.DataLost(models.MessageKey(userID, queueName, currentID))
		}

		message := models.Message{
			UserID:    userID,
			QueueName: queueName,
			Content:   msgContent,
			Index:     nextIdx,
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// Add adds a message
func (s *Message) Add(m *models.Message) error {
	_, getErr := s.One(m.UserID, m.QueueName, m.Index)
	log.Biz.Printf("Get Error: %v, time: %v\n", getErr, time.Now())

	if getErr == errors.DataNotFound {
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
