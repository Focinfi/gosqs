package storage

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

func (k *kv) Put(key string, value string) error {
	k.data[key] = value
	return nil
}

func (k *kv) Append(key string, value string) error {
	oldVal, ok := k.Get(key)
	if ok {
		k.data[key] = oldVal + value
	}

	k.data[key] = value
	return nil
}

func (k *kv) Remove(key string) error {
	delete(k.data, key)
	return nil
}

var defaultKV = &kv{}
