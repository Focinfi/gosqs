package service

import (
	"net/http"
	"time"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage"
)

// Service for one user info
type Service struct {
	*database
	*agent.Agent
}

// ReceivehMessage receives message
func (s *Service) ReceivehMessage(userID int64, queueName, content string, index int64) error {
	return s.database.ReceivehMessage(userID, queueName, content, index)
}

// RegisterClient registers client
func (s *Service) RegisterClient(c *models.Client) error {
	err := s.database.RegisterClient(c)

	if err != nil {
		return err
	}

	consumer := storage.NewConsumer(c, config.Config().ClientDefaultPriority)
	return s.Cache.AddConsumer(consumer)
}

func (s *Service) startPushMessage() {
	consumerChan := s.Cache.PopConsumerChan()
	for i := 0; i < config.Config().MaxPushWorkCount; i++ {
		go s.pushMessage(consumerChan)
	}
}

func (s *Service) pushMessage(ch <-chan models.Consumer) {
	for {
		consumer := <-ch
		log.Biz.Println("START PUSHMESSAGE")
		log.Biz.Printf("CONSUMER: %#v\n", consumer)
		now := time.Now().Unix()
		client := consumer.Client()

		c, err := s.Client.One(client.UserID, client.ID, client.QueueName)
		// client is removed, discard this consumer
		if err == errors.ClientNotFound {
			if err := s.Cache.RemoveConsumer(consumer); err != nil {
				log.DB.Error(err)
			}
			continue
		}

		if err != nil {
			log.DB.Error(err)
			s.Cache.PushConsumer(consumer)
		}

		// remove consumer if out of control
		if client.Publisher != c.Publisher {
			if err := s.Cache.RemoveConsumer(consumer); err != nil {
				log.DB.Error(err)
			}
			continue
		}

		// update client with c
		*client = *c

		message, err := s.Message.Next(client.UserID, client.QueueName, client.RecentMessageIndex, now)
		log.Biz.Printf("MESSAGE: %v, err: %v\n", message, err)
		if err != nil {
			log.DB.Errorln(err)
			consumer.IncPriority(-2)
			s.Cache.PushConsumer(consumer)
			continue
		}

		// all messages has been pushed
		if message == nil {
			client.RecentPushedAt = now
			if period := now - models.GroupID(client.RecentMessageIndex); period > 3 {
				client.RecentMessageIndex = models.GenIndex0(now - 3)
			}
			// failed to update client, should fix and give away the control of consumer
			if err := s.Client.Update(client); err != nil {
				log.DB.Errorln(err)
				continue
			}

			s.Cache.PushConsumer(consumer)
			continue
		}

		// push it the all client
		pushed := s.PushMessage(client.Addresses, message.Content)
		select {
		case <-time.After(time.Second * 5):
			consumer.IncPriority(-1)
		case <-pushed:
			log.Biz.Println("PUSHED")
			consumer.Client().RecentMessageIndex = message.Index
			consumer.IncPriority(1)
			client.RecentReceivedAt = now
		}

		client.RecentPushedAt = now
		// failed to update client, should fix and give away the control of consumer
		if err := s.Client.Update(client); err != nil {
			log.DB.Errorln(err)
			continue
		}

		s.Cache.PushConsumer(consumer)
	}
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.Agent = agent.New(defaultService, address)

	defaultService.startPushMessage()
	log.Biz.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
