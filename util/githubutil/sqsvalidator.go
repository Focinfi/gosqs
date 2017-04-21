package githubutil

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
	apiURL       = "https://api.github.com/repos/Focinfi/sqs/stargazers"
	updatePeriod = time.Second * 10
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

// SQSValidator auth by github stargazers
type SQSValidator struct {
	sync.RWMutex
	stargazerSlice []string
}

// NewGithubValidator allocates and return a new SQSValidator
func NewGithubValidator() *SQSValidator {
	return &SQSValidator{
		stargazerSlice: []string{},
	}
}

var DefaultValidator = NewGithubValidator()

// Start start the background services
func (auther *SQSValidator) Start() {
	if err := auther.updateStargazerSlice(); err != nil {
		panic(err)
	}
	go auther.updateRegularly()
}

func (auther *SQSValidator) ContainsLogin(login string) bool {
	auther.RLock()
	defer auther.RUnlock()

	idx := sort.SearchStrings(auther.stargazerSlice, login)
	return idx < len(auther.stargazerSlice) && auther.stargazerSlice[idx] == login
}

// Validate validates use accessKey as a the github login,
// secretKey encrypted the data of accessKey.
func (auther *SQSValidator) Validate(accessKey string, secretKey string) error {
	if !auther.ContainsLogin(accessKey) {
		return errors.UserNotFound
	}

	return nil
}

// update stargazerSlice regularly
func (auther *SQSValidator) updateRegularly() {
	ticker := time.NewTicker(updatePeriod)
	for {
		<-ticker.C
		auther.updateStargazerSlice()
	}
}

func (auther *SQSValidator) updateStargazerSlice() error {
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
