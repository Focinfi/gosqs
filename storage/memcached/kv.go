package memcached

import (
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/log"
	"github.com/bradfitz/gomemcache/memcache"
)

// KV uses memcache as backend
type KV struct {
	db *memcache.Client
}

// Get returns a the value of the key
func (kv *KV) Get(key string) (string, bool) {
	item, err := kv.db.Get(key)
	if err == memcache.ErrCacheMiss {
		return "", false
	}

	if err != nil {
		log.DB.Errorln(err)
	}

	if item == nil {
		return "", false
	}

	return string(item.Value), true
}

// Put puts the value for the key
func (kv *KV) Put(key string, value string) error {
	return kv.db.Set(&memcache.Item{Key: key, Value: []byte(value)})
}

// Delete deletes the item for the given key
func (kv *KV) Delete(key string) error {
	return kv.db.Delete(key)
}

// New returns a new KV client
func New() (*KV, error) {
	client := memcache.New(config.Config().MemcachedEndpoints...)
	return &KV{db: client}, nil
}
