package storage

import (
	"time"

	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/goheap"
)

// Cache is for temporary data storage
type Cache struct {
	store *Storage
	pl    models.PriorityList
}

// PopConsumerChan returns a output Consumer channel
func (cache *Cache) PopConsumerChan() <-chan models.Consumer {
	ch := make(chan models.Consumer)
	go func() {
		log.Biz.Println("SETUP POPCONSUMER")
		for {
			c := cache.pl.Pop()
			if c != nil {
				ch <- c
			}
		}
	}()

	return ch
}

// PushConsumerAt push consumer into cache
func (cache *Cache) PushConsumerAt(c models.Consumer, after time.Duration) error {
	if after > 0 {
		time.AfterFunc(after, func() {
			cache.pl.Push(c)
		})
	} else {
		cache.pl.Push(c)
	}
	return nil
}

// NewCache returns a new cache
func NewCache(s *Storage) *Cache {
	return &Cache{
		store: s,
		pl:    goheap.New(),
	}
}

// NewConsumer returns a new Consumer based on the client
func NewConsumer(client *models.Client, priority int) models.Consumer {
	return goheap.NewConsumer(client, priority)
}
