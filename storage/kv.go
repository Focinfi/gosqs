package storage

import (
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
)

var defaultKV models.KV

func init() {
	var kv models.KV
	// etcd
	kv, err := etcd.NewKV()
	if err != nil {
		log.Internal.Panic(err)
	}

	// goma
	// kv = gomap.New()

	// memcahed
	// kv, err := memcached.New()

	if err != nil {
		log.DB.Panic(err)
	}

	defaultKV = kv
}
