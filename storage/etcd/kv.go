package etcd

import (
	"context"

	"github.com/Focinfi/gosqs/errors"
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
func (kv *KV) Get(key string) (string, error) {
	res, err := kv.db.Get(context.Background(), key)

	if err != nil {
		return "", err
	}

	if len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value), nil
	}

	return "", errors.DataNotFound
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

// NewKV returns a new EtcdKV
func NewKV() (*KV, error) {
	cli, err := New()

	if err != nil {
		return nil, err
	}

	return &KV{db: cli}, nil
}
