package storage

import "github.com/Focinfi/sqs/models"

// Client for storage of clients
type Client struct {
	data interface{}
}

// Add adds client
func (client *Client) Add(c *models.Client) error {
	return nil
}

// DefaultClient default client
var DefaultClient = &Client{}
