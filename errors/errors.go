package errors

import (
	"fmt"
)

// DuplicateQueue error for duplicate queue
var DuplicateQueue = NewBizErr("duplicate queue")

// DuplicateMessage error for duplicate message
var DuplicateMessage = NewBizErr("duplicate message")

// ErrUserNotFound error for unknown user
var ErrUserNotFound = NewBizErr("user not found")

// QueueNotFound error for unknown queue
var QueueNotFound = NewBizErr("queue not found")

// MessageNotFound error for unknown message
var MessageNotFound = NewBizErr("message not found")

// ClientNotFound error for unknown client
var ClientNotFound = NewBizErr("client not found")

// DataLost returns a internal error for losting data
func DataLost(key string) error {
	return NewInternalErr(fmt.Sprintf("data lost: key= %s", key))
}

// DataBroken returns a internal error for broken data
func DataBroken(key string) error {
	return NewInternalErr(fmt.Sprintf("data broken: key= %s", key))
}

// FailedEncoding returns a internal error for encoding error
func FailedEncoding(data interface{}) error {
	return NewInternalErr(fmt.Sprintf("failed encoding for data: %v", data))
}

// Biz detects the biz errors
type Biz interface {
	error
	isBiz() bool
}

// IsBiz for detecting if err is Biz
func IsBiz(err error) bool {
	_, ok := err.(Biz)
	return ok
}

// Internal detects the internal errors
type Internal interface {
	error
	isInternal() bool
}

// IsInternal for detecting if err is Internal
func IsInternal(err error) bool {
	_, ok := err.(Internal)
	return ok
}

type bizErr struct {
	message string
}

// NewBizErr returns a new bizErr
func NewBizErr(message string) Biz {
	return &bizErr{
		message: message,
	}
}

func (err bizErr) Error() string {
	return err.message
}

func (err bizErr) isBiz() bool {
	return true
}

type internalErr struct {
	message string
}

// NewInternalErr returns a new internalErr
func NewInternalErr(message string) Internal {
	return &internalErr{
		message: message,
	}
}

func (err internalErr) Error() string {
	return err.message
}

func (err internalErr) isInternal() bool {
	return true
}
