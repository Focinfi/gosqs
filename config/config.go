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

	// secret
	BaseSecret string `default:"aslMep28nkvfiYTuWuF7OIp2A5sMb5ewOu2UwO/1PEI=" env:"SQS_BASE_SECRET"`

	// auth
	UserGithubLoginKey string `default: "user_github_login_key" env:"SQS_USER_GITHUB_LOGIN_KEY"`

	// etcd addrs and the the meta data key
	EtcdEndpoints []string `default:"['127.0.0.1:2379']" env:"SQS_ETCD_ADDRS"`

	// master
	DefaultMasterAddress string `default:"127.0.0.1:5446" env:"SQS_DEFAULT_MESSAGE_ADDRESS"`

	// node
	PullMessageCount      int `default:"10" env:"SQS_PULL_MESSAGE_COUNT"`
	MaxMessageIDRangeSize int `default:"10" env:"SQS_MAX_MESSAGE_ID_RANGE_SIZE"`

	// message
	MaxTryMessageCount int `default:"5" env:"SQS_MAX_TRY_MESSAGE_COUNT"`

	// admin
	AdminAddr string `default:"127.0.0.1:54460" env:"SQS_ADMIN_ADDR"`

	// TODO: choose log collector
	LogOut io.Writer

	// queue
	MaxQueueCountPerUser int `default:"10" env:"SQS_QUEUE_COUNT_PRE_USER"`

	SQLDB struct {
		Adapter  string `default:"mysql"`
		Name     string `default:"sqs" env:"MYSQL_INSTANCE_NAME"`
		Host     string `default:"127.0.0.1" env:"MYSQL_PORT_3306_TCP_ADDR"`
		Port     string `default:"3306" env:"MYSQL_PORT_3306_TCP_PORT"`
		User     string `default:"sqs" env:"MYSQL_USERNAME"`
		Password string `default:"" env:"MYSQL_PASSWORD"`
		Protocol string `default:"tcp" env:"MYSQL_PORT_3306_TCP_PROTO"`
	}

	Email struct {
		SMTP     string `default:"smtpdm.aliyun.com" env:"SQS_EMAIL_SMTP"`
		Port     int    `default:"25" env:"SQS_EMAIL_PORT"`
		From     string `default:"noreply@sqsadmin.club" env:"SQS_EMAIL_FROM"`
		User     string `default:"noreply@sqsadmin.club" env:"SQS_EMAIL_USER"`
		Password string `env:"SQS_EMAIL_PASSWORD"`
	}
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
