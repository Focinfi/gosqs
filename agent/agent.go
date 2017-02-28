package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessagePusher push message
type MessagePusher interface {
	PushMessage(userID int64, name, content string) error
}

// Agent for receiving message and push them
type Agent struct {
	MessagePusher
	http.Handler
}

// New allocates and returns a new Agent
func New(mp MessagePusher) *Agent {
	s := gin.Default()
	agent := &Agent{Handler: s, MessagePusher: mp}

	group := s.Group("/", throttling, parsing, auth)
	group.POST("/message", agent.ReceiveMessage)

	return agent
}

// ReceiveMessage serve message pushing via http
func (agent *Agent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		UserID    int64  `json:"userID"`
		QueueName string `json:"queueName"`
		Content   string `json:"content"`
	}

	params := &messageParam{}
	if err := ctx.BindJSON(params); err != nil {
		responseAndAbort(ctx, StatusWrongParam)
		return
	}

	err := agent.MessagePusher.PushMessage(params.UserID, params.QueueName, params.Content)
	if err != nil {
		responseAndAbort(ctx, StatusIsBusy(err))
		return
	}

	responseOK(ctx)
}

// StartDeliveryMessage deliveries messages to all online subsribers
func (agent *Agent) StartDeliveryMessage() {

}
