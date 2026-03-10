package log

import (
	"fmt"
	"log/slog"
	"strings"
)

// Level defines custom log levels for the logging system.
// We define custom levels that extend slog's built-in levels to support
// Panic and Fatal levels which are not part of the standard slog package.
type Level slog.Level

const (
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug = Level(slog.LevelDebug) // -4
	// LevelInfo level. General operational entries about what's going on inside the application.
	LevelInfo = Level(slog.LevelInfo) // 0
	// LevelWarning level. Non-critical entries that deserve eyes.
	LevelWarning = Level(slog.LevelWarn) // 4
	// LevelError level. Used for errors that should definitely be noted.
	LevelError = Level(slog.LevelError) // 8
	// LevelFatal level. Logs and then calls `os.Exit(1)`.
	LevelFatal = Level(slog.LevelError + 4) // 12
	// LevelPanic level. Highest level of severity. Logs and then calls panic.
	LevelPanic = Level(slog.LevelError + 8) // 16
)

const (
	// DebugLevel is an alias for LevelDebug level.
	// Deprecated: use LevelDebug instead, DebugLevel will be removed in v1.18.
	DebugLevel = LevelDebug
	// InfoLevel is an alias for LevelInfo level.
	// Deprecated: use LevelInfo instead, InfoLevel will be removed in v1.18.
	InfoLevel = LevelInfo
	// WarningLevel is an alias for LevelWarning level.
	// Deprecated: use LevelWarning instead, WarningLevel will be removed in v1.18.
	WarningLevel = LevelWarning
	// ErrorLevel is an alias for LevelError level.
	// Deprecated: use LevelError instead, ErrorLevel will be removed in v1.18.
	ErrorLevel = LevelError
	// FatalLevel is an alias for LevelFatal level.
	// Deprecated: use LevelFatal instead, FatalLevel will be removed in v1.18.
	FatalLevel = LevelFatal
	// PanicLevel is an alias for LevelPanic level.
	// Deprecated: use LevelPanic instead, PanicLevel will be removed in v1.18.
	PanicLevel = LevelPanic
)

// String converts the Level to a string. E.g. LevelPanic becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}
	return "unknown"
}

// MarshalText implements encoding.TextMarshaler.
func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case LevelDebug:
		return []byte("debug"), nil
	case LevelInfo:
		return []byte("info"), nil
	case LevelWarning:
		return []byte("warning"), nil
	case LevelError:
		return []byte("error"), nil
	case LevelFatal:
		return []byte("fatal"), nil
	case LevelPanic:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid log level %d", level)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = l

	return nil
}

// Level implements the slog.Leveler interface.
func (level Level) Level() slog.Level {
	return slog.Level(level)
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return LevelPanic, nil
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarning, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid log Level: %q", lvl)
}
