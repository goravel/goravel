package schedule

import (
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/support/color"
)

type Logger struct {
	log   log.Log
	debug bool
}

func NewLogger(log log.Log, debug bool) *Logger {
	return &Logger{
		debug: debug,
		log:   log,
	}
}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	if !log.debug {
		return
	}
	color.Successf("%s %v\n", msg, keysAndValues)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {
	log.log.Error(msg, keysAndValues)
}
