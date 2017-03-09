package storage

import (
	"container/heap"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
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
	log.Biz.Printf("AddConsumer: %#v", consumer.Client)
	return cache.PushConsumer(consumer, 0)
}

// PopConsumer returns a client ready to be pushed message
func (cache *Cache) PopConsumer() <-chan *models.Consumer {
	ch := make(chan *models.Consumer)
	go func() {
		log.Biz.Println("SETUP POPCONSUMER")
		for {
			if cache.consumers.Len() > 0 {
				log.Biz.Printf("POP_CONSUMER: %v\n", cache.consumers.Len())
				c := heap.Pop(cache.consumers).(*models.Consumer)
				log.Biz.Printf("POP_CONSUMER: Pushed=%v\n", c)
				ch <- c
			}
		}
	}()

	return ch
}

// PushConsumer push consumer into cache
func (cache *Cache) PushConsumer(c *models.Consumer, after time.Duration) error {
	log.Biz.Printf("CONSUMER-PUSH-BEFORE-LEN: %d\n", cache.consumers.Len())
	if after > 0 {
		time.AfterFunc(after, func() {
			heap.Push(cache.consumers, c)
			log.Biz.Printf("CONSUMER-PUSH-AFTER-LEN-After: %d\n", cache.consumers.Len())
		})
	} else {
		heap.Push(cache.consumers, c)
		log.Biz.Printf("CONSUMER-PUSH-AFTER-LEN: %d\n", cache.consumers.Len())
	}
	return nil
}

// NewCache returns a new cache
func NewCache(s *Storage) *Cache {
	return &Cache{
		store:     s,
		consumers: &models.PriorityConsumer{},
	}
}
