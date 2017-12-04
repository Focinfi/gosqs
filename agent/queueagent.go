package agent

import (
	"net/http"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/log"
	"github.com/Focinfi/gosqs/models"
	"github.com/Focinfi/gosqs/util/httputil"
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
	group := s.Group("/", setAccessControlAllowHeaders)
	group.OPTIONS("/messages")
	group.OPTIONS("/message")
	group.OPTIONS("/messageID")
	group.OPTIONS("/receivedMessageID")
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
		httputil.ResponseErr(ctx, err)
		return
	}
	userID, err := getUserID(params.Token)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	err = a.QueueService.ReportMaxReceivedMessageID(userID, params.QueueName, params.SquadName, params.MessageID)
	httputil.ResponseErr(ctx, err)
}

// handlePullMessages for pulling message
func (a *QueueAgent) handlePullMessages(ctx *gin.Context) {
	params := &models.NodeRequestParams{}
	if err := binding.JSON.Bind(ctx.Request, params); err != nil {
		log.Biz.Error(err)
		httputil.ResponseErr(ctx, err)
		return
	}

	userID, err := getUserID(params.Token)
	if err != nil {
		log.Biz.Error(err)
		httputil.ResponseErr(ctx, err)
		return
	}

	messages, err := a.QueueService.PullMessages(userID, params.QueueName, params.SquadName, config.Config.PullMessageCount)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	httputil.ResponseOKData(ctx, messages)
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
		httputil.ResponseErr(ctx, err)
		return
	}

	userID, err := getUserID(params.Token)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	err = a.QueueService.PushMessage(userID, params.QueueName, params.Content, params.MessageID)
	httputil.ResponseErr(ctx, err)
}

// handleApplyMessageIDRange try to apply the message id range for a queue
func (a *QueueAgent) handleApplyMessageIDRange(ctx *gin.Context) {
	var params = struct {
		models.NodeRequestParams
		Size int `json:"size"`
	}{}
	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	userID, err := getUserID(params.Token)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	maxID, err := a.QueueService.ApplyMessageIDRange(userID, params.QueueName, params.Size)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	var res = struct {
		MessageIDBegin int64 `json:"message_id_begin"`
		MessageIDEnd   int64 `json:"message_id_end"`
	}{
		MessageIDBegin: maxID - int64(params.Size-1),
		MessageIDEnd:   maxID,
	}

	httputil.ResponseOKData(ctx, res)
}

// handleGetStatus response the info of current node
func (a *QueueAgent) handleGetStatus(ctx *gin.Context) {
	httputil.ResponseOKData(ctx, a.QueueService.Info())
}
