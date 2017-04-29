package agent

import (
	"reflect"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/errors"
	"github.com/Focinfi/gosqs/util/githubutil"
	"github.com/Focinfi/gosqs/util/token"
)

// Validator authenticates the given accessKey and secretKey.
type Validator interface {
	Validate(accessKey string, secretKey string) (err error)
}

// ValidatorFunc func to interface helper
type ValidatorFunc func(accessKey string, secretKey string) (err error)

// Validate implements the Validator interface
func (v ValidatorFunc) Validate(accessKey string, secretKey string) (err error) {
	return v(accessKey, secretKey)
}

// Auth authes
type Auth struct {
	Validators []Validator
}

// Authenticate authenticates the accessKey and secretKey.
// Try to find the corresponding userID of the accessKey.
func (a Auth) Authenticate(accessKey string, secretKey string) (err error) {
	params, err := token.Default.Verify(secretKey, config.Config.BaseSecret)
	if err != nil {
		return errors.UserAuthError(err.Error())
	}
	if !reflect.DeepEqual(accessKey, params[config.Config.UserGithubLoginKey]) {
		return errors.UserAuthError("broken secrect_key")
	}

	for _, validator := range a.Validators {
		if err := validator.Validate(accessKey, secretKey); err == nil {
			return nil
		}
	}

	return errors.UserAuthError("failed to auth")
}

var defaultAuth *Auth

func init() {
	githubutil.DefaultValidator.Start()

	defaultAuth = &Auth{
		Validators: []Validator{testAuth, githubutil.DefaultValidator},
	}
}
