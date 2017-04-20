package agent

import (
	"testing"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/util/token"
)

func Test(t *testing.T) {
	validator := NewGithubValidator()
	validator.Start()

	accessKey := "Focinfi"
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: accessKey}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	err = validator.Validate(accessKey, secretKey)
	if err != nil {
		t.Error(err)
	}
}
