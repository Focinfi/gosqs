package agent

import "github.com/Focinfi/sqs/models"

// MasterService can distributes a node for a consume
type MasterService interface {
	AssignNode(userID int64, queueName string, squadName string) (string, error)
	Join(info models.NodeInfo)
}

// QueueService defines what a queue admin should do
type QueueService interface {
	ApplyMessageIDRange(userID int64, queueName string, size int) (maxID int64, err error)
	PushMessage(userID int64, queueName, content string, index int64) error
	PullMessage(userID int64, queueName, squadName string, length int) ([]models.Message, error)
	ReportMaxReceivedMessageID(userID int64, queueName, squadName string, messageID int64) error
	Info() models.NodeInfo
}
