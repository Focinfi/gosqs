package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/external"
	"github.com/Focinfi/sqs/util/token"
)

const (
	userGithubLoginKey = "userGithubLogin"
	apiURL             = "https://api.github.com/repos/Focinfi/sqs/stargazerSlice"
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

// GithubAuther auth by github stargazers
type GithubAuther struct {
	sync.Mutex
	stargazerSlice []string
}

// NewGithubAuther allocates and return a new GithubAuther
func NewGithubAuther() *GithubAuther {
	return &GithubAuther{
		stargazerSlice: []string{},
	}
}

// Start start the background services
func (auther *GithubAuther) Start() {
	go auther.updateStargazerSlice()
}

// Auth auth use accessKey as a the github login,
// secretKey encrypted the data of accessKey.
func (auther *GithubAuther) Auth(accessKey string, secretKey string) (int64, error) {
	idx := sort.SearchStrings(auther.stargazerSlice, accessKey)
	if idx > len(auther.stargazerSlice) || auther.stargazerSlice[idx] != accessKey {
		return -1, errors.UserNotFound
	}

	params, err := token.Default.Verify(secretKey, config.Config.BaseSecret)
	if err != nil {
		return -1, errors.UserAuthError(err.Error())
	}

	if !reflect.DeepEqual(secretKey, params[userGithubLoginKey]) {
		return -1, errors.UserAuthError("broken secrect_key")
	}

	return external.GetUserIDByUniqueID(accessKey)
}

// update stargazerSlice regularly
func (auther *GithubAuther) updateRegularly() {
	ticker := time.NewTicker(updatePeriod)
	for {
		<-ticker.C
		auther.updateStargazerSlice()
	}
}

func (auther *GithubAuther) updateStargazerSlice() error {
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respParam := stargazers{}
	if err := json.Unmarshal(respBytes, respParam); err != nil {
		return err
	}

	stargazersSlice := respParam.toSlice()
	sort.StringSlice(stargazersSlice).Sort()
	auther.Lock()
	auther.stargazerSlice = stargazersSlice
	auther.Unlock()
	return nil
}
