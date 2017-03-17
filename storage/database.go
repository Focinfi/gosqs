package storage

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
	"github.com/Focinfi/sqs/storage/gomap"
	"github.com/Focinfi/sqs/storage/oncekv"
	"github.com/Focinfi/sqs/storage/redis"
)

var defaultKV models.KV
var defaultIncrementer models.Incrementer
var etcdIncrementer *etcd.Incrementer
var etcdKV *etcd.KV
var etcdWatcher *etcd.Watcher
var mapKV *gomap.KV
var onceKV *oncekv.KV
var redisPriorityList *redis.PriorityList

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
	etcdWatcher = etcd.NewWatcher(etcdKV)

	// etcd incrementer
	etcdIncrementer = etcd.NewIncrementer(etcdKV)
	defaultIncrementer = etcdIncrementer

	// mapkv
	mapKV = gomap.New()

	if kv, err := oncekv.NewKV(); err != nil {
		panic(err)
	} else {
		onceKV = kv
	}

	// redis
	if pl, err := redis.New(); err != nil {
		panic(err)
	} else {
		redisPriorityList = pl
	}
}
