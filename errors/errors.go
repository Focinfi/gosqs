package errors

import (
	"errors"
)

// ErrDuplicateQueue error for duplicate queue
var ErrDuplicateQueue = errors.New("duplicate queue")

// ErrUserNotFound error for unknown user
var ErrUserNotFound = errors.New("user not found")

// ErrQueueNotFound error for unknown queue
var ErrQueueNotFound = errors.New("queue not found")

// ErrMessageNotFound error for unknown message
var ErrMessageNotFound = errors.New("message not found")
