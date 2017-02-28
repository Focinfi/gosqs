package agent

import "github.com/gin-gonic/gin"

// ReceiveMessage serve message pushing via http
func (agent *Agent) ReceiveMessage(ctx *gin.Context) {
	type messageParam struct {
		UserID    int64  `json:"userID"`
		QueueName string `json:"queueName"`
		Content   string `json:"content"`
	}

	params := &messageParam{}
	if err := ctx.BindJSON(params); err != nil {
		responseAndAbort(ctx, StatusWrongParam)
		return
	}

	err := agent.PushMessage(params.UserID, params.QueueName, params.Content)
	if err != nil {
		responseAndAbort(ctx, StatusIsBusy(err))
		return
	}

	responseOK(ctx)
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
		responseAndAbort(ctx, StatusWrongParam)
		return
	}

	err := agent.RegisterClient(param.UserID, param.ClientID, param.QueueName)
	if err != nil {
		responseAndAbort(ctx, StatusBizError(failedRegister, err))
		return
	}

	responseOK(ctx)
}

// StartDeliveryMessage deliveries messages to all online subsribers
func (agent *Agent) StartDeliveryMessage() {

}
