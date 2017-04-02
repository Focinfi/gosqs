package agent

import (
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// TODO: to delegate to service
func (a *Agent) getMessages(ctx *gin.Context) {}

// TODO: to delegate to service
func (a *Agent) reportMessageID(ctx *gin.Context) {}

// ReceiveMessage serve message pushing via http
func (a *Agent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		UserID    int64  `json:"user_id"`
		QueueName string `json:"queue_name"`
		Content   string `json:"content"`
		Index     int64  `josn:"index"`
	}

	params := &messageParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	err := a.QueueService.ReceiveMessage(params.UserID, params.QueueName, params.Content, params.Index)
	responseErr(ctx, err)
}

// RegisterClient registers so can get the message
func (a *Agent) RegisterClient(ctx *gin.Context) {
	type registerParam struct {
		UserID    int64    `json:"user_id"`
		ClientID  int64    `json:"client_id"`
		QueueName string   `json:"queue_name"`
		Addresses []string `json:"addresses"`
	}

	params := &registerParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	client := &models.Client{
		UserID:    params.UserID,
		ID:        params.ClientID,
		QueueName: params.QueueName,
		Publisher: a.Address,
		Addresses: params.Addresses,
	}

	err := a.QueueService.RegisterClient(client)
	responseErr(ctx, err)
}

// ApplyMessageIDRange try to apply the message id range for a queue
func (a *Agent) ApplyMessageIDRange(ctx *gin.Context) {
	var params = struct {
		UserID    int64  `json:"user_id"`
		QueueName string `json:"queue_name"`
		Size      int    `json:"size"`
	}{}

	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		log.Biz.Infoln(err)
		responseErr(ctx, err)
		return
	}

	maxID, err := a.QueueService.ApplyMessageIDRange(params.UserID, params.QueueName, params.Size)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	var res = struct {
		MessageIDBegin int64 `json:"message_id_begin"`
		MessageIDEnd   int64 `json:"message_id_end"`
	}{
		MessageIDBegin: maxID - int64(params.Size-1),
		MessageIDEnd:   maxID,
	}

	responseOKData(ctx, res)
}
