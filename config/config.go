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
	ClientDefaultPriority      int
	MaxConsumerSize            int
	MaxPushWorkCount           int
}

// Config returns the Configuration based on envroinment
func Config() Configuration {
	env := os.Getenv("sqs-env")
	if env == "" {
		env = developEnv
	}

	switch env {
	case productionEnv:
		return Configuration{
			ClientControlTimeoutSecond: 300,
			MaxConsumerSize:            1000000,
			ClientDefaultPriority:      100,
			MaxPushWorkCount:           16,
		}
	case developEnv:
		return Configuration{
			ClientControlTimeoutSecond: 3,
			MaxConsumerSize:            10,
			ClientDefaultPriority:      10,
			MaxPushWorkCount:           4,
		}
	default:
		return Configuration{
			ClientControlTimeoutSecond: 3,
			MaxConsumerSize:            10,
			ClientDefaultPriority:      10,
			MaxPushWorkCount:           4,
		}
	}
}
