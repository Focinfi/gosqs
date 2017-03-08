package models

import "fmt"

// Client for client model
type Client struct {
	ID                 int64  `json:"-"`
	UserID             int64  `json:"-"`
	QueueName          string `json:"-"`
	RecentMessageIndex int64
	Addresses          []string
	Publisher          string
	// RecentPushedAt Unix timestamp(s)
	RecentPushedAt int64
}

// ClientKeyPerfix for prefix storage key
const ClientKeyPerfix = "sqs.client"

// ClientKey for client key
func ClientKey(userID, clientID int64, queueName string) string {
	return fmt.Sprintf("%s.%d.%d.%s", ClientKeyPerfix, userID, clientID, queueName)
}

// Key return key for the c
func (c Client) Key() string {
	return ClientKey(c.UserID, c.ID, c.QueueName)
}
