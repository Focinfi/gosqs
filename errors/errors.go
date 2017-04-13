package errors

import (
	"errors"
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
	duplicateQueue              = 1001
	duplicateMessage            = 1002
	messageOutOfData            = 1003
	userNotFound                = 1004
	queueNotFound               = 1005
	messageNotFound             = 1006
	clientNotFound              = 1007
	duplicateSquad              = 1008
	clientHasAlreadyRegistered  = 1009
	serviceOverload             = 1010
	applyMessageIDRangeOversize = 1011
	messageIndexOutOfRange      = 1012
	dataNotFound                = 1013
)

// DuplicateQueue error for duplicate queue
var DuplicateQueue = NewBizErr("duplicate queue", duplicateMessage)

// DuplicateMessage error for duplicate message
var DuplicateMessage = NewBizErr("duplicate message", duplicateMessage)

// DuplicateSquad error for duplicate squad
var DuplicateSquad = NewBizErr("duplicate squad", duplicateSquad)

// MessageOutOfDate error for out-of-date message
var MessageOutOfDate = NewBizErr("message is out of date", messageOutOfData)

// MessageIndexOutOfRange error for out-of-range message index
var MessageIndexOutOfRange = NewBizErr("message index is out of range", messageIndexOutOfRange)

// UserNotFound error for unknown user
var UserNotFound = NewBizErr("user not found", userNotFound)

// QueueNotFound error for unknown queue
var QueueNotFound = NewBizErr("queue not found", queueNotFound)

// MessageNotFound error for unknown message
var MessageNotFound = NewBizErr("message not found", messageNotFound)

// ClientNotFound error for unknown client
var ClientNotFound = NewBizErr("client not found", clientNotFound)

// ApplyMessageIDRangeOversize error for oversize message id range application
var ApplyMessageIDRangeOversize = NewBizErr("apply message id range oversize", applyMessageIDRangeOversize)

// NoConsumer error for no consumer
var NoConsumer = New("cosumer queue is empty")

// ClientHasAlreadyRegistered error for client has already registered
var ClientHasAlreadyRegistered = NewBizErr("client has already registered", clientHasAlreadyRegistered)

// ServiceOverload error for node is overload
var ServiceOverload = NewBizErr("node is overload", serviceOverload)

// DataNotFound error for data not fouond
var DataNotFound = NewBizErr("data is not found", dataNotFound)

// DBQueryTimeout returns a Internal for a db query
func DBQueryTimeout(db, key string) Internal {
	return NewInternalErrorf(fmt.Sprintf("db: %s, key: %s, query timeout", db, key))
}

// DataLost returns a internal error for losting data
func DataLost(key string) error {
	return NewInternalErrorf(fmt.Sprintf("data lost: key= %s", key))
}

// DataBroken returns a internal error for broken data
func DataBroken(key string, err error) error {
	return NewInternalErrorf(fmt.Sprintf("data broken: key=%s, err: %s", key, err))
}

// FailedEncoding returns a internal error for encoding error
func FailedEncoding(data interface{}) error {
	return NewInternalErrorf(fmt.Sprintf("failed encoding for data: %v", data))
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

// NewInternalWrap returns a new internalErr
func NewInternalWrap(err error) Internal {
	return &internalErr{
		message: err.Error(),
	}
}

// NewInternalErrorf returns a new internalErr
func NewInternalErrorf(format string, a ...interface{}) Internal {
	return &internalErr{
		message: fmt.Sprintf(format, a...),
	}
}

func (err internalErr) Error() string {
	return err.message
}

func (err internalErr) isInternal() bool {
	return true
}

// New returns a new error with the message
func New(message string) error {
	return errors.New(message)
}
