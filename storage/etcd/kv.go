package etcd

import (
	"context"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/log"
	"github.com/coreos/etcd/clientv3"
)

// KV uses etcd as a KV
type KV struct {
	db *clientv3.Client
}

// Close close db
func (kv *KV) Close() {
	kv.db.Close()
}

// Get gets the value for the key
func (kv *KV) Get(key string) (string, bool) {

	res, err := kv.db.Get(context.Background(), key)

	if err != nil {
		log.DB.Error(err)
		return "", false
	}

	if len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value), true
	}

	return "", false
}

// Put puts key/value
func (kv *KV) Put(key string, value string) error {
	_, err := kv.db.Put(context.Background(), key, value)
	return err
}

// Delete deletes key/value
func (kv *KV) Delete(key string) error {
	_, err := kv.db.Delete(context.Background(), key)
	return err
}

// New returns a new EtcdKV
func New() (*KV, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Config().EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return &KV{db: cli}, nil
}
