package admin

import "time"

func messageIndex() int64 {
	return time.Now().UnixNano()
}
