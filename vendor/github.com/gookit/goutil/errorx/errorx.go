// Package errorx provide an enhanced error implements for go,
// allow with stacktraces and wrap another error.
package errorx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Causer interface for get first cause error
type Causer interface {
	// Cause returns the first cause error by call err.Cause().
	// Otherwise, will returns current error.
	Cause() error
}

// Unwrapper interface for get previous error
type Unwrapper interface {
	// Unwrap returns previous error by call err.Unwrap().
	// Otherwise, will returns nil.
	Unwrap() error
}

// XErrorFace interface
type XErrorFace interface {
	error
	Causer
	Unwrapper
}

// Exception interface
// type Exception interface {
// 	XErrorFace
// 	Code() string
// 	Message() string
// 	StackString() string
// }

/*************************************************************
 * implements XErrorFace interface
 *************************************************************/

// ErrorX struct
//
// TIPS:
//
//	fmt pkg call order: Format > GoString > Error > String
type ErrorX struct {
	// trace stack
	*stack
	prev error
	msg  string
}

// Cause implements Causer.
func (e *ErrorX) Cause() error {
	if e.prev == nil {
		return e
	}

	if ex, ok := e.prev.(*ErrorX); ok {
		return ex.Cause()
	}
	return e.prev
}

// Unwrap implements Unwrapper.
func (e *ErrorX) Unwrap() error {
	return e.prev
}

// Format error, will output stack information.
func (e *ErrorX) Format(s fmt.State, verb rune) {
	// format current error: only output on have msg
	if len(e.msg) > 0 {
		_, _ = io.WriteString(s, e.msg)
		if e.stack != nil {
			e.stack.Format(s, verb)
		}
	}

	// format prev error
	if e.prev == nil {
		return
	}

	_, _ = s.Write([]byte("\nPrevious: "))
	if ex, ok := e.prev.(*ErrorX); ok {
		ex.Format(s, verb)
	} else {
		_, _ = s.Write([]byte(e.prev.Error()))
	}
}

// GoString to GO string, contains stack information.
// printing an error with %#v will produce useful information.
func (e *ErrorX) GoString() string {
	// var sb strings.Builder
	var buf bytes.Buffer
	_, _ = e.WriteTo(&buf)
	return buf.String()
}

// Error msg string, not contains stack information.
func (e *ErrorX) Error() string {
	var buf bytes.Buffer
	e.writeMsgTo(&buf)
	return buf.String()
}

// String error to string, contains stack information.
func (e *ErrorX) String() string {
	return e.GoString()
}

// WriteTo write the error to a writer, contains stack information.
func (e *ErrorX) WriteTo(w io.Writer) (n int64, err error) {
	// current error: only output on have msg
	if len(e.msg) > 0 {
		_, _ = w.Write([]byte(e.msg))

		// with stack
		if e.stack != nil {
			_, _ = e.stack.WriteTo(w)
		}
	}

	// with prev error
	if e.prev != nil {
		_, _ = io.WriteString(w, "\nPrevious: ")

		if ex, ok := e.prev.(*ErrorX); ok {
			_, _ = ex.WriteTo(w)
		} else {
			_, _ = io.WriteString(w, e.prev.Error())
		}
	}
	return
}

// Message error message of current
func (e *ErrorX) Message() string {
	return e.msg
}

// StackString returns error stack string of current.
func (e *ErrorX) StackString() string {
	if e.stack != nil {
		return e.stack.String()
	}
	return ""
}

// writeMsgTo write the error msg to a writer
func (e *ErrorX) writeMsgTo(w io.Writer) {
	// current error
	if len(e.msg) > 0 {
		_, _ = w.Write([]byte(e.msg))
	}

	// with prev error
	if e.prev != nil {
		_, _ = w.Write([]byte("; "))
		if ex, ok := e.prev.(*ErrorX); ok {
			ex.writeMsgTo(w)
		} else {
			_, _ = io.WriteString(w, e.prev.Error())
		}
	}
}

// CallerFunc returns the error caller func. if stack is nil, will return nil
func (e *ErrorX) CallerFunc() *Func {
	if e.stack == nil {
		return nil
	}
	return FuncForPC(e.stack.CallerPC())
}

// Location information for the caller func. more please see CallerFunc
//
// Returns eg:
//
//	github.com/gookit/goutil/errorx_test.TestWithPrev(), errorx_test.go:34
func (e *ErrorX) Location() string {
	if e.stack == nil {
		return "unknown"
	}
	return e.CallerFunc().Location()
}

/*************************************************************
 * new error with call stacks
 *************************************************************/

// New error message and with caller stacks
func New(msg string) error {
	return &ErrorX{
		msg:   msg,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// Newf error with format message, and with caller stacks.
// alias of Errorf()
func Newf(tpl string, vars ...any) error {
	return &ErrorX{
		msg:   fmt.Sprintf(tpl, vars...),
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// Errorf error with format message, and with caller stacks
func Errorf(tpl string, vars ...any) error {
	return &ErrorX{
		msg:   fmt.Sprintf(tpl, vars...),
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// With prev error and error message, and with caller stacks
func With(err error, msg string) error {
	return &ErrorX{
		msg:   msg,
		prev:  err,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// Withf error and with format message, and with caller stacks
func Withf(err error, tpl string, vars ...any) error {
	return &ErrorX{
		msg:   fmt.Sprintf(tpl, vars...),
		prev:  err,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// WithPrev error and message, and with caller stacks. alias of With()
func WithPrev(err error, msg string) error {
	return &ErrorX{
		msg:   msg,
		prev:  err,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// WithPrevf error and with format message, and with caller stacks. alias of Withf()
func WithPrevf(err error, tpl string, vars ...any) error {
	return &ErrorX{
		msg:   fmt.Sprintf(tpl, vars...),
		prev:  err,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

/*************************************************************
 * wrap go error with call stacks
 *************************************************************/

// WithStack wrap a go error with a stacked trace. If err is nil, will return nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	return &ErrorX{
		msg: err.Error(),
		// prev:  err,
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// Traced warp a go error and with caller stacks. alias of WithStack()
func Traced(err error) error {
	if err == nil {
		return nil
	}

	return &ErrorX{
		msg:   err.Error(),
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// Stacked warp a go error and with caller stacks. alias of WithStack()
func Stacked(err error) error {
	if err == nil {
		return nil
	}

	return &ErrorX{
		msg:   err.Error(),
		stack: callersStack(stdOpt.SkipDepth, stdOpt.TraceDepth),
	}
}

// WithOptions new error with some option func
func WithOptions(msg string, fns ...func(opt *ErrStackOpt)) error {
	opt := newErrOpt()
	for _, fn := range fns {
		fn(opt)
	}

	return &ErrorX{
		msg:   msg,
		stack: callersStack(opt.SkipDepth, opt.TraceDepth),
	}
}

/*************************************************************
 * helper func for wrap error without stacks
 *************************************************************/

// Wrap error and with message, but not with stack
func Wrap(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}

	return &ErrorX{
		msg:  msg,
		prev: err,
	}
}

// Wrapf error with format message, but not with stack
func Wrapf(err error, tpl string, vars ...any) error {
	if err == nil {
		return fmt.Errorf(tpl, vars...)
	}

	return &ErrorX{
		msg:  fmt.Sprintf(tpl, vars...),
		prev: err,
	}
}
