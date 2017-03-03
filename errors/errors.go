package errors

import (
	"fmt"
)

const (
	// ParamFormatErr for wrong param format
	ParamFormatErr = 998
	// InternalErr for internal error code
	InternalErr = 999
	// NoErr for successful request
	NoErr = 1000

	// biz error code
	duplicateQueue             = 1001
	duplicateMessage           = 1002
	messageOutOfData           = 1003
	userNotFound               = 1004
	queueNotFound              = 1005
	messageNotFound            = 1006
	clientNotFound             = 1007
	duplicateClient            = 1008
	clientHasAlreadyRegistered = 1009
	serviceOverload            = 1010
)

// DuplicateQueue error for duplicate queue
var DuplicateQueue = NewBizErr("duplicate queue", duplicateMessage)

// DuplicateMessage error for duplicate message
var DuplicateMessage = NewBizErr("duplicate message", duplicateMessage)

// DuplicateClient error for duplicate lient
var DuplicateClient = NewBizErr("duplicate lient", duplicateClient)

// MessageOutOfData erros for out-of-date message
var MessageOutOfData = NewBizErr("message is out of date", messageOutOfData)

// UserNotFound error for unknown user
var UserNotFound = NewBizErr("user not found", userNotFound)

// QueueNotFound error for unknown queue
var QueueNotFound = NewBizErr("queue not found", queueNotFound)

// MessageNotFound error for unknown message
var MessageNotFound = NewBizErr("message not found", messageNotFound)

// ClientNotFound error for unknown client
var ClientNotFound = NewBizErr("client not found", clientNotFound)

// ClientHasAlreadyRegistered error for client has already registered
var ClientHasAlreadyRegistered = NewBizErr("client has already registered", clientHasAlreadyRegistered)

// ServiceOverload error for service is overload
var ServiceOverload = NewBizErr("service is overload", serviceOverload)

// DataLost returns a internal error for losting data
func DataLost(key string) error {
	return NewInternalErr(fmt.Sprintf("data lost: key= %s", key))
}

// DataBroken returns a internal error for broken data
func DataBroken(key string, err error) error {
	return NewInternalErr(fmt.Sprintf("data broken: key=%s, err: %s", key, err))
}

// FailedEncoding returns a internal error for encoding error
func FailedEncoding(data interface{}) error {
	return NewInternalErr(fmt.Sprintf("failed encoding for data: %v", data))
}

// Biz detects the biz errors
type Biz interface {
	error
	BizCode() int
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
	code    int
	message string
}

// NewBizErr returns a new bizErr
func NewBizErr(message string, code int) Biz {
	return &bizErr{
		code:    code,
		message: message,
	}
}

func (err bizErr) Error() string {
	return err.message
}

func (err bizErr) BizCode() int {
	return err.code
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
