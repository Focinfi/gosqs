package gomap

import (
	"fmt"
	"sync"
)

// KV use a map as a K/V database
type KV struct {
	sync.RWMutex
	data map[string]string
}

// New returns a new KV
func New() *KV {
	return &KV{data: map[string]string{}}
}

// Get gets the value and existence for the key
func (k *KV) Get(key string) (string, bool) {
	k.RLock()
	defer k.RUnlock()

	value, ok := k.data[key]
	return value, ok
}

// Put puts the value for the key
func (k *KV) Put(key, value string) (err error) {
	track(func() {
		err = k.put(key, value)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Put")

	return
}

// Append appends sthe value for the key
func (k *KV) Append(key, value string) (err error) {
	track(func() {
		err = k.append(key, value)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Append")
	return
}

// Delete delete the item for the key
func (k *KV) Delete(key string) (err error) {
	track(func() {
		err = k.remove(key)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Delete")
	return
}

func (k *KV) put(key string, value string) error {
	k.Lock()
	defer k.Unlock()

	k.data[key] = value
	return nil
}

func (k *KV) append(key string, value string) error {
	k.Lock()
	defer k.Unlock()

	oldVal, ok := k.Get(key)
	if ok {
		k.data[key] = oldVal + value
	}

	k.data[key] = value
	return nil
}

func (k *KV) remove(key string) error {
	k.Lock()
	defer k.Unlock()

	delete(k.data, key)
	return nil
}
