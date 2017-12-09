package urlutil

import (
	"fmt"
	"strings"

	"github.com/Focinfi/gosqs/config"
)

// MakeURL make a url for some abbreviation addr like ":12345"
// TODO: fit in more scene
func MakeURL(addr string) string {
	protocol := "http"
	if config.Config.Env.IsProduction() {
		protocol = "https"
	}

	if strings.HasPrefix(addr, ":") {
		return fmt.Sprintf("%s://127.0.0.1%s", protocol, addr)
	}

	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		return fmt.Sprintf("%s://%s", protocol, addr)
	}

	return addr
}
