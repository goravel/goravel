package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("logging", map[string]interface{}{
		//Default Log Channel
		//This option defines the default log channel that gets used when writing
		//messages to the logs. The name specified in this option should match
		//one of the channels defined in the "channels" configuration array.
		"default": config.Env("LOG_CHANNEL", "stack"),

		//Log Channels
		//Here you may configure the log channels for your application.
		//Available Drivers: "single", "daily", "custom", "stack"
		//Available Level: "debug", "info", "warning", "error", "fatal", "panic"
		"channels": map[string]interface{}{
			"stack": map[string]interface{}{
				"driver":   "stack",
				"channels": []string{"single", "daily", "test"},
			},
			"single": map[string]interface{}{
				"driver": "single",
				"path":   "storage/logs/goravel.log",
				"level":  config.Env("LOG_LEVEL", "debug"),
			},
			"daily": map[string]interface{}{
				"driver": "daily",
				"path":   "storage/logs/goravel.log",
				"level":  config.Env("LOG_LEVEL", "debug"),
				"days":   7,
			},
			"test": map[string]interface{}{
				"driver": "custom",
				"path":   "storage/logs/test.log",
				"via":    &Logger{},
			},
		},
	})
}

type Logger struct {
}

// Handle 传入通道配置路径
func (logger *Logger) Handle(channel string) (log.Hook, error) {
	return &Hook{channel: channel}, nil
}

type Hook struct {
	channel string
}

// Levels 要监控的等级
func (h *Hook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarningLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

// Fire 当触发时执行的逻辑
func (h *Hook) Fire(entry *log.Entry) error {
	logPath := facades.Config.GetString(h.channel + ".path")
	err := os.MkdirAll(path.Dir(logPath), os.ModePerm)
	if err != nil {
		return errors.New("Create dir fail:" + err.Error())
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return errors.New("Failed to log to file:" + err.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(fmt.Sprintf("level=%v time=%v message=%s\n", entry.GetLevel(), entry.GetTime(), entry.GetMessage()))
	write.Flush()

	return nil
}
