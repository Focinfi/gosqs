package models

import (
	"fmt"
	"time"
)

const SquadKeyPrefix = "sqs.squad"

func SquadKey(userID int64, queueName string, squadName string) string {
	return fmt.Sprintf("%s.%d.%s.%s", SquadKeyPrefix, userID, queueName, squadName)
}

// Squad is a record for one queue processed index
type Squad struct {
	Name              string
	UserID            int64
	QueueName         string
	ReceivedMessageID int64
	RecentPushedAt    time.Time
}

func (s *Squad) Key() string {
	return SquadKey(s.UserID, s.QueueName, s.Name)
}
