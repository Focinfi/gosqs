package gomap

import (
	"encoding/json"
	"fmt"

	"github.com/Focinfi/sqs/log"
)

const startFormat = "-Start--%s-------\n"
const endFormat = "-End----%s-------\n"

func track(fn func(), a ...interface{}) {
	log.Biz.Printf(startFormat, fmt.Sprint(a...))
	fn()
	log.Biz.Printf(endFormat, fmt.Sprint(a...))
}

func trackf(fn func(), format string, a ...interface{}) {
	log.Biz.Printf(startFormat, fmt.Sprintf(format, a...))
	fn()
	log.Biz.Printf(endFormat, fmt.Sprintf(format, a...))
}

// JSONIndentFormat returns JSON format with indent
func JSONIndentFormat(value interface{}) string {
	b, _ := json.MarshalIndent(value, "", "  ")
	return string(b)
}
