package urlutil

import (
	"fmt"
	"strings"
)

// MakeURL make a url for some abbreviation addr like ":12345"
// TODO: fit in more scene
func MakeURL(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return fmt.Sprintf("http://127.0.0.1%s", addr)
	}

	if !strings.HasPrefix(addr, "http://") || !strings.HasPrefix(addr, "https://") {
		return fmt.Sprintf("http://%s", addr)
	}

	return addr
}
