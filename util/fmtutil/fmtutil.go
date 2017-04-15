package fmtutil

import "encoding/json"

// JSONIndentFormat returns JSON format with indent
func JSONIndentFormat(value interface{}) string {
	b, _ := json.MarshalIndent(value, "", "  ")
	return string(b)
}
