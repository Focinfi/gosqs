package storage

import (
	"math/rand"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/redis"
)

// Cache is for temporary data storage
type Cache struct {
	store   *Storage
	pl      models.PriorityList
	watcher models.Watcher
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

			after := time.Now().Unix() - c.Client().RecentReceivedAt
			log.Biz.Printf("AFTER: %#v, REC_AT: %#v\n", after, c.Client().RecentReceivedAt)

			// retry in a random time
			if after < int64(config.Config().MaxRetryConsumerSeconds) {
				rand.Seed(time.Now().UnixNano())
				waitTime := time.Millisecond * 100 * time.Duration(rand.Float64()*float64(after))
				time.AfterFunc(waitTime, func() {
					ch <- c
				})

				continue
			}

			// watch the changed
			key := models.QueueRecentMessageIDKey(c.Client().UserID, c.Client().QueueName)
			watchChan := cache.watcher.Watch(key)
			timeout := time.Second * time.Duration(config.Config().MaxRetryConsumerSeconds*2)
			go func() {
				select {
				case <-time.After(timeout):
				case change := <-watchChan:
					log.Biz.Infoln("NOTIFICATION: ", change)
					if change == "" {
						log.DB.Infof("watcher faild for key: %s\n", key)
					}
				}

				ch <- c
				log.Biz.Infoln("Pop Consumer")
			}()

		}
	}()

	return ch
}

// AddConsumer adds c into pl
func (cache *Cache) AddConsumer(c models.Consumer) error {
	return cache.pl.Add(c)
}

// PushConsumer pushes the c into cache
func (cache *Cache) PushConsumer(c models.Consumer) error {
	return cache.pl.Push(c)
}

// RemoveConsumer removes consumer
func (cache *Cache) RemoveConsumer(c models.Consumer) error {
	return cache.pl.Remove(c)
}

// NewConsumer returns a new Consumer based on the client
func NewConsumer(client *models.Client, priority int) models.Consumer {
	return redis.NewConsumer(client, priority)
}
