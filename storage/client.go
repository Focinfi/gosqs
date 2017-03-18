package storage

import (
	"encoding/json"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
)

// Client for storage of clients
type Client struct {
	store *Storage
	db    models.KV
}

// One returns a client
func (s *Client) One(userID int64, clientID int64, queueName string) (*models.Client, error) {
	key := models.ClientKey(userID, clientID, queueName)
	value, err := s.db.Get(key)
	if err == errors.DBNotFound {
		return nil, errors.ClientNotFound
	}

	client := &models.Client{}
	err = json.Unmarshal([]byte(value), client)
	if err != nil {
		return nil, errors.DataBroken(key, err)
	}

	return client, nil
}

// Add adds or update client
func (s *Client) Add(c *models.Client) error {
	key := models.ClientKey(c.UserID, c.ID, c.QueueName)
	data, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErrorf(err.Error())
	}

	err = s.db.Put(key, string(data))
	if err != nil {
		return errors.NewInternalErrorf(err.Error())
	}

	return nil
}

// Update updates the RecentMessageIndex for the client
func (s *Client) Update(c *models.Client) error {
	log.Biz.Printf("TO UPDATE CLIENT: %#v\n", c)
	_, err := s.One(c.UserID, c.ID, c.QueueName)
	if err != nil {
		return err
	}

	key := models.ClientKey(c.UserID, c.ID, c.QueueName)
	data, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErrorf(err.Error())
	}

	err = s.db.Put(key, string(data))
	if err != nil {
		return errors.NewInternalErrorf(err.Error())
	}

	return nil
}
