package agent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Focinfi/sqs/errors"
)

const (
	userGithubLoginKey = "userGithubLogin"
	apiURL             = "https://api.github.com/repos/Focinfi/sqs/stargazers"
	updatePeriod       = time.Second * 10
)

type stargazer struct {
	Login string `json:"login"`
}

type stargazers []stargazer

func (s stargazers) toSlice() []string {
	slice := make([]string, len(s))
	for i, stargazer := range s {
		slice[i] = stargazer.Login
	}

	return slice
}

// GithubValidator auth by github stargazers
type GithubValidator struct {
	sync.Mutex
	stargazerSlice []string
}

// NewGithubValidator allocates and return a new GithubValidator
func NewGithubValidator() *GithubValidator {
	return &GithubValidator{
		stargazerSlice: []string{},
	}
}

var githubValidator = NewGithubValidator()

// Start start the background services
func (auther *GithubValidator) Start() {
	if err := auther.updateStargazerSlice(); err != nil {
		panic(err)
	}
	go auther.updateRegularly()
}

// Validate validates use accessKey as a the github login,
// secretKey encrypted the data of accessKey.
func (auther *GithubValidator) Validate(accessKey string, secretKey string) error {

	idx := sort.SearchStrings(auther.stargazerSlice, accessKey)
	if idx > len(auther.stargazerSlice) || auther.stargazerSlice[idx] != accessKey {
		return errors.UserNotFound
	}

	return nil
}

// update stargazerSlice regularly
func (auther *GithubValidator) updateRegularly() {
	ticker := time.NewTicker(updatePeriod)
	for {
		<-ticker.C
		auther.updateStargazerSlice()
	}
}

func (auther *GithubValidator) updateStargazerSlice() error {
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respParam := []stargazer{}
	if err := json.Unmarshal(respBytes, &respParam); err != nil {
		return err
	}

	stargazersSlice := stargazers(respParam).toSlice()
	sort.StringSlice(stargazersSlice).Sort()
	auther.Lock()
	auther.stargazerSlice = stargazersSlice
	auther.Unlock()
	return nil
}
