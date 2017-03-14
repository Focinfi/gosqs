package etcd

import (
	"context"

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

// Txn calls kv.db.Txn
func (kv *KV) Txn(ctx context.Context) clientv3.Txn {
	return kv.db.Txn(ctx)
}

// NewKV returns a new EtcdKV
func NewKV() (*KV, error) {
	cli, err := New()

	if err != nil {
		return nil, err
	}

	return &KV{db: cli}, nil
}
