package storage

import (
	"encoding/json"
	"fmt"

	"sort"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

const messageKeyPrefix = "sqs.message"

func messageListKey(userID int64, queueName string, gorupID int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, gorupID)
}

func messageKey(userID int64, queueName string, index int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, index)
}

// Message for message storage
type Message struct {
	store *Storage
	db    KV
}

// All returns all messages index list
func (s *Message) All(userID int64, queueName string, gorupID int64, filters ...interface{}) ([]int64, error) {
	key := messageListKey(userID, queueName, gorupID)
	all, ok := s.db.Get(key)
	if !ok {
		all = "[]"
	}

	list := []int64{}
	if err := json.Unmarshal([]byte(all), &list); err != nil {
		return nil, errors.DataBroken(key)
	}

	return list, nil
}

// One returns a message string
func (s *Message) One(userID int64, queueName string, index int64) (string, bool) {
	key := messageKey(userID, queueName, index)
	return s.db.Get(key)
}

// Add adds a message
func (s *Message) Add(m *models.Message) error {
	if _, ok := s.One(m.UserID, m.QueueName, m.Index); ok {
		return errors.DuplicateMessage
	}

	all, err := s.All(m.UserID, m.QueueName, m.GroupID())
	if err != nil {
		return err
	}

	less := func(i, j int) bool { return all[i] < all[j] }
	if !sort.SliceIsSorted(all, less) {
		sort.Slice(all, less)
	}

	if len(all) > 0 && all[len(all)-1] > m.Index {
		return errors.MessageOutOfData
	}

	all = append(all, m.Index)
	data, err := json.Marshal(all)
	if err != nil {
		return errors.FailedEncoding(all)
	}

	err = s.db.Put(messageListKey(m.UserID, m.QueueName, m.GroupID()), string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(messageKey(m.UserID, m.QueueName, m.Index), m.Content)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}
