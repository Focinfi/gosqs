package githubutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/errors"
)

const (
	apiURL       = "https://api.github.com/repos/Focinfi/gosqs/stargazers"
	updatePeriod = time.Second * 10
)

var (
	accessToken     = config.Config.GithubAccessToken
	stargazerAPIURL = fmt.Sprintf("%s?access_token=%s", apiURL, accessToken)
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

// DefaultValidator for a default export SQSValidator
var DefaultValidator = NewGithubValidator()

// Start start the background services
func (v *SQSValidator) Start() {
	if err := v.updateStargazerSlice(); err != nil {
		panic(err)
	}
	go v.updateRegularly()
}

// ContainsLogin for checking the login
func (v *SQSValidator) ContainsLogin(login string) bool {
	v.RLock()
	defer v.RUnlock()

	idx := sort.SearchStrings(v.stargazerSlice, login)
	return idx < len(v.stargazerSlice) && v.stargazerSlice[idx] == login
}

// Validate validates use accessKey as a the github login,
// secretKey encrypted the data of accessKey.
func (v *SQSValidator) Validate(accessKey string, secretKey string) error {
	if !v.ContainsLogin(accessKey) {
		return errors.UserNotFound
	}

	return nil
}

// update stargazerSlice regularly
func (v *SQSValidator) updateRegularly() {
	ticker := time.NewTicker(updatePeriod)
	for {
		<-ticker.C
		v.updateStargazerSlice()
	}
}

func (v *SQSValidator) updateStargazerSlice() error {
	resp, err := http.Get(stargazerAPIURL)
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
	v.Lock()
	v.stargazerSlice = stargazersSlice
	v.Unlock()
	return nil
}
