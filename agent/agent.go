package agent

import (
	"net/http"

	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
)

// QueueService defines what a queue admin should do
type QueueService interface {
	ReceiveMessage(userID int64, queueName, content string, index int64) error
	RegisterClient(client *models.Client) error
	ApplyMessageIDRange(userID int64, queueName string, size int) (maxID int64, err error)
}

// Agent for receiving message and push them
type Agent struct {
	Address string
	http.Handler
	QueueService
}

// New allocates and returns a new Agent
func New(admin QueueService, address string) *Agent {
	agent := &Agent{
		Address:      address,
		QueueService: admin,
	}
	agent.routing()
	return agent
}

func (agent *Agent) routing() {
	s := gin.Default()
	group := s.Group("/")
	group.POST("/message", agent.ReceiveMessage)
	group.POST("/register", agent.RegisterClient)
	group.PUT("/messageID", agent.ApplyMessageIDRange)
	agent.Handler = s
}
