package strconvutil

import (
	"fmt"
	"strconv"
)

// ParseInt64 parses the s in base 10 and 64 bitSize
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Int64toa return a int64 string value
func Int64toa(i int64) string {
	return fmt.Sprintf("%d", i)
}
