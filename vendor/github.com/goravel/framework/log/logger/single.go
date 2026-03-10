package logger

import (
	"os"
	"path/filepath"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
)

type Single struct {
	config config.Config
	json   foundation.Json
}

func NewSingle(config config.Config, json foundation.Json) *Single {
	return &Single{
		config: config,
		json:   json,
	}
}

func (single *Single) Handle(channel string) (log.Handler, error) {
	logPath := single.config.GetString(channel + ".path")
	if logPath == "" {
		return nil, errors.LogEmptyLogFilePath
	}

	logPath = filepath.Join(support.RelativePath, logPath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	level := GetLevelFromString(single.config.GetString(channel + ".level"))
	formatter := single.config.GetString(channel+".formatter", FormatterText)

	return NewIOHandler(file, single.config, single.json, level, formatter), nil
}

// GetLevelFromString converts a string log level to log.Level.
func GetLevelFromString(level string) log.Level {
	l, err := log.ParseLevel(level)
	if err != nil {
		return log.LevelDebug
	}
	return l
}
