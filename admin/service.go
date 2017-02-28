package admin

import (
	"log"
	"net/http"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/external"
)

// Service for one user info
type Service struct {
	*agent.Agent
}

// root for initialization
var root = external.UserFunc(func() int64 { return 1 })

// defaultService for initialization
var defaultService = Service{
	Agent: agent.New(defualtQueues),
}

// Start starts services
func Start(address string) {
	log.Fatal(http.ListenAndServe(address, defaultService.Agent))
}
