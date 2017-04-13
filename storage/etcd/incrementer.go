package etcd

import (
	"context"

	"fmt"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/util/strconvutil"
	"github.com/coreos/etcd/clientv3"
)

// Incrementer implements the models.Incrementer
type Incrementer struct {
	kv *KV
}

// Increment try to increment the value of the key
func (inc *Incrementer) Increment(key string, size int) (int64, error) {
	curIDVal, getErr := inc.kv.Get(key)

	// try set the value to size
	if getErr == errors.DataNotFound {
		maxIDVal := fmt.Sprintf("%d", size)
		_, err := inc.kv.db.Txn(context.TODO()).
			If(clientv3.CreateRevision(key)).
			Then(clientv3.OpPut(key, maxIDVal)).
			Commit()

		if err != nil {
			return -1, err
		}

		return int64(size), nil
	}

	if getErr != nil {
		return -1, getErr
	}

	curID, getErr := strconvutil.ParseInt64(curIDVal)
	if getErr != nil {
		return -1, errors.DataBroken(key, getErr)
	}

	maxID := curID + int64(size)
	maxIDVal := fmt.Sprintf("%d", maxID)

	_, getErr = inc.kv.db.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.Value(key), "=", curIDVal)).
		Then(clientv3.OpPut(key, maxIDVal)).
		Commit()

	if getErr != nil {
		return -1, getErr
	}

	return maxID, nil
}

// NewIncrementer returns a new Incrementer
func NewIncrementer(kv *KV) *Incrementer {
	return &Incrementer{kv: kv}
}
