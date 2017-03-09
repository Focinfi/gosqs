package storage

import "github.com/Focinfi/sqs/storage/memcached"

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Delete(key string) error
}

var defaultKV KV

func init() {
	// kv, err := etcd.New()
	// if err != nil {
	// 	log.Internal.Panic(err)
	// }

	// defaultKV = kv
	defaultKV, _ = memcached.New()
}
