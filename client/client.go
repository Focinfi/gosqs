package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/util/urlutil"
)

const (
	jsonHTTPHeader = "application/json"

	registerURLFormat                = "%s/register"
	applyMessageIDURLFormat          = "%s/messageID"
	pushMessageURLFormat             = "%s/message"
	pullMessageURLFormat             = "%s/messages"
	reportReceivedMessageIDURLFormat = "%s/receivedMessageID"
)

type Option struct {
	// Endpoint for main server
	Endpoint string
	//AccessKey for sqs basic key
	AccessKey string
	// Secret for user auth
	SecretKey string
}

// Client for one sqs client
type Client struct {
	opt *Option
}

// QueueClient for one query client
type QueueClient struct {
	endpoint string
	serving  string
	BaseInfo
}

// BaseInfo for one client basic info
type BaseInfo struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	QueueName string `json:"queue_name"`
}

type pushMessageParam struct {
	BaseInfo
	MessageID int64  `json:"message_id"`
	Content   string `json:"content"`
}

type applyMessageIDParam struct {
	BaseInfo
	Size int `json:"size"`
}

type applyMessageResponseParam struct {
	MessageIDBegin int64 `json:"message_id_begin"`
	MessageIDEnd   int64 `json:"message_id_end"`
}

type reportReceivedParam struct {
	BaseInfo
	MessageID int64 `json:"message_id"`
}

type Message struct {
	MessageID int64
	Content   string
}

// Queue returns a new QueueClient with the given name
func (cli *Client) Queue(name string) *QueueClient {
	return &QueueClient{
		endpoint: cli.opt.Endpoint,
		BaseInfo: BaseInfo{
			AccessKey: cli.opt.AccessKey,
			SecretKey: cli.opt.SecretKey,
			QueueName: name,
		},
	}
}

func (cli *QueueClient) Register() error {
	b, err := json.Marshal(cli.BaseInfo)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(registerURLFormat, urlutil.MakeURL(cli.endpoint))
	resp, err := http.Post(url, jsonHTTPHeader, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// PushMessage pushes a message
func (cli *QueueClient) PushMessage(content string) error {
	// apply a id
	id, err := cli.applyMessageID()
	if err != nil {
		return err
	}

	param := &pushMessageParam{
		MessageID: id,
		BaseInfo:  cli.BaseInfo,
	}

	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(pushMessageURLFormat, urlutil.MakeURL(cli.serving))
	resp, err := http.Post(url, jsonHTTPHeader, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push message, resonse status code: %d\n", resp.StatusCode)
	}

	return nil
}

// PullMessage for pull message request
func (cli *QueueClient) PullMessage() ([]Message, error) {
	url := fmt.Sprintf(pullMessageURLFormat, urlutil.MakeURL(cli.serving))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	messages := []Message{}
	if err := json.Unmarshal(b, messages); err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		go cli.reportReceived(messages[len(messages)-1].MessageID)
	}

	return messages, nil
}

// reportReceived reports the last received message id
func (cli *QueueClient) reportReceived(messageID int64) error {
	url := fmt.Sprintf(reportReceivedMessageIDURLFormat, cli.serving)
	param := &reportReceivedParam{
		BaseInfo:  cli.BaseInfo,
		MessageID: messageID,
	}

	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	var delay time.Duration

	for {
		select {
		case <-time.After(delay):
			resp, err := http.Post(url, jsonHTTPHeader, bytes.NewReader(b))
			if err != nil {
				log.Service.Errorf("can not report received message id, err: %v\n", err)
				delay = (delay + 1) * time.Millisecond * 500
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Service.Errorf("can not report received message id, status code: %d\n", resp.StatusCode)
				continue
			}

			return nil
		}
	}
}

func (cli *QueueClient) applyMessageID() (int64, error) {
	param := &applyMessageIDParam{cli.BaseInfo, 1}
	b, err := json.Marshal(param)
	if err != nil {
		return -1, err
	}

	url := fmt.Sprintf(applyMessageIDURLFormat, urlutil.MakeURL(cli.serving))
	resp, err := http.Post(url, jsonHTTPHeader, bytes.NewReader(b))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	messageID := &applyMessageResponseParam{}
	if err := json.Unmarshal(respBytes, messageID); err != nil {
		return -1, err
	}

	if messageID.MessageIDEnd < messageID.MessageIDBegin {
		return -1, errors.New("GET /appyMessageID response data broken: end < begin")
	}

	return messageID.MessageIDEnd, nil
}
