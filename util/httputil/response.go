package httputil

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Focinfi/gosqs/errors"
	"github.com/Focinfi/gosqs/models"
	"github.com/gin-gonic/gin"
)

// StatusIsBusy for internal error
func StatusIsBusy(err error) *models.HTTPStatusMeta {
	log.Printf("error: %v\n", err)
	return &models.HTTPStatusMeta{
		Code:    errors.InternalErr,
		Message: "Service is busy, please try again later.",
	}
}

// StatusBadRequest for wrong param format status
func StatusBadRequest(err error) *models.HTTPStatusMeta {
	return &models.HTTPStatusMeta{
		Code:    errors.ParamFormatErr,
		Message: fmt.Sprintf("Wrong format of parameters, err: %v", err),
	}
}

// StatusOK for successful request
var StatusOK = &models.HTTPStatusMeta{Code: errors.NoErr}

// ResponseJOSN put the data as JSON format into the reponse body
func ResponseJOSN(ctx *gin.Context, meta *models.HTTPStatusMeta, data interface{}, isAbort bool) {
	ctx.JSON(http.StatusOK, models.HTTPStatus{HTTPStatusMeta: *meta, Data: data})
	if isAbort {
		ctx.Abort()
	}
}

// ResponseOK just 200 without body data
func ResponseOK(ctx *gin.Context) {
	ResponseJOSN(ctx, StatusOK, nil, true)
}

// ResponseOKData responses 200 and JSON data
func ResponseOKData(ctx *gin.Context, data interface{}) {
	ResponseJOSN(ctx, StatusOK, data, true)
}

// ResponseErr responses depends on the err type
func ResponseErr(ctx *gin.Context, err error) {
	if err == nil {
		ResponseOK(ctx)
		return
	}

	if bizErr, ok := err.(errors.Biz); ok {
		ResponseJOSN(ctx, &models.HTTPStatusMeta{Code: bizErr.BizCode(), Message: bizErr.Error()}, nil, true)
		return
	}

	if internalErr, ok := err.(errors.Internal); ok {
		ResponseJOSN(ctx, StatusIsBusy(internalErr), nil, true)
		return
	}

	ResponseJOSN(ctx, StatusBadRequest(err), nil, true)
	return
}
