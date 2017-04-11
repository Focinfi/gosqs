package node

import (
	"net/http"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
)

var (
	defaultPullMessageCount = config.Config().PullMessageCount
)

// Service for one user info
type Service struct {
	masterAddr string
	*database
	*agent.QueueAgent
}

func (s *Service) PullMessage(userID int64, queueName, squadName string, length int) ([]models.Message, error) {
	squad, err := s.Squad.One(userID, queueName, squadName)
	if err != nil {
		return nil, err
	}

	// TODO: to confirm the concurrent requests result
	return s.Message.Nextn(userID, queueName, squad.ReceivedMessageID, defaultPullMessageCount)
}

func (s *Service) ReportMaxReceivedMessageID(userID int64, queueName, squadName string, messageID int64) error {
	squad, err := s.database.Squad.One(userID, queueName, squadName)
	if err != nil {
		return err
	}

	if squad.ReceivedMessageID >= messageID {
		return nil
	}

	squad.ReceivedMessageID = messageID
	return s.Squad.Update(squad)
}

// PushMessage receives message
func (s *Service) PushMessage(userID int64, queueName, content string, index int64) error {
	maxID, err := s.Queue.MessageMaxID(userID, queueName)
	if err != nil {
		return err
	}

	if index > maxID {
		return errors.MessageIndexOutOfRange
	}

	return s.database.PushMessage(userID, queueName, content, index)
}

// ApplyMessageIDRange tries to apply a range a free message id
func (s *Service) ApplyMessageIDRange(userID int64, queueName string, size int) (maxID int64, err error) {
	if size > config.Config().MaxMessgeIDRangeSize {
		return -1, errors.ApplyMessageIDRangeOversize
	}

	return s.Queue.ApplyMessageIDRange(userID, queueName, size)
}

// Start starts services
func Start(address string) {
	var defaultService = &Service{database: db}
	defaultService.QueueAgent = agent.NewQueueAgent(defaultService, address)

	log.Biz.Fatal(http.ListenAndServe(address, defaultService.QueueAgent))
}
