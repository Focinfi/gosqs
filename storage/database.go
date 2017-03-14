package storage

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
	"github.com/Focinfi/sqs/storage/gomap"
	"github.com/Focinfi/sqs/storage/memcached"
	"github.com/Focinfi/sqs/storage/redis"
)

var defaultKV models.KV
var defaultIncrementer models.Incrementer
var etcdIncrementer models.Incrementer
var etcdKV models.KV
var etcdWatcher models.Watcher
var memcachedKV models.KV
var redisPriorityList models.PriorityList

func init() {
	// gomap
	defaultKV = gomap.New()

	// etcd kv
	if kv, err := etcd.NewKV(); err != nil {
		panic(err)
	} else {
		etcdKV = kv
	}

	// etcd watcher
	if watcher, err := etcd.NewWatcher(); err != nil {
		panic(err)
	} else {
		etcdWatcher = watcher
	}

	// etcd incrementer
	if incrementer, err := etcd.NewIncrementer(); err != nil {
		panic(err)
	} else {
		etcdIncrementer = incrementer
		defaultIncrementer = incrementer
	}

	// memcahed
	if kv, err := memcached.New(); err != nil {
		panic(err)
	} else {
		memcachedKV = kv
	}

	// redis
	if pl, err := redis.New(); err != nil {
		panic(err)
	} else {
		redisPriorityList = pl
	}
}
