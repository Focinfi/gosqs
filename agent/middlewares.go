package agent

import (
	"errors"

	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/util/strconvutil"
	"github.com/Focinfi/sqs/util/token"
	"github.com/gin-gonic/gin"
)

const (
	userIDKey = "userID"
)

var (
	tokenExpiration = time.Duration(-1) // no expiration
	tokener         = token.Default
	baseSecret      = config.Config.BaseSecret
)

// throttling protects our server from overload
func throttling(ctx *gin.Context) {
}

// parsing parses params in the req for the following middleware
func parsing(ctx *gin.Context) {

}

// auth authenticates identity for the req
func auth(ctx *gin.Context) {
}

func getUserID(code string) (int64, error) {
	data, err := tokener.Verify(code, baseSecret)
	if err != nil {
		return -1, err
	}

	idStr, ok := data[userIDKey]
	if !ok {
		return -1, errors.New("broken code")
	}

	return strconvutil.ParseInt64(idStr.(string))
}

func makeToken(userID int64) (string, error) {
	return tokener.Make(baseSecret, map[string]interface{}{userIDKey: strconvutil.Int64toa(userID)}, tokenExpiration)
}
