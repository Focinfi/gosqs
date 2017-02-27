package admin

import (
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/queue"
)

// Service for one user info
type Service struct {
	user   external.User
	queues map[string]queue.Queue
}

// root for initialization
var root = external.UserFunc(func() int64 { return 1 })

// defaultService for initialization
var defaultService = Service{user: root, queues: map[string]queue.Queue{}}

// Start starts services
func Start(address string) {

}
