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
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{userGithubLoginKey: accessKey}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	err = validator.Validate(accessKey, secretKey)
	if err != nil {
		t.Error(err)
	}
}
