package log

import (
	"context"
	"log/slog"

	"github.com/goravel/framework/contracts/log"
)

// slogAdapter wraps a log.Handler to implement slog.Handler
type slogAdapter struct {
	handler log.Handler
}

func (h *slogAdapter) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(log.Level(level))
}

func (h *slogAdapter) Handle(ctx context.Context, record slog.Record) error {
	entry := FromSlogRecord(record)
	return h.handler.Handle(entry)
}

func (h *slogAdapter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *slogAdapter) WithGroup(name string) slog.Handler {
	return h
}

func HandlerToSlogHandler(handler log.Handler) slog.Handler {
	return &slogAdapter{handler: handler}
}

// hookAdapter wraps a log.Hook to implement log.Handler for backward compatibility.
// Deprecated: Use Handler directly instead, hookAdapter will be removed in v1.18.
type hookAdapter struct {
	hook log.Hook
}

func (h *hookAdapter) Enabled(level log.Level) bool {
	for _, l := range h.hook.Levels() {
		if l == level {
			return true
		}
	}
	return false
}

func (h *hookAdapter) Handle(entry log.Entry) error {
	return h.hook.Fire(entry)
}

// HookToHandler converts a Hook to a Handler for backward compatibility.
// Deprecated: Use Handler directly instead, HookToHandler will be removed in v1.18.
func HookToHandler(hook log.Hook) log.Handler {
	return &hookAdapter{hook: hook}
}
