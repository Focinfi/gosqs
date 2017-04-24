package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT implements the Tokener
type JWT struct {
	*jwt.Token
}

// NewJWT allocates and returns a new Token with name and expiration mins
func NewJWT() *JWT {
	return &JWT{Token: jwt.New(jwt.SigningMethodHS256)}
}

// Make makes a jwt token string.
// if expiration <= 0, no expiration.
func (j *JWT) Make(secret string, params map[string]interface{}, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	if expiration > 0 {
		claims["exp"] = time.Now().Add(expiration).Unix()
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	for k, v := range params {
		t.Claims.(jwt.MapClaims)[k] = v
	}

	return t.SignedString([]byte(secret))
}

// Verify verifies the code
func (j *JWT) Verify(code string, secret string) (map[string]interface{}, error) {
	t, err := jwt.Parse(code, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return t.Claims.(jwt.MapClaims), nil
}
