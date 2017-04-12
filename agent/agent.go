package agent

import (
	"net/http"

	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
)

// MasterService can distributes a node for a consume
type MasterService interface {
	AssignNode(userID int64, queueName string, squadName string) (string, error)
	//Heartbeat()
}

// QueueService defines what a queue admin should do
type QueueService interface {
	ApplyMessageIDRange(userID int64, queueName string, size int) (maxID int64, err error)
	PushMessage(userID int64, queueName, content string, index int64) error
	PullMessage(userID int64, queueName, squadName string, length int) ([]models.Message, error)
	ReportMaxReceivedMessageID(userID int64, queueName, squadName string, messageID int64) error
}

type MasterAgent struct {
	Address string
	http.Handler
	MasterService
}

func NewMasterAgent(service MasterService, address string) *MasterAgent {
	agt := &MasterAgent{
		Address:       address,
		MasterService: service,
	}

	agt.masterAgentRouting()
	return agt
}

// QueueAgent x
type QueueAgent struct {
	Address string
	http.Handler
	QueueService
}

func (a *MasterAgent) masterAgentRouting() {
	s := gin.Default()
	group := s.Group("/")
	group.POST("/applyNode", a.ApplyNode)
	a.Handler = s
}

// NewQueueAgent allocates and returns a new QueueAgent
func NewQueueAgent(service QueueService, address string) *QueueAgent {
	agt := &QueueAgent{
		Address:      address,
		QueueService: service,
	}
	agt.queueAgentRouting()
	return agt
}

func (a *QueueAgent) queueAgentRouting() {
	s := gin.Default()
	group := s.Group("/")
	group.GET("/messages", a.PullMessages)
	group.POST("/message", a.ReceiveMessage)
	group.PUT("/messageID", a.ApplyMessageIDRange)
	group.PUT("/receivedMessageID", a.ReportMaxReceivedMessageID)
	a.Handler = s
}
