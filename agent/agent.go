package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessagePusher push message
type MessagePusher interface {
	PushMessage(userID int64, name, content string) error
}

// ClientRegistrar register clients
type ClientRegistrar interface {
	RegisterClient(userID int64, clientID int64, queueName string) error
}

// QueueAdmin defines what a queue admin should do
type QueueAdmin interface {
	MessagePusher
	ClientRegistrar
}

// Agent for receiving message and push them
type Agent struct {
	http.Handler
	QueueAdmin
}

// New allocates and returns a new Agent
func New(admin QueueAdmin) *Agent {
	agent := &Agent{
		QueueAdmin: admin,
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
