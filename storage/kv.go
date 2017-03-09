package storage

import (
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/storage/etcd"
)

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Delete(key string) error
}

var defaultKV KV

func init() {
	kv, err := etcd.New()
	if err != nil {
		log.Internal.Panic(err)
	}

	defaultKV = kv
}
