package storage

import (
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
)

var defaultKV models.KV
var defaultIncrementer models.Incrementer

func initKV() {
	var kv models.KV
	// etcd
	kv, err := etcd.NewKV()
	if err != nil {
		panic(err)
	}

	// goma
	// kv = gomap.New()

	// memcahed
	// kv, err := memcached.New()
	// if err != nil {
	// 	panic(err)
	// }

	defaultKV = kv
}

func initIncrementer() {
	incrementer, err := etcd.NewIncrementer()
	if err != nil {
		panic(err)
	}

	defaultIncrementer = incrementer
}

func init() {
	initKV()
	initIncrementer()
}
