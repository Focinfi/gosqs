package fmtutil

import (
	"encoding/json"
	"fmt"
)

// JSONIndentFormat returns JSON format with indent
func JSONIndentFormat(value interface{}) string {
	b, _ := json.MarshalIndent(value, "", "  ")
	return string(b)
}

const defaultFormat = "[%s] "

// Format contains a prefix
type Format struct {
	Prefix string
}

// NewFormat returns a new Format
func NewFormat(prefix string) *Format {
	return &Format{
		Prefix: fmt.Sprintf(defaultFormat, prefix),
	}
}

// Sprintf adds prefix before format
func (f Format) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(f.Prefix+format, a...)
}

// Sprintln adds prefix before format
func (f Format) Sprintln(args ...interface{}) string {
	return f.Prefix + fmt.Sprintln(args...)
}

// Println adds prefix before format
func (f Format) Println(args ...interface{}) {
	fmt.Println(args...)
}
