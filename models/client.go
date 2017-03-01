package models

// Client for client model
type Client struct {
	ID                 int64
	UserID             int64
	QueueName          string
	RecentMessageIndex int64
}
