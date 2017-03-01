package storage

import (
	"encoding/json"
	"fmt"
)

const startFormat = "-Start--%s-------\n"
const endFormat = "-End----%s-------\n"

func track(fn func(), a ...interface{}) {
	fmt.Printf(startFormat, fmt.Sprint(a...))
	fn()
	fmt.Printf(endFormat, fmt.Sprint(a...))
}

func trackf(fn func(), format string, a ...interface{}) {
	fmt.Printf(startFormat, fmt.Sprintf(format, a...))
	fn()
	fmt.Printf(endFormat, fmt.Sprintf(format, a...))
}

// JSONIndentFormat returns JOSN format with indent
func JSONIndentFormat(value interface{}) string {
	b, _ := json.MarshalIndent(value, "", "  ")
	return string(b)
}
