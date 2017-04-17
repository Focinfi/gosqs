package gomap

import (
	"fmt"

	"github.com/Focinfi/sqs/log"
)

const startFormat = "-Start--%s-------\n"
const endFormat = "-End----%s-------\n"

func track(fn func(), a ...interface{}) {
	log.DB.Infof(startFormat, fmt.Sprint(a...))
	fn()
	log.Internal.Printf(endFormat, fmt.Sprint(a...))
}
