package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status contains info of a failed request
type Status struct {
	Code    int    `json:"code,string"`
	Message string `json:"message"`
}

const (
	okCode         = 1000
	wrongParamCode = 1001
	isBusyCode     = 1002
)

// StatusWrongParam for wrong param response
var StatusWrongParam = &Status{Code: wrongParamCode}

// StatusIsBusy for wrong param response
func StatusIsBusy(err error) *Status {
	return &Status{
		Code:    isBusyCode,
		Message: err.Error(),
	}
}

// StatusOK for successful request
var StatusOK = &Status{Code: okCode}

func response(ctx *gin.Context, err *Status, isAbort bool) {
	ctx.JSON(http.StatusOK, err)
	if isAbort {
		ctx.Abort()
	}
}

func responseOK(ctx *gin.Context) {
	response(ctx, StatusOK, true)
}

func responseAndAbort(ctx *gin.Context, err *Status) {
	response(ctx, err, true)
}
