package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// QueueService defines what a queue admin should do
type QueueService interface {
	PushMessage(userID int64, name, content string) error
	RegisterClient(userID int64, clientID int64, queueName string) error
}

// Agent for receiving message and push them
type Agent struct {
	http.Handler
	QueueService
}

// New allocates and returns a new Agent
func New(admin QueueService) *Agent {
	agent := &Agent{
		QueueService: admin,
	}
	agent.routing()
	return agent
}

func (agent *Agent) routing() {
	s := gin.Default()
	group := s.Group("/")
	group.POST("/message", agent.ReceiveMessage)
	group.POST("/register", agent.Register)
	agent.Handler = s
}
