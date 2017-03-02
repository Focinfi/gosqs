package models

// Client for client model
type Client struct {
	ID                 int64  `json:"-"`
	UserID             int64  `json:"-"`
	QueueName          string `json:"-"`
	RecentMessageIndex int64
	Addresses          []string
	Publisher          string
	// RecentReceivedAt Unix timestamp(s)
	RecentReceivedAt int64
}
