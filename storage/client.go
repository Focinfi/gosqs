package storage

import (
	"encoding/json"
	"fmt"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

// Client for storage of clients
type Client struct {
	store *Storage
	db    KV
}

// One returns a client
func (s *Client) One(userID int64, clientID int64, queueName string) (*models.Client, error) {
	key := models.ClientKey(userID, clientID, queueName)
	value, ok := s.db.Get(key)
	if !ok {
		return nil, errors.ClientNotFound
	}

	client := &models.Client{}
	err := json.Unmarshal([]byte(value), client)
	if err != nil {
		return nil, errors.DataBroken(key, err)
	}

	return client, nil
}

// Add adds client
func (s *Client) Add(c *models.Client) error {
	client, err := s.One(c.UserID, c.ID, c.QueueName)
	if err != errors.ClientNotFound {
		return err
	}
	if client != nil {
		return errors.DuplicateClient
	}

	key := models.ClientKey(c.UserID, c.ID, c.QueueName)
	data, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(key, string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// Update updates the RecentMessageIndex for the client
func (s *Client) Update(c *models.Client) error {
	fmt.Printf("TO UPDATE CLIENT: %#v\n", c)
	_, err := s.One(c.UserID, c.ID, c.QueueName)
	if err != nil {
		return err
	}

	key := models.ClientKey(c.UserID, c.ID, c.QueueName)
	data, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	err = s.db.Put(key, string(data))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}
