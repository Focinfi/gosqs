package agent

import (
	"net/http"

	"log"

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
	failedRegister = 1003
)

// StatusWrongParam for wrong param response
var StatusWrongParam = &Status{Code: wrongParamCode}

// StatusIsBusy for internal error
func StatusIsBusy(err error) *Status {
	log.Printf("error: %v\n", err)
	return &Status{
		Code:    isBusyCode,
		Message: "Service is busy, please try again later.",
	}
}

// StatusBizError for biz logic error
func StatusBizError(code int, err error) *Status {
	return &Status{
		Code:    code,
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
