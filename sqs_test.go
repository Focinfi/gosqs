package sqs

import (
	"testing"
	"time"

	"github.com/Focinfi/sqs/client"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/master"
	"github.com/Focinfi/sqs/node"
	"github.com/Focinfi/sqs/util/token"
)

func Test(t *testing.T) {
	go master.NewService(":54661").Start()
	// wait a moment
	time.Sleep(time.Second)
	go node.New(":54662", 54462, ":54661").Start()

	time.Sleep(time.Second)
	accessKey := "Focinfi"
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: accessKey}, time.Hour)
	cli := client.New(config.Config.DefaultMasterAddress, accessKey, secretKey)
	queueCli, err := cli.Queue(example.Greeting, example.Home)
	if err != nil {
		t.Fatal("failed to create a queue, err:", err)
	}

	if err := queueCli.ApplyNode(); err != nil {
		t.Fatal("failed to apply node, err:", err)
	}

	if err := queueCli.PushMessage("foo"); err != nil {
		t.Fatal("failed to push message, err:", err)
	}

	time.Sleep(time.Second)
	queueCli.PullMessages(func(messages []client.Message) error {
		if len(messages) != 0 || messages[0].Content != "foo" {
			t.Fatalf("can not pull messages, got: %v", messages)
		}
		return nil
	})
}
