package agent

import (
	"net/http"
	"net/url"
)

// StartPushing start pushing message
func (a *Agent) StartPushing() error {
	return nil
}

// PushMessage push message to all clients
func (a *Agent) PushMessage(addresses []string, message string) chan bool {
	if len(addresses) == 0 || message == "" {
		return nil
	}

	pushed := make(chan bool)
	for _, address := range addresses {
		addr := address
		go func() {
			resp, err := http.PostForm(addr, url.Values{"message": {message}})
			if err == nil && resp.StatusCode == http.StatusOK {
				pushed <- true
			}

		}()
	}

	return pushed
}
