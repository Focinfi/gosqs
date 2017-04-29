package gosqs

import (
	"testing"
	"time"

	"github.com/Focinfi/gosqs/client"
	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/example"
	"github.com/Focinfi/gosqs/master"
	"github.com/Focinfi/gosqs/node"
	"github.com/Focinfi/gosqs/util/token"
)

func Test(t *testing.T) {
	masterAddr := ":54661"
	go master.NewService(masterAddr).Start()
	// wait a moment
	time.Sleep(time.Second)
	go node.New(":54462", 54462, masterAddr).Start()

	time.Sleep(time.Second)
	accessKey := "test"
	paramsKey := config.Config.UserGithubLoginKey
	secretKey, err := token.Default.Make(config.Config.BaseSecret, map[string]interface{}{paramsKey: accessKey}, -1)
	cli := client.New(masterAddr, accessKey, secretKey)
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
