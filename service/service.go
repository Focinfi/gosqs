package service

import (
	"net/http"
	"time"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/config"
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
	if err := s.database.RegisterClient(c); err != nil {
		return err
	}

	consumer := storage.NewConsumer(c, config.Config().ClientDefaultPriority)
	return s.Cache.PushConsumerAt(consumer, 0)
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
		log.Biz.Printf("CONSUMER: %#v\n", consumer.Client())
		now := time.Now().Unix()
		client := consumer.Client()

		if c, err := s.Client.One(client.UserID, client.ID, client.QueueName); err != nil {
			log.DB.Error(err)
			s.Cache.PushConsumerAt(consumer, time.Millisecond)
		} else if client.Publisher != c.Publisher { // remove consumer if out of control
			continue
		}

		message, err := s.Message.Next(client.UserID, client.QueueName, client.RecentMessageIndex, now)
		log.Biz.Printf("MESSAGE: %v, err: %v\n", message, err)
		if err != nil {
			log.DB.Errorln(err)
			consumer.IncPriority(-10)
			s.Cache.PushConsumerAt(consumer, time.Second)
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

			consumer.IncPriority(-1)
			s.Cache.PushConsumerAt(consumer, time.Second)
			continue
		}

		// push it the all client
		pushed := s.PushMessage(client.Addresses, message.Content)
		after := time.Millisecond
		select {
		case <-time.After(time.Second * 5):
			consumer.IncPriority(-2)
			after = time.Second * time.Duration(5)
		case <-pushed:
			log.Biz.Println("PUSHED")
			consumer.Client().RecentMessageIndex = message.Index
			client.RecentPushedAt = time.Now().Unix()
		}

		// failed to update client, should fix and give away the control of consumer
		if err := s.Client.Update(client); err != nil {
			log.DB.Errorln(err)
			continue
		}

		s.Cache.PushConsumerAt(consumer, after)
	}
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.Agent = agent.New(defaultService, address)

	defaultService.startPushMessage()
	log.Biz.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
