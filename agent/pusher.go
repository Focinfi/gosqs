package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// PushMessage push message to all clients
func (Agent) PushMessage(addresses []string, message string) chan bool {
	if len(addresses) == 0 || message == "" {
		return nil
	}

	pushed := make(chan bool)
	var lock sync.Mutex
	var done bool
	for _, address := range addresses {
		addr := address
		go func() {
			resp, err := http.PostForm(addr, url.Values{"message": {message}})
			fmt.Printf("RESP: %v, ERR: %v\n", resp, err)
			if err == nil && resp.StatusCode == http.StatusOK {
				if done {
					return
				}

				lock.Lock()
				done = true
				lock.Unlock()

				fmt.Printf("DONE: %v\n", done)

				pushed <- true
			}

			fmt.Printf("PUSH MSG ERROR: %v\n", err)
		}()
	}

	return pushed
}
