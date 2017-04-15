package agent

import (
	"net/http"

	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// MasterAgent for a agent of master node
type MasterAgent struct {
	Address string
	http.Handler
	MasterService
}

// NewMasterAgent allocates a new MasterAgent to handle the HTTP API
func NewMasterAgent(service MasterService, address string) *MasterAgent {
	agt := &MasterAgent{
		Address:       address,
		MasterService: service,
	}

	agt.masterAgentRouting()
	return agt
}

func (a *MasterAgent) masterAgentRouting() {
	s := gin.Default()
	group := s.Group("/")
	group.POST("/applyNode", a.handleApplyNode)
	group.POST("/join", a.handleJoinNode)
	a.Handler = s
}

// handleJoinNode for joining a new node
func (a *MasterAgent) handleJoinNode(ctx *gin.Context) {
	params := models.NodeInfo{}
	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		log.Biz.Error(err)
		responseErr(ctx, err)
		return
	}

	a.MasterService.AddNode(params)
}

// handleApplyNode register a consumer which is ready to pull messages for a squad of a queue
func (a *MasterAgent) handleApplyNode(ctx *gin.Context) {
	params := &struct {
		models.UserAuth
		models.NodeRequestParams
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
	tokenCode, err := makeToken(userID)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"node": node, "token": tokenCode})
}

// handlePullMessages for pulling message
func (a *QueueAgent) handlePullMessages(ctx *gin.Context) {
	params := &models.NodeRequestParams{}
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

	log.Internal.Infoln("[handlePullMessages] userID:", userID)
	messages, err := a.QueueService.PullMessages(userID, params.QueueName, params.SquadName, 10)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"messages": messages})
}
