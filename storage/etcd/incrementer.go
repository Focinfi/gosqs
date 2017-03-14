package etcd

import (
	"context"

	"fmt"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/util"
	"github.com/coreos/etcd/clientv3"
)

// Incrementer implements the models.Incrementer
type Incrementer struct {
	db *KV
}

// Increment try to increment the value of the key
func (inc *Incrementer) Increment(key string, size int) (int64, error) {
	curIDVal, ok := inc.db.Get(key)

	// try set the value to size
	if !ok {
		maxIDVal := fmt.Sprintf("%d", size)
		_, err := inc.db.Txn(context.TODO()).
			If(clientv3.CreateRevision(key)).
			Then(clientv3.OpPut(key, maxIDVal)).
			Commit()

		if err != nil {
			return -1, err
		}

		return int64(size), nil
	}

	curID, err := util.ParseInt64(curIDVal)
	if err != nil {
		return -1, errors.DataBroken(key, err)
	}

	maxID := curID + int64(size)
	maxIDVal := fmt.Sprintf("%d", maxID)

	_, err = inc.db.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.Value(key), "=", curIDVal)).
		Then(clientv3.OpPut(key, maxIDVal)).
		Commit()

	if err != nil {
		return -1, err
	}

	return maxID, nil
}

// NewIncrementer returns a new Incrementer
func NewIncrementer(kv *KV) (*Incrementer, error) {
	return &Incrementer{db: kv}, nil
}
