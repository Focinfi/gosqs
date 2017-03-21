package oncekv

import (
	"time"

	"github.com/Focinfi/oncekv/client"
	"github.com/Focinfi/sqs/errors"
)

const requestTimeout = time.Millisecond * 300

// KV for kv storage
type KV struct {
	db *client.KV
}

// Get get the value of the key
func (kv *KV) Get(key string) (string, error) {

	val, err := kv.db.Get(key)
	if err == client.ErrDataNotFound {
		return "", errors.DBNotFound
	}

	if err != nil {
		return "", errors.NewInternalWrap(err)
	}

	return val, nil
}

// Put put the key/value pair
func (kv *KV) Put(key, value string) error {
	return kv.db.Put(key, value)
}

// Delete deletes the key
func (kv *KV) Delete(key string) error {
	return nil
}

// NewKV returns a new kv
func NewKV() (*KV, error) {
	db, err := client.DefaultKV()
	if err != nil {
		return nil, err
	}

	return &KV{db: db}, nil
}
