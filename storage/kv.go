package storage

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/memcached"
)

var defaultKV models.KV

func init() {
	// kv, err := etcd.New()
	// if err != nil {
	// 	log.Internal.Panic(err)
	// }

	// defaultKV = kv
	defaultKV, _ = memcached.New()
}
