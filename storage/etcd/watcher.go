package etcd

import (
	"context"

	"github.com/Focinfi/sqs/log"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// Watcher watches a key
type Watcher struct {
	kv *KV
}

// Watch watches a key
func (w *Watcher) Watch(key string) <-chan string {
	ch := make(chan string)

	watchChan := w.kv.db.Watch(context.Background(), key)

	go func() {
		res := <-watchChan
		log.DB.Infof("etcd watch: %#v, %#v\n", res, res.Events)
		if res.Canceled {
			ch <- ""
			return
		}

		for i := len(res.Events) - 1; i >= 0; i-- {
			event := res.Events[i]
			switch event.Type {
			case mvccpb.DELETE:
				ch <- ""
				return
			case mvccpb.PUT:
				ch <- string(event.Kv.Value)
				return
			}
		}

		ch <- ""
	}()

	return ch
}

// Close closes the watch channel
func (w *Watcher) Close() error {
	return w.kv.db.Watcher.Close()
}

// NewWatcher returns a new watcher
func NewWatcher(kv *KV) *Watcher {
	return &Watcher{kv: kv}
}
