package storage

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/redis"
)

// Cache is for temporary data storage
type Cache struct {
	store *Storage
	pl    models.PriorityList
}

// PopConsumerChan returns a output Consumer channel
func (cache *Cache) PopConsumerChan() <-chan models.Consumer {
	ch := make(chan models.Consumer, config.Config().MaxPushWorkCount)
	go func() {
		log.Biz.Println("SETUP POPCONSUMER")
		for {
			c, err := cache.pl.Pop()
			if err == errors.NoConsumer {
				continue
			}

			if err != nil {
				log.DB.Error(err)
				continue
			}

			after := float64(time.Now().Unix() - c.Client().RecentReceivedAt)
			fmt.Printf("AFTER: %#v, REC_AT: %#v\n", after, c.Client().RecentReceivedAt)
			time.AfterFunc(time.Millisecond*100*time.Duration(rand.Float64()*float64(after)), func() {
				ch <- c
			})

		}
	}()

	return ch
}

// PushConsumer pushes the c into cache
func (cache *Cache) PushConsumer(c models.Consumer) error {
	return cache.pl.Push(c)
}

// NewCache returns a new cache
func NewCache(s *Storage) *Cache {
	var pl models.PriorityList

	// goheap
	// pl = goheap.New()

	pl, err := redis.New()
	if err != nil {
		log.DB.Panic(err)
	}

	return &Cache{
		store: s,
		pl:    pl,
	}
}

// NewConsumer returns a new Consumer based on the client
func NewConsumer(client *models.Client, priority int) models.Consumer {
	// return goheap.NewConsumer(client, priority)
	return redis.NewConsumer(client, priority)
}
