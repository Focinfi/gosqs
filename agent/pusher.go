package agent

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/Focinfi/sqs/log"
)

// PushMessage push message to all clients
func (QueueAgent) PushMessage(addresses []string, message string) chan bool {
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
			log.Biz.Printf("RESP: %v, ERR: %v\n", resp, err)
			if err == nil && resp.StatusCode == http.StatusOK {
				if done {
					return
				}

				lock.Lock()
				done = true
				lock.Unlock()
				pushed <- true
			}
		}()
	}

	return pushed
}
