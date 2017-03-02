package service

import (
	"log"
	"net/http"
	"time"

	"github.com/Focinfi/sqs/agent"
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
	if err := s.applyClient(c); err != nil {
		return err
	}

	if err := s.database.RegisterClient(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) applyClient(c *models.Client) error {
	if err := s.Cache.AddClient(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) startPushMessage() {
	for {
		client := <-s.Cache.Client()
		message, err := s.Message.Next(client.UserID, client.QueueName, client.RecentMessageIndex)
		if err != nil {
			// log error to fix
			log.Println(err)
			continue
		}

		// all messages has been pushed
		if message == nil {
			client.RecentReceivedAt = time.Now().Unix()
			if err := s.Client.Update(client); err != nil {
				log.Println(err)
				continue
			}
		}

		// push it the all client
		pushed := s.PushMessage(client.Addresses, message.Content)
		go func() {
			select {
			case <-time.After(time.Second * 5):
				client.RecentReceivedAt = time.Now().Unix()
				if err := s.Client.Update(client); err != nil {
					log.Println(err)
				}
			case <-pushed:
				client.RecentMessageIndex = message.Index
				client.RecentReceivedAt = time.Now().Unix()
				if err := s.Client.Update(client); err != nil {
					log.Println(err)
				}
			}
		}()
	}
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.Agent = agent.New(defaultService, address)

	go defaultService.startPushMessage()
	log.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
