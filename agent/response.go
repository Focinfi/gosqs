package agent

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Focinfi/sqs/errors"
	"github.com/gin-gonic/gin"
)

// Status contains Info of a failed request
type Status struct {
	Code    int         `json:"code,string"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// StatusIsBusy for internal error
func StatusIsBusy(err error) *Status {
	log.Printf("error: %v\n", err)
	return &Status{
		Code:    errors.InternalErr,
		Message: "Service is busy, please try again later.",
	}
}

// StatusBadRequest for wrong param format status
func StatusBadRequest(err error) *Status {
	return &Status{
		Code:    errors.ParamFormatErr,
		Message: fmt.Sprintf("Wrong format of parameters, err: %v", err),
	}
}

// StatusOK for successful request
var StatusOK = &Status{Code: errors.NoErr}

func responseJOSN(ctx *gin.Context, status *Status, isAbort bool) {
	ctx.JSON(http.StatusOK, status)
	if isAbort {
		ctx.Abort()
	}
}

func responseOK(ctx *gin.Context) {
	responseJOSN(ctx, StatusOK, true)
}

func responseOKData(ctx *gin.Context, data interface{}) {
	responseJOSN(ctx, &Status{Code: errors.NoErr, Data: data}, true)
}

func responseAndAbort(ctx *gin.Context, err *Status) {
	responseJOSN(ctx, err, true)
}

func responseErr(ctx *gin.Context, err error) {
	if err == nil {
		responseOK(ctx)
		return
	}

	if bizErr, ok := err.(errors.Biz); ok {
		responseJOSN(ctx, &Status{Code: bizErr.BizCode(), Message: bizErr.Error()}, true)
		return
	}

	if internalErr, ok := err.(errors.Internal); ok {
		responseJOSN(ctx, StatusIsBusy(internalErr), true)
		return
	}

	responseJOSN(ctx, StatusBadRequest(err), true)
	return
}
