package config

import (
	"io"
	"os"

	"time"

	"github.com/gin-gonic/gin"
)

const (
	testEnv       = "test"
	developEnv    = "develop"
	productionEnv = "production"
)

const (
	test       = "test"
	develop    = "develop"
	production = "production"
)

// Environment for application environment
type Environment string

// IsProduction returns if the env equals to production
func (e Environment) IsProduction() bool {
	return e == production
}

// IsDevelop returns if the env equals to develop
func (e Environment) IsDevelop() bool {
	return e == develop
}

// IsTest returns if the env equals to develop
func (e Environment) IsTest() bool {
	return e == test
}

var env = Environment(develop)

// Env returns the env
func Env() Environment {
	return env
}

// Configuration defines configuration
type Configuration struct {
	ClientControlTimeoutSecond int64
	ClientDefaultPriority      int
	MaxConsumerSize            int
	MaxPushWorkCount           int
	MaxRetryConsumerSeconds    int
	LogOut                     io.Writer
	EtcdEndpoints              []string
	MemcachedEndpoints         []string
	RedisAdrr                  string
	RedisPwd                   string
	MaxMessgeIDRangeSize       int
	MaxTryMessageCount         int
	OncekvMetaRefreshPeroid    time.Duration
	IdealKVResponseDuration    time.Duration
	PullMessageCount           int
	DefaultMasterAddress       string
}

func newDefaultConfig() Configuration {
	return Configuration{
		ClientControlTimeoutSecond: 3,
		MaxConsumerSize:            10,
		ClientDefaultPriority:      10,
		MaxPushWorkCount:           4,
		LogOut:                     os.Stdout,
		EtcdEndpoints:              []string{"localhost:2379"},
		MemcachedEndpoints:         []string{"localhost:11211"},
		RedisAdrr:                  "localhost:6379",
		RedisPwd:                   "",
		MaxRetryConsumerSeconds:    3,
		MaxMessgeIDRangeSize:       10,
		MaxTryMessageCount:         3,
		OncekvMetaRefreshPeroid:    time.Second,
		IdealKVResponseDuration:    time.Millisecond * 50,
		PullMessageCount:           5,
		DefaultMasterAddress:       "127.0.0.1:5446",
	}
}

// Config returns the Configuration based on envroinment
func Config() Configuration {

	switch env {
	case productionEnv:
		return Configuration{
			ClientControlTimeoutSecond: 300,
			MaxConsumerSize:            1000000,
			MaxRetryConsumerSeconds:    10,
			ClientDefaultPriority:      100,
			MaxPushWorkCount:           16,
			LogOut:                     os.Stdout,
			MaxMessgeIDRangeSize:       100,
			MaxTryMessageCount:         10,
			OncekvMetaRefreshPeroid:    time.Second,
			IdealKVResponseDuration:    time.Millisecond * 50,
			PullMessageCount:           10,
		}
	case developEnv:
		return newDefaultConfig()
	default:
		return newDefaultConfig()
	}
}

func init() {
	if e := os.Getenv("SQS_ENV"); e != "" {
		env = Environment(e)
	}

	if Env().IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
}
