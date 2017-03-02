package storage

import "github.com/Focinfi/sqs/models"

// Cache is for temporary data storage
type Cache struct {
	store *Storage
	KV
}

// AddClient add or refresh client ready to be pushed message
func (cache *Cache) AddClient(c *models.Client) error {
	return nil
}

// Client returns a client ready to be pushed message
func (cache *Cache) Client() <-chan *models.Client {
	return make(chan *models.Client)
}
