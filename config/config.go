package config

import (
	"os"
)

const (
	testEnv       = "test"
	developEnv    = "develop"
	productionEnv = "production"
)

// Configuration defines configuration
type Configuration struct {
	ClientControlTimeoutSecond int64
}

// Config returns the Configuration based on envroinment
func Config() Configuration {
	env := os.Getenv("sqs-env")
	if env == "" {
		env = developEnv
	}

	switch env {
	case productionEnv:
		return Configuration{ClientControlTimeoutSecond: 300}
	case developEnv:
		return Configuration{ClientControlTimeoutSecond: 3}
	default:
		return Configuration{ClientControlTimeoutSecond: 3}
	}
}
