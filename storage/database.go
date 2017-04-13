package storage

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
	"github.com/Focinfi/sqs/storage/gomap"
)

var defaultKV models.KV
var defaultIncrementer models.Incrementer

// EtcdKV uses etcd as a key/value storage
var EtcdKV *etcd.KV
var etcdIncrementer *etcd.Incrementer
var etcdWatcher *etcd.Watcher
var mapKV *gomap.KV

func init() {
	// map
	defaultKV = gomap.New()

	// etcd db
	if kv, err := etcd.NewKV(); err != nil {
		panic(err)
	} else {
		EtcdKV = kv
	}

	// etcd watcher
	etcdWatcher = etcd.NewWatcher(EtcdKV)

	// etcd incrementer
	etcdIncrementer = etcd.NewIncrementer(EtcdKV)
	defaultIncrementer = etcdIncrementer

	// mapkv
	mapKV = gomap.New()

	//if kv, err := oncekv.NewKV(); err != nil {
	//	panic(err)
	//} else {
	//	onceKV = kv
	//}
}
