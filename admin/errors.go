package admin

import (
	"errors"
)

// ErrDuplicateQueue error for duplicate queue
var ErrDuplicateQueue = errors.New("admin: duplicate queue")
