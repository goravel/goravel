package errorx

import (
	"fmt"
	"strconv"
	"strings"
)

// ErrorCoder interface
type ErrorCoder interface {
	error
	Code() int
}

// ErrorR useful for web service replay/response.
// code == 0 is successful. otherwise, is failed.
type ErrorR interface {
	ErrorCoder
	fmt.Stringer
	IsSuc() bool
	IsFail() bool
}

// error reply struct
type errorR struct {
	code int
	msg  string
}

// NewR code with error response
func NewR(code int, msg string) ErrorR {
	return &errorR{code: code, msg: msg}
}

// Fail code with error response
func Fail(code int, msg string) ErrorR {
	return &errorR{code: code, msg: msg}
}

// Failf code with error response
func Failf(code int, tpl string, v ...any) ErrorR {
	return &errorR{code: code, msg: fmt.Sprintf(tpl, v...)}
}

// Suc success response reply
func Suc(msg string) ErrorR {
	return &errorR{code: 0, msg: msg}
}

// IsSuc code value check
func (e *errorR) IsSuc() bool {
	return e.code == 0
}

// IsFail code value check
func (e *errorR) IsFail() bool {
	return e.code != 0
}

// Code value
func (e *errorR) Code() int {
	return e.code
}

// Error string
func (e *errorR) Error() string {
	return e.msg
}

// String get
func (e *errorR) String() string {
	return e.msg + "(code: " + strconv.FormatInt(int64(e.code), 10) + ")"
}

// GoString get.
func (e *errorR) GoString() string {
	return e.String()
}

// ErrorM multi error map
type ErrorM map[string]error

// ErrMap alias of ErrorM
type ErrMap = ErrorM

// Error string
func (e ErrorM) Error() string {
	var sb strings.Builder
	for name, err := range e {
		sb.WriteString(name)
		sb.WriteByte(':')
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ErrorOrNil error
func (e ErrorM) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// IsEmpty error
func (e ErrorM) IsEmpty() bool {
	return len(e) == 0
}

// One error
func (e ErrorM) One() error {
	for _, err := range e {
		return err
	}
	return nil
}

// Errors multi error list
type Errors []error

// ErrList alias for Errors
type ErrList = Errors

// Error string
func (es Errors) Error() string {
	var sb strings.Builder
	for _, err := range es {
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ErrorOrNil error
func (es Errors) ErrorOrNil() error {
	if len(es) == 0 {
		return nil
	}
	return es
}

// IsEmpty error
func (es Errors) IsEmpty() bool {
	return len(es) == 0
}

// First error
func (es Errors) First() error {
	if len(es) > 0 {
		return es[0]
	}
	return nil
}
