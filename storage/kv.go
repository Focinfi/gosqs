package storage

// KV defines underlying key/value database
type KV interface {
	Get(key string) (string, bool)
	Put(key string, value string) error
	Append(key string, value string) error
	Remove(key string) error
}

type kv struct{}

func (k kv) Get(key string) (string, bool)         { return "", false }
func (k kv) Put(key string, value string) error    { return nil }
func (k kv) Append(key string, value string) error { return nil }
func (k kv) Remove(key string) error               { return nil }

var defaultKV = kv{}
