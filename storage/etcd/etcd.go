package etcd

import (
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/coreos/etcd/clientv3"
)

// New returns a new etcd client
func New() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Config().EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return cli, nil
}
