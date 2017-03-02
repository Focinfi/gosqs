package agent

import "github.com/gin-gonic/gin"

// ReceiveMessage serve message pushing via http
func (agent *Agent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		UserID    int64  `json:"userID"`
		QueueName string `json:"queueName"`
		Content   string `json:"content"`
		Index     int64  `josn:"index"`
	}

	params := &messageParam{}
	if err := ctx.BindJSON(params); err != nil {
		response(ctx, err)
		return
	}

	err := agent.PushMessage(params.UserID, params.QueueName, params.Content, params.Index)
	response(ctx, err)
}

// Register registers so can get the message
func (agent *Agent) Register(ctx *gin.Context) {
	type registerParam struct {
		UserID    int64  `json:"userID"`
		ClientID  int64  `json:"clientID"`
		QueueName string `json:"queueName"`
	}

	param := &registerParam{}
	if err := ctx.BindJSON(param); err != nil {
		response(ctx, err)
		return
	}

	err := agent.RegisterClient(param.UserID, param.ClientID, param.QueueName)
	response(ctx, err)
}
