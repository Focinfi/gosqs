package config

import (
	"io"
	"os"

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

// Envroinment for application envroinment
type Envroinment string

// IsProduction returns if the env equals to production
func (e Envroinment) IsProduction() bool {
	return e == production
}

// IsDevelop returns if the env equals to develop
func (e Envroinment) IsDevelop() bool {
	return e == develop
}

// IsTest returns if the env equals to develop
func (e Envroinment) IsTest() bool {
	return e == develop
}

var env = Envroinment(develop)

// Env returns the env
func Env() Envroinment {
	return env
}

// Configuration defines configuration
type Configuration struct {
	ClientControlTimeoutSecond int64
	ClientDefaultPriority      int
	MaxConsumerSize            int
	MaxPushWorkCount           int
	LogOut                     io.Writer
}

func newDefaultConfig() Configuration {
	return Configuration{
		ClientControlTimeoutSecond: 3,
		MaxConsumerSize:            10,
		ClientDefaultPriority:      10,
		MaxPushWorkCount:           4,
		LogOut:                     os.Stdout,
	}
}

// Config returns the Configuration based on envroinment
func Config() Configuration {

	switch env {
	case productionEnv:
		return Configuration{
			ClientControlTimeoutSecond: 300,
			MaxConsumerSize:            1000000,
			ClientDefaultPriority:      100,
			MaxPushWorkCount:           16,
			LogOut:                     os.Stdout,
		}
	case developEnv:
		return newDefaultConfig()
	default:
		return newDefaultConfig()
	}
}

func init() {
	if e := os.Getenv("SQS_ENV"); e != "" {
		env = Envroinment(e)
	}

	if Env().IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
}
