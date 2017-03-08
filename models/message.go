package models

import (
	"fmt"
	"math/rand"
)

const BaseUnit = int64(268435456) // 2 ** 28

// Message contains info
type Message struct {
	UserID    int64
	QueueName string
	Content   string
	// Index: |1bit|35bit for timestamp(m)|28bit for sequence|
	// example:
	//    if   timestamp = 1488350906, sequence = 1001
	//		then index = timestamp * 2**28 + sequence = 381017832937
	Index int64
}

// GroupID returns message group id based on index
func (m Message) GroupID() int64 {
	return GroupID(m.Index)
}

// GroupID returns message group id based on index
func GroupID(index int64) int64 {
	return index / BaseUnit
}

// GenIndexRandom gen index for the time
func GenIndexRandom(timestamp int64) int64 {
	return GenIndex0(timestamp) + rand.Int63n(1000000)
}

// GenIndex0 gen index for the time, random index in n
func GenIndex0(timestamp int64) int64 {
	return timestamp * BaseUnit
}

// MessageIndex for one entry of message index
type MessageIndex struct {
	Timestamp int64
	Indexes   []int64
}

const messageKeyPrefix = "sqs.message"

// MessageListKey for message list storage key
func MessageListKey(userID int64, queueName string, gorupID int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, gorupID)
}

// MessageKey for message storage key
func MessageKey(userID int64, queueName string, index int64) string {
	return fmt.Sprintf("%s.%d.%s.%d", messageKeyPrefix, userID, queueName, index)
}
