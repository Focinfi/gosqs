package agent

import (
	"net/http"
	"path"

	"github.com/Focinfi/gosqs/admin"
	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/errors"
	"github.com/Focinfi/gosqs/external"
	"github.com/Focinfi/gosqs/log"
	"github.com/Focinfi/gosqs/models"
	"github.com/Focinfi/gosqs/util/httputil"
	"github.com/Focinfi/gosqs/util/urlutil"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var root = path.Join(config.Root(), "master")

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
	s.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/home")
	})
	// fe
	s.LoadHTMLFiles(path.Join(root, "index.html"))
	group := s.Group("/", setAccessControlAllowHeaders)
	group.StaticFS("/static", http.Dir(path.Join(root, "static")))
	group.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File(path.Join(root, "favicon.ico"))
	})
	group.GET("/home", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// nodes
	group.OPTIONS("/applyNode")
	group.POST("/applyNode", a.handleApplyNode)
	group.POST("/join", a.handleJoinNode)

	// sendKeys
	group.GET("/sendGithubEmailSecretKey/:login", admin.SendGithubEmailSecretKey)
	a.Handler = s
}

// handleJoinNode for joining a new node
func (a *MasterAgent) handleJoinNode(ctx *gin.Context) {
	params := models.NodeInfo{}
	if err := binding.JSON.Bind(ctx.Request, &params); err != nil {
		log.Biz.Error(err)
		httputil.ResponseErr(ctx, err)
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
		httputil.ResponseErr(ctx, err)
		return
	}

	err := defaultAuth.Authenticate(params.AccessKey, params.SecretKey)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	userID, getUserIDErr := external.DefaultUserStore.GetUserIDByUniqueID(params.AccessKey)
	if getUserIDErr == errors.DataNotFound {
		id, err := external.DefaultUserStore.CreateUserByUniqueID(params.AccessKey)
		if err != nil {
			httputil.ResponseErr(ctx, err)
			return
		}

		userID = id
	}
	if getUserIDErr != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	node, err := a.MasterService.AssignNode(userID, params.QueueName, params.SquadName)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	tokenCode, err := makeToken(userID)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	httputil.ResponseOKData(ctx, gin.H{"node": urlutil.MakeURL(node), "token": tokenCode})
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

	messages, err := a.QueueService.PullMessages(userID, params.QueueName, params.SquadName, 10)
	if err != nil {
		httputil.ResponseErr(ctx, err)
		return
	}

	httputil.ResponseOKData(ctx, messages)
}
