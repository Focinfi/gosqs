package admin

import (
	"github.com/Focinfi/sqs/queue"
)

var queues = []queue.Queue{}

// AddQueue adds a queue into queues
func AddQueue(name string, configs ...queue.Config) {
	q := queue.New(name, configs...)
	queues = append(queues, q)
}
