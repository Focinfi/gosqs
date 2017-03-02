package service

import (
	"log"
	"net/http"

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
	// check availability
	// s.Availability()

	// add into list
	if err := s.Cache.AddClient(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) startPushMessage() error {
	go func() {
		for {
			// client := <-s.Cache.Client()
			// s.Agent.PushMessage(client.Addresses, message)
		}
	}()

	return nil
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.Agent = agent.New(defaultService, address)

	if err := defaultService.startPushMessage(); err != nil {
		panic(err)
	}
	log.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
