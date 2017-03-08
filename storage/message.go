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
		return nil, errors.DataBroken(key, err)
	}

	return list, nil
}

// One returns a message string
func (s *Message) One(userID int64, queueName string, index int64) (string, bool) {
	key := messageKey(userID, queueName, index)
	return s.db.Get(key)
}

// Next for next message of current Message
func (s *Message) Next(userID int64, queueName string, index int64, timestamp int64) (*models.Message, error) {
	orginID := models.GroupID(index)
	groupID := orginID
	nowGroupID := timestamp
	fmt.Printf("GROUP_ID: %d, NOW_GORUP_ID:%d\n", groupID, nowGroupID)

	for {
		// the message with
		if nowGroupID <= groupID {
			break
		}
		all, err := s.All(userID, queueName, groupID, nil)
		if err != nil {
			return nil, err
		}

		// if this group id empty try the next
		if len(all) == 0 {
			groupID++
			continue
		}
		fmt.Printf("MESSAGE-ALL: %v\n", all)
		var nextIdx int64
		if index/models.BaseUnit == 0 || orginID < groupID {
			nextIdx = all[0]
		} else {
			i := sort.Search(len(all), func(i int) bool { return all[i] >= index })
			if i < len(all) && all[i] == index {
				// last one in the group with the id groupID
				if i == len(all)-1 {
					// try next group
					groupID++
					continue
				}

				// got the next message index
				nextIdx = all[i+1]
			} else {
				return nil, errors.DataLost(messageKey(userID, queueName, index))
			}
		}

		message, ok := s.One(userID, queueName, nextIdx)
		if !ok {
			return nil, errors.DataLost(messageKey(userID, queueName, nextIdx))
		}

		fmt.Printf("NEXT INDEX: %d\n", nextIdx)
		return &models.Message{
			UserID:    userID,
			QueueName: queueName,
			Index:     nextIdx,
			Content:   message,
		}, nil
	}

	return nil, nil
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
