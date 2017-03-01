package storage

import "fmt"

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Append(key string, value string) error
	Remove(key string) error
}

type kv struct {
	data map[string]string
}

func (k *kv) Get(key string) (string, bool) {
	value, ok := k.data[key]
	return value, ok
}

func (k *kv) Put(key, value string) (err error) {
	track(func() {
		err = k.put(key, value)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Put")

	return
}

func (k *kv) Append(key, value string) (err error) {
	track(func() {
		err = k.append(key, value)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Append")
	return
}

func (k *kv) Remove(key string) (err error) {
	track(func() {
		err = k.remove(key)
		fmt.Println(JSONIndentFormat(k.data))
	}, "Remove")
	return
}

func (k *kv) put(key string, value string) error {
	k.data[key] = value
	return nil
}

func (k *kv) append(key string, value string) error {
	oldVal, ok := k.Get(key)
	if ok {
		k.data[key] = oldVal + value
	}

	k.data[key] = value
	return nil
}

func (k *kv) remove(key string) error {
	delete(k.data, key)
	return nil
}

var defaultKV = &kv{data: map[string]string{}}
