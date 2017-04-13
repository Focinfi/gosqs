package token

import (
	"time"
)

// Tokener defines the tokener behavior
type Tokener interface {
}

func Make(secret string, params map[string]string, expiration time.Duration) (string, error) {
	// TODO: to impelement
	return "mock.token", nil
}

// Verify check the code and returns the data pairs
func Verify(code string) (map[string]string, error) {
	// TODO: to implement
	return map[string]string{"userID": "1"}, nil
}
