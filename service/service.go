package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/models"
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

	return s.Cache.AddConsumer(c, config.Config().ClientDefaultPriority)
}

func (s *Service) startPushMessage() {
	consumerChan := s.Cache.PopConsumer()
	for i := 0; i < config.Config().MaxPushWorkCount; i++ {
		go s.pushMessage(consumerChan)
	}
}

func (s *Service) pushMessage(ch <-chan *models.Consumer) {
	for {
		consumer := <-ch
		fmt.Println("START PUSHMESSAGE")
		fmt.Printf("CONSUMER: %#v\n", consumer)
		client := consumer.Client

		// remove consumer if out of control
		if consumer.Publisher != s.Agent.Address {
			continue
		}

		fmt.Printf("%#v\n", consumer.Client)
		message, err := s.Message.Next(consumer.UserID, consumer.QueueName, consumer.RecentMessageIndex)
		if err != nil {
			// log error to fix
			log.Println(err)
			s.Cache.PushConsumer(consumer)
			continue
		}

		// all messages has been pushed
		if message == nil {
			client.RecentPushedAt = time.Now().Unix()
			if err := s.Client.Update(client); err != nil {
				log.Println(err)
			}

			s.Cache.PushConsumer(consumer)
			continue
		}

		// push it the all client
		pushed := s.PushMessage(client.Addresses, message.Content)
		select {
		case <-time.After(time.Second * 5):
			consumer.Priority -= 2

		case <-pushed:
			consumer.Priority--
		}

		client.RecentPushedAt = time.Now().Unix()
		if err := s.Client.Update(client); err != nil {
			log.Println(err)
		}

		s.Cache.PushConsumer(consumer)
	}
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.Agent = agent.New(defaultService, address)

	defaultService.startPushMessage()
	log.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
