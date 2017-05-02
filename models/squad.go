package models

import (
	"fmt"
	"time"
)

// SquadKeyPrefix prefix of the storage key for squad
const SquadKeyPrefix = "sqs.squad"

// SquadKey returns the key for the given params
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

// Key returns the key of this Squad
func (s *Squad) Key() string {
	return SquadKey(s.UserID, s.QueueName, s.Name)
}
