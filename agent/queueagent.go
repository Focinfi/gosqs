package agent

import (
	"net/http"

	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// QueueAgent for a queue HTTP agent
type QueueAgent struct {
	Address string
	http.Handler
	QueueService
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
	group.POST("/messages", a.handlePullMessages)
	group.POST("/message", a.handlePushMessage)
	group.POST("/messageID", a.handleApplyMessageIDRange)
	group.POST("/receivedMessageID", a.handleReportMaxReceivedMessageID)
	group.GET("/stats", a.handleGetStatus)
	a.Handler = s
}

// handleReportMaxReceivedMessageID handles the request for reporting the max id of received messages
func (a *QueueAgent) handleReportMaxReceivedMessageID(ctx *gin.Context) {
	params := &struct {
		models.NodeRequestParams
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

// handlePushMessage serve message pushing via http
func (a *QueueAgent) handlePushMessage(ctx *gin.Context) {
	type messageParam struct {
		models.NodeRequestParams
		Content   string `json:"content"`
		MessageID int64  `json:"message_id"`
	}
	params := &messageParam{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		responseErr(ctx, err)
		return
	}
	log.Internal.Infoln("[handlePushMessage]", params)

	userID, err := getUserID(params.Token)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	err = a.QueueService.PushMessage(userID, params.QueueName, params.Content, params.MessageID)
	responseErr(ctx, err)
}

// handleApplyMessageIDRange try to apply the message id range for a queue
func (a *QueueAgent) handleApplyMessageIDRange(ctx *gin.Context) {
	var params = struct {
		models.NodeRequestParams
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

// handleGetStatus response the info of current node
func (a *QueueAgent) handleGetStatus(ctx *gin.Context) {
	responseOKData(ctx, a.QueueService.Info())
}
