package log

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"

	"github.com/dromara/carbon/v2"
	"github.com/rotisserie/eris"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
)

type Writer struct {
	logger *slog.Logger
	ctx    context.Context
	entry  *Entry // nil for base writer, only set when fluent methods are called
}

func NewWriter(ctx context.Context, logger *slog.Logger) log.Writer {
	return &Writer{
		logger: logger,
		ctx:    ctx,
		entry:  nil,
	}
}

func (r *Writer) Debug(args ...any) {
	r.log(log.LevelDebug, fmt.Sprint(args...))
}

func (r *Writer) Debugf(format string, args ...any) {
	r.log(log.LevelDebug, fmt.Sprintf(format, args...))
}

func (r *Writer) Info(args ...any) {
	r.log(log.LevelInfo, fmt.Sprint(args...))
}

func (r *Writer) Infof(format string, args ...any) {
	r.log(log.LevelInfo, fmt.Sprintf(format, args...))
}

func (r *Writer) Warning(args ...any) {
	r.log(log.LevelWarning, fmt.Sprint(args...))
}

func (r *Writer) Warningf(format string, args ...any) {
	r.log(log.LevelWarning, fmt.Sprintf(format, args...))
}

func (r *Writer) Error(args ...any) {
	msg := fmt.Sprint(args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelError, msg)
}

func (r *Writer) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelError, msg)
}

func (r *Writer) Fatal(args ...any) {
	msg := fmt.Sprint(args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelFatal, msg)
	os.Exit(1)
}

func (r *Writer) Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelFatal, msg)
	os.Exit(1)
}

func (r *Writer) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelPanic, msg)
	panic(msg)
}

func (r *Writer) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	nw := r.ensureEntry()
	nw.withStackTrace(msg)
	nw.log(log.LevelPanic, msg)
	panic(msg)
}

// Code set a code or slug that describes the error.
func (r *Writer) Code(code string) log.Writer {
	nw := r.clone()
	nw.entry.code = code
	return nw
}

// Hint set a hint for faster debugging.
func (r *Writer) Hint(hint string) log.Writer {
	nw := r.clone()
	nw.entry.hint = hint
	return nw
}

// In sets the feature category or domain in which the log entry is relevant.
func (r *Writer) In(domain string) log.Writer {
	nw := r.clone()
	nw.entry.domain = domain
	return nw
}

// Owner set the name/email of the colleague/team responsible for handling this error.
func (r *Writer) Owner(owner any) log.Writer {
	nw := r.clone()
	nw.entry.owner = owner
	return nw
}

// Request supplies a http.Request.
func (r *Writer) Request(req http.ContextRequest) log.Writer {
	nw := r.clone()
	if req != nil {
		nw.entry.request = map[string]any{
			"method": req.Method(),
			"uri":    req.FullUrl(),
			"header": req.Headers(),
			"body":   req.All(),
		}
	}
	return nw
}

// Response supplies a http.Response.
func (r *Writer) Response(res http.ContextResponse) log.Writer {
	nw := r.clone()
	if res != nil {
		nw.entry.response = map[string]any{
			"status": res.Origin().Status(),
			"header": res.Origin().Header(),
			"body":   res.Origin().Body(),
			"size":   res.Origin().Size(),
		}
	}
	return nw
}

// Tags add multiple tags, describing the feature returning an error.
func (r *Writer) Tags(tags ...string) log.Writer {
	nw := r.clone()
	nw.entry.tags = append(nw.entry.tags, tags...)
	return nw
}

// User sets the user associated with the log entry.
func (r *Writer) User(user any) log.Writer {
	nw := r.clone()
	nw.entry.user = user
	return nw
}

// With adds key-value pairs to the context of the log entry.
func (r *Writer) With(data map[string]any) log.Writer {
	nw := r.clone()
	maps.Copy(nw.entry.with, data)
	return nw
}

// WithTrace adds a stack trace to the log entry.
func (r *Writer) WithTrace() log.Writer {
	nw := r.clone()
	nw.withStackTrace("")
	return nw
}

func (r *Writer) log(level log.Level, msg string) {
	entry := r.entry
	if entry == nil {
		// For direct log calls without fluent methods, acquire a fresh entry
		entry = acquireEntry()
		entry.ctx = r.ctx
	}

	entry.time = carbon.Now().StdTime()
	entry.message = msg
	entry.level = level

	_ = r.logger.Handler().Handle(entry.ctx, entry.ToSlogRecord())
	releaseEntry(entry)
}

func (r *Writer) withStackTrace(message string) {
	erisNew := eris.New(message)
	if erisNew == nil {
		return
	}

	format := eris.NewDefaultJSONFormat(eris.FormatOptions{
		InvertOutput: true,
		WithTrace:    true,
		InvertTrace:  true,
	})
	r.entry.stacktrace = eris.ToCustomJSON(erisNew, format)
}

func (r *Writer) getEntry() *Entry {
	if r.entry == nil {
		entry := acquireEntry()
		entry.ctx = r.ctx
		return entry
	}
	return r.entry
}

func (r *Writer) clone() *Writer {
	entry := r.getEntry()
	return &Writer{
		logger: r.logger,
		ctx:    r.ctx,
		entry:  entry,
	}
}

func (r *Writer) ensureEntry() *Writer {
	if r.entry != nil {
		return r
	}
	return r.clone()
}
