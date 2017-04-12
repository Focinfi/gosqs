package agent

import (
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/node"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type basicParam struct {
	UserID    int64  `json:"user_id"`
	QueueName string `json:"queue_name"`
	SquadName string `json:"squad_name,omitempty"`
}

// JoinNode for joining a new node
func (a *MasterAgent) JoinNode(ctx *gin.Context) {
	params := node.Info{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	a.MasterService.Join(params)
}

// ApplyNode register a consumer which is ready to pull messages for a squad of a queue
func (a *MasterAgent) ApplyNode(ctx *gin.Context) {
	params := &basicParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	node, err := a.MasterService.AssignNode(params.UserID, params.QueueName, params.SquadName)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"node": node})
}

func (a *QueueAgent) PullMessages(ctx *gin.Context) {
}

func (a *QueueAgent) ReportMaxReceivedMessageID(ctx *gin.Context) {
	params := &struct {
		basicParam
		MessageID int64 `json:"message_id"`
	}{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	err := a.QueueService.ReportMaxReceivedMessageID(params.UserID, params.QueueName, params.SquadName, params.MessageID)
	responseErr(ctx, err)
}

// PushMessage serve message pushing via http
func (a *QueueAgent) ReceiveMessage(ctx *gin.Context) {
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

	err := a.QueueService.PushMessage(params.UserID, params.QueueName, params.Content, params.Index)
	responseErr(ctx, err)
}

// ApplyMessageIDRange try to apply the message id range for a queue
func (a *QueueAgent) ApplyMessageIDRange(ctx *gin.Context) {
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

func (a *QueueAgent) Info(ctx *gin.Context) {
	responseOKData(ctx, a.QueueService.Info())
}
