package token

import (
	"time"
)

// Tokener defines the abilities of tokener
type Tokener interface {
	Make(secret string, params map[string]interface{}, expiration time.Duration) (string, error)
	Verify(code string, secret string) (map[string]interface{}, error)
}

// Default is a ready to use Tokener
var Default Tokener

func init() {
	Default = NewJWT()
}
