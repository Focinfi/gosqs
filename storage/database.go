package storage

import (
	"fmt"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/storage/etcd"
	"github.com/Focinfi/sqs/storage/gomap"
	"github.com/Focinfi/sqs/storage/oncekv"
)

// ClusterMetaKV for nodes cluster
var ClusterMetaKV models.KV
var sqsMetaKV models.KV
var messageKV models.KV

// different backend
var mapKV *gomap.KV
var onceKV *oncekv.KV
var etcdKV *etcd.KV
var etcdIncrementer *etcd.Incrementer

func init() {
	//mapKV
	mapKV = gomap.New()

	// etcd db
	if kv, err := etcd.NewKV(); err != nil {
		panic(err)
	} else {
		etcdKV = kv
	}
	// incrementer
	etcdIncrementer = etcd.NewIncrementer(etcdKV)

	// oncekv
	if kv, err := oncekv.NewKV(); err != nil {
		panic(err)
	} else {
		onceKV = kv
	}

	if config.Config.Env.IsProduction() {
		ClusterMetaKV = etcdKV
		sqsMetaKV = etcdKV
		messageKV = onceKV
	} else if config.Config.Env.IsDevelop() {
		ClusterMetaKV = etcdKV
		sqsMetaKV = etcdKV
		messageKV = mapKV
	} else if config.Config.Env.IsTest() {
		ClusterMetaKV = mapKV
		sqsMetaKV = mapKV
		messageKV = mapKV
	} else {
		panic(fmt.Sprintf("env '%s' is not allowed", config.Config.Env))
	}
}
