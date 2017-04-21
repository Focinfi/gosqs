package agent

import (
	"reflect"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/util/githubutil"
	"github.com/Focinfi/sqs/util/token"
)

// Validator authenticates the given accessKey and secretKey.
type Validator interface {
	Validate(accessKey string, secretKey string) (err error)
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
		Validators: []Validator{githubutil.DefaultValidator},
	}
}
