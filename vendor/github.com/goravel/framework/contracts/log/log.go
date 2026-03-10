package log

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/http"
)

const (
	DriverStack  = "stack"
	DriverSingle = "single"
	DriverDaily  = "daily"
	DriverOtel   = "otel"
	DriverCustom = "custom"
)

const (
	// Deprecated: use DriverStack instead, StackDriver will be removed in v1.18.
	StackDriver = DriverStack
	// Deprecated: use DriverSingle instead, SingleDriver will be removed in v1.18.
	SingleDriver = DriverSingle
	// Deprecated: use DriverDaily instead, DriverDaily will be removed in v1.18.
	DailyDriver = DriverDaily
	// Deprecated: use DriverCustom instead, CustomDriver will be removed in v1.18.
	CustomDriver = DriverCustom
)

type Data map[string]any

type Log interface {
	// WithContext adds a context to the logger.
	WithContext(ctx context.Context) Log
	// Channel return a writer for a specific channel.
	Channel(channel string) Log
	// Stack return a writer for multiple channels.
	Stack(channels []string) Log
	// Writer returns the underlying Writer.
	Writer
}

type Writer interface {
	// Debug logs a message at LevelDebug.
	Debug(args ...any)
	// Debugf is equivalent to Debug, but with support for fmt.Printf-style arguments.
	Debugf(format string, args ...any)
	// Info logs a message at LevelInfo.
	Info(args ...any)
	// Infof is equivalent to Info, but with support for fmt.Printf-style arguments.
	Infof(format string, args ...any)
	// Warning logs a message at LevelWarning.
	Warning(args ...any)
	// Warningf is equivalent to Warning, but with support for fmt.Printf-style arguments.
	Warningf(format string, args ...any)
	// Error logs a message at LevelError.
	Error(args ...any)
	// Errorf is equivalent to Error, but with support for fmt.Printf-style arguments.
	Errorf(format string, args ...any)
	// Fatal logs a message at LevelFatal.
	Fatal(args ...any)
	// Fatalf is equivalent to Fatal, but with support for fmt.Printf-style arguments.
	Fatalf(format string, args ...any)
	// Panic logs a message at LevelPanic.
	Panic(args ...any)
	// Panicf is equivalent to Panic, but with support for fmt.Printf-style arguments.
	Panicf(format string, args ...any)
	// Code set a code or slug that describes the error.
	// Error messages are intended to be read by humans, but such code is expected to
	// be read by machines and even transported over different services.
	Code(code string) Writer
	// Hint set a hint for faster debugging.
	Hint(hint string) Writer
	// In sets the feature category or domain in which the log entry is relevant.
	In(domain string) Writer
	// Owner set the name/email of the colleague/team responsible for handling this error.
	// Useful for alerting purpose.
	Owner(owner any) Writer
	// Request supplies a http.Request.
	Request(req http.ContextRequest) Writer
	// Response supplies a http.Response.
	Response(res http.ContextResponse) Writer
	// Tags add multiple tags, describing the feature returning an error.
	Tags(tags ...string) Writer
	// User sets the user associated with the log entry.
	User(user any) Writer
	// With adds key-value pairs to the context of the log entry
	With(data map[string]any) Writer
	// WithTrace adds a stack trace to the log entry.
	WithTrace() Writer
}

// Logger is the interface for custom log drivers.
type Logger interface {
	// Handle returns a Handler for the given channel config path.
	Handle(channel string) (Handler, error)
}

// Handler is the interface for log handlers.
type Handler interface {
	// Enabled reports whether the handler is enabled for the given level.
	Enabled(Level) bool
	// Handle handles the log entry.
	Handle(Entry) error
}

// Hook is the interface for log hooks.
// Deprecated: Use Handler instead, Hook will be removed in v1.18.
type Hook interface {
	// Levels monitoring level
	Levels() []Level
	// Fire executes logic when trigger
	Fire(Entry) error
}

type Entry interface {
	// Code returns the associated code.
	Code() string
	// Context returns the context of the entry.
	Context() context.Context
	// Data returns the data of the entry.
	Data() Data
	// Domain returns the domain of the entry.
	Domain() string
	// Hint returns the hint of the entry.
	Hint() string
	// Level returns the level of the entry.
	Level() Level
	// Message returns the message of the entry.
	Message() string
	// Owner returns the log's owner.
	Owner() any
	// Request returns the request data.
	Request() map[string]any
	// Response returns the response data.
	Response() map[string]any
	// Tags returns the list of tags.
	Tags() []string
	// Time returns the timestamp of the entry.
	Time() time.Time
	// Trace returns the stack trace or trace data.
	Trace() map[string]any
	// User returns the user information.
	User() any
	// With returns additional context data.
	With() map[string]any
}
