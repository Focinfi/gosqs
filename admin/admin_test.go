package admin

import (
	"testing"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/util/token"
)

func TestSendKeysEmail(t *testing.T) {
	accessKey := "Focinfi"
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: accessKey}, -1)
	t.Log(secretKey, err)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if err := sendSecretKeyToEmail("focinfi@qq.com", accessKey, secretKey); err != nil {
	// 	t.Error(err)
	// }
}
