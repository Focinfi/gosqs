package log

import (
	"log"
	"os"
)

// Log for logging
type Log struct {
	*log.Logger
}

var logger = Log{log.New(os.Stdout, "sqs", log.Ltime)}
