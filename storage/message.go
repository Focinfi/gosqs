package storage

import (
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

// One returns a message string
func (s *Message) One(userID int64, queueName string, index int64) (string, error) {
	key := models.MessageKey(userID, queueName, index)

	return s.db.Get(key)
}

// Nextn returns the next n messages after the given currentID, try MaxTryMessageCount entries after maxMassageID
func (s *Message) Nextn(userID int64, queueName string, currentID int64, maxMassageID int64, n int) ([]models.Message, error) {
	nextIdx := currentID + 1
	upperID := maxMassageID + int64(config.Config.MaxTryMessageCount)
	log.Internal.Infoln("Nextn", nextIdx, upperID)
	messages := []models.Message{}

	for nextIdx <= upperID {
		if len(messages) >= n {
			return messages, nil
		}

		msgContent, err := s.One(userID, queueName, nextIdx)
		log.DB.Infoln("Nextn", nextIdx, err, msgContent)
		if err == errors.DataNotFound {
			nextIdx++
			continue
		}

		if err != nil {
			return nil, err
		}

		if msgContent == "" {
			return nil, errors.DataLost(models.MessageKey(userID, queueName, nextIdx))
		}

		message := models.Message{
			UserID:    userID,
			QueueName: queueName,
			Content:   msgContent,
			Index:     nextIdx,
		}

		messages = append(messages, message)
		nextIdx++
	}

	return messages, nil
}

// Add adds a message
func (s *Message) Add(m *models.Message) error {
	_, getErr := s.One(m.UserID, m.QueueName, m.Index)
	log.DB.Infof("Get(%d.%s.%d) Error: %v, time: %v\n", m.UserID, m.QueueName, m.Index, getErr, time.Now())

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
