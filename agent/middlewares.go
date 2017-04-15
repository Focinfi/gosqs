package agent

import (
	"errors"

	"time"

	"github.com/Focinfi/sqs/util/strconvutil"
	"github.com/Focinfi/sqs/util/token"
	"github.com/gin-gonic/gin"
)

const (
	userIDKey = "userID"
)

var (
	// TODO: move into config
	secret          = "sqs.secret"
	tokenExpiration = time.Hour
)

// throttling protects our server from overload
func throttling(ctx *gin.Context) {
}

// parsing parses params in the req for the following middlewares
func parsing(ctx *gin.Context) {

}

// auth authenticates identity for the req
func auth(ctx *gin.Context) {
}

func getUserID(code string) (int64, error) {
	data, err := token.Verify(code)
	if err != nil {
		return -1, err
	}

	idStr, ok := data[userIDKey]
	if !ok {
		return -1, errors.New("broken code")
	}

	return strconvutil.ParseInt64(idStr)
}

func makeToken(userID int64) (string, error) {
	return token.Make(secret, map[string]string{userIDKey: strconvutil.Int64toa(userID)}, tokenExpiration)
}
