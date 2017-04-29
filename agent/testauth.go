package agent

import "fmt"

var testAuth = ValidatorFunc(func(accessKey string, secretKey string) (err error) {
	if accessKey == "test" {
		return nil
	}

	return fmt.Errorf("%s is not a test user\n", accessKey)
})
