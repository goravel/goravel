package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
)

// Formatter types for log output
const (
	FormatterText = "text"
	FormatterJson = "json"
)

type IOHandler struct {
	writer    io.Writer
	config    config.Config
	json      foundation.Json
	level     slog.Leveler
	formatter string
}

func NewIOHandler(w io.Writer, config config.Config, json foundation.Json, level slog.Leveler, formatter string) *IOHandler {
	return &IOHandler{
		writer:    w,
		config:    config,
		json:      json,
		level:     level,
		formatter: formatter,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *IOHandler) Enabled(level log.Level) bool {
	return level.Level() >= h.level.Level()
}

// Handle handles the Record.
func (h *IOHandler) Handle(entry log.Entry) error {
	switch h.formatter {
	case FormatterJson:
		return h.handleJSON(entry)
	case FormatterText:
		return h.handleText(entry)
	default:
		return errors.LogFormatterNotSupported.Args(h.formatter)
	}
}

// handleText formats the log entry as human-readable text.
func (h *IOHandler) handleText(entry log.Entry) error {
	var b bytes.Buffer

	timestamp := carbon.FromStdTime(entry.Time(), carbon.DefaultTimezone()).ToDateTimeMilliString()
	env := h.config.GetString("app.env")

	_, err := fmt.Fprintf(&b, "[%s] %s.%s: %s\n", timestamp, env, entry.Level().String(), entry.Message())
	if err != nil {
		return err
	}

	// Format Entry
	if v := entry.Code(); v != "" {
		_, _ = fmt.Fprintf(&b, "[Code] %+v\n", v)
	}
	if v := entry.Context(); v != nil {
		values := make(map[any]any)
		getContextValues(v, values)
		if len(values) > 0 {
			_, _ = fmt.Fprintf(&b, "[Context] %+v\n", values)
		}
	}
	if v := entry.Domain(); v != "" {
		_, _ = fmt.Fprintf(&b, "[Domain] %+v\n", v)
	}
	if v := entry.Hint(); v != "" {
		_, _ = fmt.Fprintf(&b, "[Hint] %+v\n", v)
	}
	if v := entry.Owner(); v != nil {
		_, _ = fmt.Fprintf(&b, "[Owner] %+v\n", v)
	}
	if v := entry.Request(); v != nil {
		_, _ = fmt.Fprintf(&b, "[Request] %+v\n", v)
	}
	if v := entry.Response(); v != nil {
		_, _ = fmt.Fprintf(&b, "[Response] %+v\n", v)
	}
	if v := entry.Trace(); v != nil {
		traces, err := formatStackTraces(h.json, v)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(&b, "[Trace] %+v", traces)
	}
	if v := entry.Tags(); len(v) > 0 {
		_, _ = fmt.Fprintf(&b, "[Tags] %+v\n", v)
	}
	if v := entry.User(); v != nil {
		_, _ = fmt.Fprintf(&b, "[User] %+v\n", v)
	}
	if v := entry.With(); len(v) > 0 {
		_, _ = fmt.Fprintf(&b, "[With] %+v\n", v)
	}

	_, err = h.writer.Write(b.Bytes())

	return err
}

// handleJSON formats the log entry as a JSON object (one line per entry).
func (h *IOHandler) handleJSON(entry log.Entry) error {
	timestamp := carbon.FromStdTime(entry.Time(), carbon.DefaultTimezone()).ToDateTimeMilliString()
	env := h.config.GetString("app.env")

	data := map[string]any{
		"time":        timestamp,
		"environment": env,
		"level":       entry.Level().String(),
		"message":     entry.Message(),
	}

	if v := entry.Code(); v != "" {
		data["code"] = v
	}
	if v := entry.Context(); v != nil {
		values := make(map[any]any)
		getContextValues(v, values)
		if len(values) > 0 {
			// Convert map[any]any to map[string]any for JSON serialization
			stringValues := make(map[string]any)
			for k, val := range values {
				stringValues[fmt.Sprintf("%v", k)] = val
			}
			data["context"] = stringValues
		}
	}
	if v := entry.Domain(); v != "" {
		data["domain"] = v
	}
	if v := entry.Hint(); v != "" {
		data["hint"] = v
	}
	if v := entry.Owner(); v != nil {
		data["owner"] = v
	}
	if v := entry.Request(); v != nil {
		data["request"] = v
	}
	if v := entry.Response(); v != nil {
		data["response"] = v
	}
	if v := entry.Trace(); v != nil {
		data["trace"] = v
	}
	if v := entry.Tags(); len(v) > 0 {
		data["tags"] = v
	}
	if v := entry.User(); v != nil {
		data["user"] = v
	}
	if v := entry.With(); len(v) > 0 {
		data["extra"] = v
	}

	jsonBytes, err := h.json.Marshal(data)
	if err != nil {
		return err
	}

	// Append newline for line-delimited JSON
	jsonBytes = append(jsonBytes, '\n')

	_, err = h.writer.Write(jsonBytes)
	return err
}

type ConsoleHandler struct {
	*IOHandler
}

func NewConsoleHandler(config config.Config, json foundation.Json, level slog.Leveler, formatter string) *ConsoleHandler {
	return &ConsoleHandler{
		IOHandler: &IOHandler{
			writer:    os.Stdout,
			config:    config,
			json:      json,
			level:     level,
			formatter: formatter,
		},
	}
}

type StackTrace struct {
	Root struct {
		Message string   `json:"message"`
		Stack   []string `json:"stack"`
	} `json:"root"`
	Wrap []struct {
		Message string `json:"message"`
		Stack   string `json:"stack"`
	} `json:"wrap"`
}

func formatStackTraces(json foundation.Json, stackTraces any) (string, error) {
	var formattedTraces strings.Builder
	data, err := json.Marshal(stackTraces)

	if err != nil {
		return "", err
	}
	var traces StackTrace
	err = json.Unmarshal(data, &traces)
	if err != nil {
		return "", err
	}
	root := traces.Root
	if len(root.Stack) > 0 {
		for _, stackStr := range root.Stack {
			formattedTraces.WriteString(formatStackTrace(stackStr))
		}
	}

	return formattedTraces.String(), nil
}

func formatStackTrace(stackStr string) string {
	lastColon := strings.LastIndex(stackStr, ":")
	if lastColon > 0 && lastColon < len(stackStr)-1 {
		secondLastColon := strings.LastIndex(stackStr[:lastColon], ":")
		if secondLastColon > 0 {
			fileLine := stackStr[secondLastColon+1:]
			method := stackStr[:secondLastColon]
			return fmt.Sprintf("%s [%s]\n", fileLine, method)
		}
	}
	return fmt.Sprintf("%s\n", stackStr)
}
