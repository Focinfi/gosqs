package agent

import (
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ReceiveMessage serve message pushing via http
func (a *Agent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		UserID    int64  `json:"userID"`
		QueueName string `json:"queueName"`
		Content   string `json:"content"`
		Index     int64  `josn:"index"`
	}

	params := &messageParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		response(ctx, err)
		return
	}

	err := a.ReceivehMessage(params.UserID, params.QueueName, params.Content, params.Index)
	response(ctx, err)
}

// RegisterClient registers so can get the message
func (a *Agent) RegisterClient(ctx *gin.Context) {
	type registerParam struct {
		UserID    int64  `json:"userID"`
		ClientID  int64  `json:"clientID"`
		QueueName string `json:"queueName"`
		Address   string `json:"address"`
	}

	params := &registerParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		response(ctx, err)
		return
	}

	client := &models.Client{
		UserID:    params.UserID,
		ID:        params.ClientID,
		QueueName: params.QueueName,
		Publisher: a.Address,
		Address:   params.Address,
	}

	err := a.QueueService.RegisterClient(client)
	response(ctx, err)
}
