package config

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/jinzhu/configor"
)

var root string

// Root returns the root path of oncekv
func Root() string {
	return root
}

// Env for application environment
type Env string

// IsProduction returns if the env equals to production
func (e Env) IsProduction() bool {
	return e == "production"
}

// IsDevelop returns if the env equals to develop
func (e Env) IsDevelop() bool {
	return e == "develop"
}

// IsTest returns if the env equals to develop
func (e Env) IsTest() bool {
	return e == "test"
}

// Config for config
var Config = struct {
	Env  Env `default:"develop" env:"SQS_ENV"`
	Root string

	// etcd addrs and the the meta data key
	EtcdEndpoints []string `default:"['127.0.0.1:2379']" env:"SQS_ETCD_ADDRS"`

	// master
	DefaultMasterAddress string `default:"127.0.0.1:5446" env:"DEFAULT_MESSAGE_ADDRESS"`

	// node
	PullMessageCount      int `default:"10" env:"PULL_MESSAGE_COUNT"`
	MaxMessageIDRangeSize int `default:"10" env:"MAX_MESSAGE_ID_RANGE_SIZE"`

	// message
	MaxTryMessageCount int `default:"5" env:"MAX_TRY_MESSAGE_COUNT"`

	// admin
	AdminAddr string `default:"127.0.0.1:54460" env:"SQS_ADMIN_ADDR"`

	// TODO: choose log collector
	LogOut io.Writer
}{}

func init() {
	if r := os.Getenv("GOPATH"); r != "" {
		root = path.Join(r, "src", "github.com", "Focinfi", "oncekv")
	} else {
		panic("sqs: envroinment param $GOPATH not set")
	}

	err := configor.Load(&Config, path.Join(root, "config", "config.json"))
	if err != nil {
		panic(err)
	}

	fmt.Println(Config)
}
