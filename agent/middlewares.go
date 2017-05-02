package agent

import (
	"errors"

	"time"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/util/strconvutil"
	"github.com/Focinfi/gosqs/util/token"
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

func setAccessControlAllowHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Accept")
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
