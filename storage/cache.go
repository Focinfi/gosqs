package storage

import (
	"container/heap"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

// Cache is for temporary data storage
type Cache struct {
	store     *Storage
	consumers heap.Interface
}

// AddConsumer add or refresh client ready to be pushed message
func (cache *Cache) AddConsumer(c *models.Client, priority int) error {
	if cache.consumers.Len() >= config.Config().MaxConsumerSize {
		return errors.ServiceOverload
	}

	consumer := models.NewConsumer(cache.consumers, c, priority)
	return cache.PushConsumer(consumer)
}

// PopConsumer returns a client ready to be pushed message
func (cache *Cache) PopConsumer() <-chan *models.Consumer {
	ch := make(chan *models.Consumer)
	go func() {
		for cache.consumers.Len() > 0 {
			ch <- heap.Pop(cache.consumers).(*models.Consumer)
		}
	}()

	return ch
}

// PushConsumer push consumer into cache
func (cache *Cache) PushConsumer(c *models.Consumer) error {
	heap.Push(cache.consumers, c)
	return nil
}

// NewCache returns a new cache
func NewCache(s *Storage) *Cache {
	return &Cache{
		store:     s,
		consumers: &models.PriorityConsumer{},
	}
}
