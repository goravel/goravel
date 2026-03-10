package log

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/goravel/framework/contracts/log"
)

var entryPool = sync.Pool{
	New: func() any {
		return &Entry{
			with: make(map[string]any),
		}
	},
}

func acquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

func releaseEntry(e *Entry) {
	e.time = time.Time{}
	e.ctx = nil
	e.owner = nil
	e.user = nil
	e.data = nil
	e.request = nil
	e.response = nil
	e.stacktrace = nil
	e.code = ""
	e.domain = ""
	e.hint = ""
	e.message = ""
	e.tags = nil
	e.level = 0

	clear(e.with)

	entryPool.Put(e)
}

type Entry struct {
	time       time.Time
	ctx        context.Context
	owner      any
	user       any
	data       log.Data
	request    map[string]any
	response   map[string]any
	stacktrace map[string]any
	with       map[string]any
	code       string
	domain     string
	hint       string
	message    string
	tags       []string
	level      log.Level
}

func (e *Entry) Code() string {
	return e.code
}

func (e *Entry) Context() context.Context {
	return e.ctx
}

func (e *Entry) Data() log.Data {
	return e.data
}

func (e *Entry) Domain() string {
	return e.domain
}

func (e *Entry) Hint() string {
	return e.hint
}

func (e *Entry) Level() log.Level {
	return e.level
}

func (e *Entry) Message() string {
	return e.message
}

func (e *Entry) Owner() any {
	return e.owner
}

func (e *Entry) Request() map[string]any {
	return e.request
}

func (e *Entry) Response() map[string]any {
	return e.response
}

func (e *Entry) Tags() []string {
	return e.tags
}

func (e *Entry) Time() time.Time {
	return e.time
}

func (e *Entry) Trace() map[string]any {
	return e.stacktrace
}

func (e *Entry) User() any {
	return e.user
}

func (e *Entry) With() map[string]any {
	return e.with
}

func (e *Entry) ToSlogRecord() slog.Record {
	r := slog.NewRecord(
		e.time,
		slog.Level(e.level),
		e.message,
		0,
	)

	if e.code != "" {
		r.Add("code", e.code)
	}
	if e.ctx != nil {
		r.Add("context", e.ctx)
	}
	if e.domain != "" {
		r.Add("domain", e.domain)
	}
	if e.hint != "" {
		r.Add("hint", e.hint)
	}
	if e.owner != nil {
		r.Add("owner", e.owner)
	}
	if e.request != nil {
		r.Add("request", e.request)
	}
	if e.response != nil {
		r.Add("response", e.response)
	}
	if e.stacktrace != nil {
		r.Add("stacktrace", e.stacktrace)
	}
	if len(e.tags) > 0 {
		r.Add("tags", e.tags)
	}
	if e.user != nil {
		r.Add("user", e.user)
	}
	if len(e.with) > 0 {
		r.Add("with", e.with)
	}

	return r
}

func FromSlogRecord(r slog.Record) *Entry {
	e := acquireEntry()
	e.time = r.Time
	e.message = r.Message
	e.level = log.Level(r.Level)

	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "code":
			e.code = a.Value.String()
			return true
		case "context":
			if a.Value.Kind() == slog.KindAny {
				if ctx, ok := a.Value.Any().(context.Context); ok {
					e.ctx = ctx
				}
			}
			return true
		case "domain":
			e.domain = a.Value.String()
			return true
		case "hint":
			e.hint = a.Value.String()
			return true
		case "message":
			e.message = a.Value.String()
			return true
		case "owner":
			e.owner = a.Value.Any()
			return true
		case "request":
			if a.Value.Kind() == slog.KindAny {
				if request, ok := a.Value.Any().(map[string]any); ok {
					e.request = request
				}
			}
			return true
		case "response":
			if a.Value.Kind() == slog.KindAny {
				if response, ok := a.Value.Any().(map[string]any); ok {
					e.response = response
				}
			}
			return true
		case "stacktrace":
			if a.Value.Kind() == slog.KindAny {
				if trace, ok := a.Value.Any().(map[string]any); ok {
					e.stacktrace = trace
				}
			}
			return true
		case "tags":
			if a.Value.Kind() == slog.KindAny {
				if tags, ok := a.Value.Any().([]string); ok {
					e.tags = tags
				}
			}
			return true
		case "user":
			e.user = a.Value.Any()
			return true
		case "with":
			if with, ok := a.Value.Any().(map[string]any); ok {
				e.with = with
			}
			return true
		default:
			e.with[a.Key] = a.Value.Any()
			return true
		}
	})

	return e
}
