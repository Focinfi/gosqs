package models

import (
	"math/rand"
)

const base = int64(268435456) // 2 ** 28

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
	return index / base
}

// GenIndexRandom gen index for the time
func GenIndexRandom(timestamp int64) int64 {
	return GenIndex0(timestamp) + rand.Int63n(1000000)
}

// GenIndex0 gen index for the time, random index in n
func GenIndex0(timestamp int64) int64 {
	return timestamp * base
}

// MessageIndex for one entry of message index
type MessageIndex struct {
	Timestamp int64
	Indexes   []int64
}
