package agent

import (
	"net/http"
)

// Agent for receiving message and push them
type Agent struct {
	http.Handler
}

// ReceiveMessage serve message pushing via http
func (agent *Agent) ReceiveMessage() {

}

// DeliveryMessage deliveries messages to all online subsribers
func (agent *Agent) DeliveryMessage() {

}
