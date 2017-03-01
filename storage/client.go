package storage

import (
	"fmt"

	"strconv"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

const clientKeyPerfix = "sqs.client"

func clientKey(userID, clientID int64, queueName string) string {
	return fmt.Sprintf("%s.%d.%d.%s", clientKeyPerfix, userID, clientID, queueName)
}

// Client for storage of clients
type Client struct {
	store *Storage
	db    KV
}

// One returns a client
func (s *Client) One(userID int64, clientID int64, queueName string) (int64, error) {
	key := clientKey(userID, clientID, queueName)
	index, ok := s.db.Get(key)
	if !ok {
		return 0, errors.ClientNotFound
	}

	i, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		return 0, errors.DataBroken(key)
	}

	return i, nil
}

// Add adds client
func (s *Client) Add(c *models.Client) error {
	_, err := s.One(c.UserID, c.ID, c.QueueName)
	if err != errors.ClientNotFound {
		return err
	}

	key := clientKey(c.UserID, c.ID, c.QueueName)
	err = s.db.Put(key, fmt.Sprintf("%d", c.RecentMessageIndex))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}

// Update updates the RecentMessageIndex for the client
func (s *Client) Update(c *models.Client) error {
	recentMessageIndex, err := s.One(c.UserID, c.ID, c.QueueName)
	if err != errors.ClientNotFound {
		return err
	}

	if c.RecentMessageIndex == recentMessageIndex {
		return nil
	}

	key := messageKey(c.UserID, c.QueueName, c.ID)
	err = s.db.Put(key, fmt.Sprintf("%d", c.RecentMessageIndex))
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	return nil
}
