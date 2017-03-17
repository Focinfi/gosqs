package oncekv

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
)

const requestTimeout = time.Second * 2

// KV for kv storage
type KV struct {
	cli *Client
}

type kvParams struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Get get the value of the key
func (kv *KV) Get(key string) (string, bool) {
	val, err := kv.getFromCache(key)
	if err == nil {
		return val, true
	}
	log.DB.Error(err)

	val, err = kv.getFromDB(key)
	if err != nil {
		return "", false
	}

	return val, true
}

// Put put key/value pair
func (kv *KV) Put(key string, value string) error {
	dbs := []string{}
	copy(dbs, kv.cli.dbs)
	var done bool
	var mux sync.Mutex

	err := errors.New("failed to set")
	success := make(chan bool)

	for _, db := range dbs {
		go func(url string) {
			err = kv.setValue(key, value, url)

			if err != nil {
				err = errors.DataLost(err.Error())
				return
			}

			mux.Lock()
			if !done {
				done = true
				go kv.cli.setFastDB(url)
				mux.Unlock()
			} else {
				mux.Unlock()
				return
			}

			success <- true
		}(db)
	}

	select {
	case <-time.After(requestTimeout):
	case <-success:
		return nil
	}

	return err
}

// Delete delete the key
func (kv *KV) Delete(key string) error {
	// TODO: remove old key for save db space
	return nil
}

func (kv *KV) setValue(key string, value string, url string) error {
	// http.Post(fmt.Sprintf("http://%s/key", url), "application-type/json", bytes.NewReader(b))
	return nil
}

func (kv *KV) getFromCache(key string) (string, error) {
	url := kv.cli.fastCache
	if url == "" {
		return kv.fetchFromCache(key)
	}

	idealDuration := config.Config().IdealKVResponseDuration
	val, duration, err := kv.fetchValue(key, idealDuration, url)
	if err != nil {
		return kv.fetchFromCache(key)
	}

	if val == "" {
		return "", errors.DataLost(key)
	}

	// for update kv.cli.fastCache
	if duration > idealDuration {
		go func() {
			if _, err := kv.fetchFromCache(key); err != nil {
				log.DB.Error(err)
			}
		}()
	}

	return val, err
}

func (kv *KV) getFromDB(key string) (string, error) {
	dbs := []string{}
	copy(dbs, kv.cli.dbs)
	var done bool
	var mux sync.Mutex

	err := errors.DataLost(key)
	data := make(chan string)

	for _, db := range dbs {
		go func(url string) {
			val, _, err := kv.fetchValue(key, requestTimeout, url)
			if err != nil {
				return
			}

			if val == "" {
				err = errors.DataLost(key)
				return
			}

			mux.Lock()
			if !done {
				done = true
				go kv.cli.setFastDB(url)
				mux.Unlock()
			} else {
				mux.Unlock()
				return
			}

			data <- val
		}(db)
	}

	select {
	case <-time.After(requestTimeout):
	case value := <-data:
		return value, nil
	}

	return "", err
}

func (kv *KV) tryParseResponse(readCloser io.ReadCloser, key string) (string, error) {
	defer readCloser.Close()

	b, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	param := &kvParams{}
	if err := json.Unmarshal(b, param); err != nil {
		return "", err
	}

	if param.Key != key {
		return "", errors.DataLost(key)
	}

	if param.Value == "" {
		return "", errors.DataBroken(key, errors.New("value is nil"))
	}

	return param.Value, nil
}

func (kv *KV) fetchValue(key string, timeout time.Duration, url string) (value string, duration time.Duration, err error) {
	begin := time.Now()
	resChan := make(chan *http.Response)
	errChan := make(chan error)

	go func() {
		res, err := http.Get(fmt.Sprintf("%s/sqs/message/%s", url, key))
		resChan <- res
		errChan <- err
	}()

	select {
	case <-time.After(config.Config().IdealKVResponseDuration):
		return "", -1, errors.DBQueryTimeout(url, key)
	case err := <-errChan:
		log.DB.Error(err)
		return "", -1, errors.NewInternalErr(err.Error())
	case res := <-resChan:
		if res.StatusCode == http.StatusOK {
			if val, err := kv.tryParseResponse(res.Body, key); err != nil {
				return val, time.Now().Sub(begin), err
			}
		}
	}

	return "", -1, errors.DataLost(key)
}

// try all caching urls, set the fastCache
func (kv *KV) fetchFromCache(key string) (string, error) {
	caches := []string{}
	copy(caches, kv.cli.caches)
	durations := make([]time.Duration, len(caches))
	var fetched bool
	var value string
	var mux sync.Mutex
	var err error
	var wg sync.WaitGroup

	for i, cache := range caches {
		go func(index int, url string) {
			wg.Add(1)
			defer wg.Done()
			val, duration, err := kv.fetchValue(key, requestTimeout, url)
			if err != nil {
				log.DB.Error(err)
				durations[index] = requestTimeout
				return
			}

			if val != "" {
				err = errors.DataLost(key)
				log.DB.Error(err)
				durations[index] = requestTimeout
				return
			}

			durations[index] = duration

			mux.Lock()
			if !fetched {
				fetched = true
				value = val
			}
			mux.Unlock()
		}(i, cache)
	}

	wg.Wait()

	go func() {
		var min = requestTimeout
		index := -1
		for i, duration := range durations {
			if duration <= min {
				min = duration
				index = i
			}
		}
		if index >= 0 {
			kv.cli.setFastCache(caches[index])
		}
	}()

	return value, err
}

// NewKV returns a new KV
func NewKV() (*KV, error) {
	cli, err := New()
	if err != nil {
		return nil, err
	}

	return &KV{cli: cli}, nil
}
