package admin

import (
	"testing"
	"time"

	"github.com/Focinfi/sqs/util/token"
	"github.com/Focinfi/sqs/config"
)

func TestSendKeysEmail(t *testing.T) {
	accessKey := "Focinfi"
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: accessKey}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if err := sendSecretKeyToEmail("focinfi@qq.com", accessKey, secretKey); err != nil {
		t.Error(err)
	}
}
