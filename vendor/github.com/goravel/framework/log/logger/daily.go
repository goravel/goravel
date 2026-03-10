package logger

import (
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/goravel/file-rotatelogs/v2"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/carbon"
)

type Daily struct {
	config config.Config
	json   foundation.Json
}

func NewDaily(config config.Config, json foundation.Json) *Daily {
	return &Daily{
		config: config,
		json:   json,
	}
}

func (daily *Daily) Handle(channel string) (log.Handler, error) {
	logPath := daily.config.GetString(channel + ".path")
	if logPath == "" {
		return nil, errors.LogEmptyLogFilePath
	}

	ext := filepath.Ext(logPath)
	logPath = strings.ReplaceAll(logPath, ext, "")
	logPath = filepath.Join(support.RelativePath, logPath)

	writer, err := rotatelogs.New(
		logPath+"-%Y-%m-%d"+ext,
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationCount(uint(daily.config.GetInt(channel+".days"))),
		// When using carbon.SetTestNow(), carbon.Now().StdTime() should always be used to get the current time.
		// Hence, WithLocation cannot be used here.
		rotatelogs.WithClock(&rotatelogsClock{}),
	)
	if err != nil {
		return nil, err
	}

	level := GetLevelFromString(daily.config.GetString(channel + ".level"))
	formatter := daily.config.GetString(channel+".formatter", FormatterText)

	return NewRotatingFileHandler(writer, daily.config, daily.json, level, formatter), nil
}

func NewRotatingFileHandler(w io.Writer, config config.Config, json foundation.Json, level slog.Leveler, formatter string) log.Handler {
	return &IOHandler{
		writer:    w,
		config:    config,
		json:      json,
		level:     level,
		formatter: formatter,
	}
}

type rotatelogsClock struct{}

func (clock *rotatelogsClock) Now() time.Time {
	return carbon.Now().StdTime()
}
