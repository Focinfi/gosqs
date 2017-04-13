package agent

import (
	"net/http"

	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MasterAgent struct {
	Address string
	http.Handler
	MasterService
}

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
	group.POST("/applyNode", a.ApplyNode)
	group.POST("/join", a.JoinNode)
	a.Handler = s
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
	tokenCode, err := makeToekn(userID)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"node": node, "token": tokenCode})
}

// PullMessages for pulling message
func (a *QueueAgent) PullMessages(ctx *gin.Context) {
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

	log.Internal.Infoln("[PullMessage] userID:", userID)
	messages, err := a.QueueService.PullMessage(userID, params.QueueName, params.SquadName, 10)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseOKData(ctx, gin.H{"messages": messages})
}
