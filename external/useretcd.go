package external

import (
	"fmt"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/storage/etcd"
	"github.com/Focinfi/sqs/util/strconvutil"
)

const (
	userStoreIDKey = "sqs.users.id"
	userStoreKey   = "sqs.users"
)

func etcdUserStoreKey(uniqueID string) string {
	return fmt.Sprintf("%s.%s", userStoreKey, uniqueID)
}

// Etcd for a UserStore using etcdkv
type Etcd struct {
	*etcd.KV
	*etcd.Incrementer
}

// NewEtcd returns a new Etcd
func NewEtcd(kv *etcd.KV) *Etcd {
	return &Etcd{KV: kv, Incrementer: etcd.NewIncrementer(kv)}
}

// GetUserIDByUniqueID gets one user using the uniqueID
func (e *Etcd) GetUserIDByUniqueID(uniqueID string) (int64, error) {
	key := etcdUserStoreKey(uniqueID)
	val, err := e.Get(key)
	if err != nil {
		return -1, err
	}

	id, err := strconvutil.ParseInt64(val)
	if err != nil {
		return -1, errors.DataBroken(key, err)
	}

	return id, nil
}

// CreateUserByUniqueID creates a user with the given uniqueID
func (e *Etcd) CreateUserByUniqueID(uniqueID string) (int64, error) {
	_, err := e.GetUserIDByUniqueID(uniqueID)
	if err == errors.DataNotFound {
		id, err := e.Increment(userStoreIDKey, 1)
		if err != nil {
			return -1, err
		}

		key := etcdUserStoreKey(uniqueID)
		return id, e.Put(key, strconvutil.Int64toa(id))
	}

	if err != nil {
		return -1, err
	}

	return -1, errors.DuplicateUser
}
