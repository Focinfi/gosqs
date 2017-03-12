package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// Watcher watches a key
type Watcher struct {
	agent *clientv3.Client
}

// Watch watches a key
func (w *Watcher) Watch(key string) <-chan string {
	ch := make(chan string)

	watchChan := w.agent.Watch(context.Background(), key)

	go func() {
		res := <-watchChan
		if res.Canceled {
			ch <- ""
			return
		}

		for i := len(res.Events); i >= 0; i-- {
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

	return nil
}

// NewWatcher returns a new watcher
func NewWatcher() (*Watcher, error) {
	cli, err := New()
	if err != nil {
		return nil, err
	}

	return &Watcher{agent: cli}, nil
}
