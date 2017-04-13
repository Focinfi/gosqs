package agent

import (
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type basicParam struct {
	Token     string `json:"token"`
	QueueName string `json:"queue_name"`
	SquadName string `json:"squad_name,omitempty"`
}

// JoinNode for joining a new node
func (a *MasterAgent) JoinNode(ctx *gin.Context) {
	params := models.NodeInfo{}
	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		log.Biz.Error(err)
		responseErr(ctx, err)
		return
	}

	a.MasterService.Join(params)
}

// ApplyNode register a consumer which is ready to pull messages for a squad of a queue
func (a *MasterAgent) ApplyNode(ctx *gin.Context) {
	params := &struct {
		models.UserAuth
		basicParam
	}{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}

	userID, err := external.GetUserWithKey(params.AccessKey, params.SecretKey)
	if err != nil {
		responseErr(ctx, err)
	}

	node, err := a.MasterService.AssignNode(userID, params.QueueName, params.SquadName)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	log.Internal.Println("AssignNode:", node)
	tokenCode, err := makeToekn(userID)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"node": node, "token": tokenCode})
}

// PullMessages for pulling message
func (a *QueueAgent) PullMessages(ctx *gin.Context) {
	params := &basicParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		log.Internal.Error(err)
		responseErr(ctx, err)
		return
	}

	userID, err := getUserID(params.Token)
	if err != nil {
		log.Internal.Error(err)
		responseErr(ctx, err)
		return
	}

	log.Internal.Infoln("[PullMessage] userID:", userID)
	messages, err := a.QueueService.PullMessage(userID, params.QueueName, params.SquadName, 10)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"messages": messages})
}

// ReportMaxReceivedMessageID handles the request for reporting the max id of received messages
func (a *QueueAgent) ReportMaxReceivedMessageID(ctx *gin.Context) {
	params := &struct {
		basicParam
		MessageID int64 `json:"message_id"`
	}{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}
	userID, err := getUserID(params.Token)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	err = a.QueueService.ReportMaxReceivedMessageID(userID, params.QueueName, params.SquadName, params.MessageID)
	responseErr(ctx, err)
}

// ReceiveMessage serve message pushing via http
func (a *QueueAgent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		basicParam
		Content   string `json:"content"`
		MessageID int64  `json:"message_id"`
	}
	params := &messageParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}
	log.Internal.Infoln("[ReceiveMessage]", params)

	userID, err := getUserID(params.Token)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	err = a.QueueService.PushMessage(userID, params.QueueName, params.Content, params.MessageID)
	responseErr(ctx, err)
}

// ApplyMessageIDRange try to apply the message id range for a queue
func (a *QueueAgent) ApplyMessageIDRange(ctx *gin.Context) {
	var params = struct {
		basicParam
		Size int `json:"size"`
	}{}
	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		log.Biz.Infoln(err)
		responseErr(ctx, err)
		return
	}

	userID, err := getUserID(params.Token)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	maxID, err := a.QueueService.ApplyMessageIDRange(userID, params.QueueName, params.Size)
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

// Info response the info of current node
func (a *QueueAgent) Info(ctx *gin.Context) {
	responseOKData(ctx, a.QueueService.Info())
}
