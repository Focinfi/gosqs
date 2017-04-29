package log

import (
	"os"

	"github.com/Focinfi/gosqs/config"
	"github.com/Sirupsen/logrus"
)

func init() {
	if config.Config.Env.IsProduction() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetOutput(config.Config.LogOut)

	// Biz
	Biz.Level = logrus.DebugLevel
	Biz.Out = os.Stdout
}

// Internal for internal logic error
var Internal = logrus.New()

// DB for database error logger
var DB = logrus.New()

// Biz for user face logic logger
var Biz = logrus.New()

// Service for cluster error logger
var Service = logrus.New()

// Service for cluster error logger
var Cluster = logrus.New()

// InternalError for logic error
func InternalError(funcName string, message interface{}) {
	Internal.WithFields(logrus.Fields{"function_name": funcName}).Error(message)
}

// DBError log database error
func DBError(sql interface{}, err error, message interface{}) {
	DB.WithFields(logrus.Fields{"sql": sql, "error": err}).Error(message)
}

// LibError for lib error
func LibError(lib string, message interface{}) {
	DB.WithFields(logrus.Fields{"lib": lib}).Error(message)
}

// ThirdPartyServiceError for third-party node error
func ThirdPartyServiceError(thirdPartyService string, err error, message interface{}, params ...string) {
	DB.WithFields(logrus.Fields{"third_party_service": thirdPartyService, "error": err, "params": params}).Info(message)
}
