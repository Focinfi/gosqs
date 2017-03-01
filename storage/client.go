package storage

import "github.com/Focinfi/sqs/models"

// Client for storage of clients
type Client struct {
	db KV
}

// One returns a client
func (client *Client) One(userID int64, queueName string, messageIndex int64) (*models.Client, error) {
	return nil, nil
}

// All returns clients
func (client *Client) All(userID int64, queueName string, filter interface{}) (map[int64]*models.Client, error) {
	return nil, nil
}

// Add adds client
func (client *Client) Add(c *models.Client) error {
	return nil
}

// Update updates the RecentMessageIndex for the client
func (client *Client) Update(recentMessageIndex int64) error {
	return nil
}

// DefaultClient default client
var DefaultClient = &Client{db: defaultKV}
