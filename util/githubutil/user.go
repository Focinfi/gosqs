package githubutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiUserURLFormat = "https://api.github.com/users/%s?access_token=%s"

// EmailForUserLogin gets the email of the given login
func EmailForUserLogin(login string) (string, error) {
	url := fmt.Sprintf(apiUserURLFormat, login, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	params := struct {
		Email string `json:"email"`
	}{}
	if err := json.Unmarshal(respBytes, &params); err != nil {
		return "", err
	}

	return params.Email, nil
}
