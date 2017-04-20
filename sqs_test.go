package sqs

import (
	"testing"
	"time"

	"github.com/Focinfi/sqs/client"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/example"
	"github.com/Focinfi/sqs/master"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/node"
	"github.com/Focinfi/sqs/util/token"
)

func Test(t *testing.T) {
	queue := models.NewQueue(1, example.Greeting)
	if err := master.AddQueue(queue); err != nil {
		panic(err)
	}

	go master.NewService(config.Config.DefaultMasterAddress).Start()
	// wait a moment
	time.Sleep(time.Second)
	go node.New(":54461", ":5446").Start()

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
