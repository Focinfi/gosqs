package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error contains info of a failed request
type Error struct {
	Code    int    `json:"code,string"`
	Message string `json:"message"`
}

const (
	wrongParamCode = 1001
	isBusyCode     = 1002
)

// ErrWrongParam for wrong param response
var ErrWrongParam = &Error{Code: wrongParamCode}

// ErrIsBusy for wrong param response
func ErrIsBusy(err error) *Error {
	return &Error{
		Code:    isBusyCode,
		Message: err.Error(),
	}
}

func response(ctx *gin.Context, err *Error, isAbort bool) {
	ctx.JSON(http.StatusOK, err)
	if isAbort {
		ctx.Abort()
	}
}

func responseOK(ctx *gin.Context) {
	response(ctx, nil, true)
}

func responseAndAbort(ctx *gin.Context, err *Error) {
	response(ctx, err, true)
}
