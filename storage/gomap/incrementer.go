package gomap

import (
	"sync"

	"github.com/Focinfi/sqs/util/strconvutil"
)

// Incrementer implements the models.Incrementer interface
type Incrementer struct {
	sync.Mutex
	kv *KV
}

func NewIncrementer(kv *KV) *Incrementer {
	return &Incrementer{kv: kv}
}

//	Increment(key string, number int) (result int64, err error)
func (i *Incrementer) Increment(key string, size int) (int64, error) {
	i.Lock()
	defer i.Unlock()

	val, err := i.kv.Get(key)
	if err != nil {
		return -1, err
	}

	cur, err := strconvutil.ParseInt64(val)
	if err != nil {
		return -1, err
	}

	res := cur + int64(size)
	if i.kv.Put(key, strconvutil.Int64toa(res)); err != nil {
		return -1, err
	}

	return res, nil
}
