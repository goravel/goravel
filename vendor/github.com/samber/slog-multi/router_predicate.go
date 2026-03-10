package slogmulti

import (
	"context"
	"log/slog"
	"strings"
)

// LevelIs returns a function that checks if the record level is in the given levels.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.LevelIs(slog.LevelInfo)).
//	    Add(fileHandler, slogmulti.LevelIs(slog.LevelError)).
//	    Handler()
//
// Args:
//
//	levels: The levels to match
//
// Returns:
//
//	A function that checks if the record level is in the given levels
func LevelIs(levels ...slog.Level) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		for _, level := range levels {
			if r.Level == level {
				return true
			}
		}
		return false
	}
}

// LevelIsNot returns a function that checks if the record level is not in the given levels.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.LevelIsNot(slog.LevelInfo)).
//	    Add(fileHandler, slogmulti.LevelIsNot(slog.LevelError)).
//	    Handler()
//
// Args:
//
//	levels: The levels to check
//
// Returns:
//
//	A function that checks if the record level is not in the given levels
func LevelIsNot(levels ...slog.Level) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		for _, level := range levels {
			if r.Level == level {
				return false
			}
		}
		return true
	}
}

// MessageIs returns a function that checks if the record message is equal to the given message.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.MessageIs("database error")).
//	    Add(fileHandler, slogmulti.MessageIs("database error")).
//	    Handler()
//
// Args:
//
//	msg: The message to check
//
// Returns:
//
//	A function that checks if the record message is equal to the given message
func MessageIs(msg string) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		return r.Message == msg
	}
}

// MessageIsNot returns a function that checks if the record message is not equal to the given message.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.MessageIsNot("database error")).
//	    Add(fileHandler, slogmulti.MessageIsNot("database error")).
//	    Handler()
//
// Args:
//
//	msg: The message to check
//
// Returns:
//
//	A function that checks if the record message is not equal to the given message
func MessageIsNot(msg string) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		return r.Message != msg
	}
}

// MessageContains returns a function that checks if the record message contains the given part.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.MessageContains("database error")).
//	    Add(fileHandler, slogmulti.MessageContains("database error")).
//	    Handler()
//
// Args:
//
//	part: The part to check
//
// Returns:
//
//	A function that checks if the record message contains the given part
func MessageContains(part string) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		return strings.Contains(r.Message, part)
	}
}

// MessageNotContains returns a function that checks if the record message does not contain the given part.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.MessageNotContains("database error")).
//	    Add(fileHandler, slogmulti.MessageNotContains("database error")).
//	    Handler()
//
// Args:
//
//	part: The part to check
//
// Returns:
//
//	A function that checks if the record message does not contain the given part
func MessageNotContains(part string) func(ctx context.Context, r slog.Record) bool {
	return func(ctx context.Context, r slog.Record) bool {
		return !strings.Contains(r.Message, part)
	}
}

// AttrValueIs returns a function that checks if the record has all specified attributes with exact values.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.AttrValueIs("scope", "influx")).
//	    Add(fileHandler, slogmulti.AttrValueIs("env", "production", "region", "us-east")).
//	    Handler()
//
// Args:
//
//	args: Pairs of attribute key (string) and expected value (any)
//
// Returns:
//
//	A function that checks if the record has all specified attributes with exact values
func AttrValueIs(args ...any) func(ctx context.Context, r slog.Record) bool {
	if len(args)%2 != 0 {
		panic("AttrValueIs requires key/value pairs")
	}
	m := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		key, ok1 := args[i].(string)
		value := args[i+1]
		if !ok1 {
			panic("AttrValueIs requires string keys")
		}
		m[key] = value
	}

	return func(ctx context.Context, r slog.Record) bool {
		count := 0
		r.Attrs(func(attr slog.Attr) bool {
			if v, ok := m[attr.Key]; ok && attr.Value.Any() == v {
				count++
				if count == len(m) {
					return false // early exit
				}
			}
			return true
		})
		return count == len(m)
	}
}

// AttrKindIs returns a function that checks if the record has an attribute with the given key and type.
// Example usage:
//
//	r := slogmulti.Router().
//	    Add(consoleHandler, slogmulti.AttrKindIs("user_id", slog.KindString)).
//	    Add(fileHandler, slogmulti.AttrKindIs("user_id", slog.KindString)).
//	    Handler()
//
// Args:
//
//	key: The attribute key to check
//	ty: The attribute type to check
//
// Returns:
//
//	A function that checks if the record has an attribute with the given key and type
func AttrKindIs(args ...any) func(ctx context.Context, r slog.Record) bool {
	if len(args)%2 != 0 {
		panic("AttrKindIs requires key/kind pairs")
	}
	m := map[string]slog.Kind{}
	for i := 0; i < len(args); i += 2 {
		key, ok1 := args[i].(string)
		ty, ok2 := args[i+1].(slog.Kind)
		if !ok1 || !ok2 {
			panic("AttrKindIs requires string keys and slog.Kind values")
		}
		m[key] = ty
	}

	return func(ctx context.Context, r slog.Record) bool {
		count := 0
		r.Attrs(func(attr slog.Attr) bool {
			if ty, ok := m[attr.Key]; ok && attr.Value.Kind() == ty {
				count++
				if count == len(m) {
					return false // early exit
				}
			}
			return true
		})
		return count == len(m)
	}
}
