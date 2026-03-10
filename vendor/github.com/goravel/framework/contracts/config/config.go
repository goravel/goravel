package config

import "time"

type Config interface {
	// Env get config from env.
	Env(envName string, defaultValue ...any) any
	// EnvString get string value from env with optional default.
	EnvString(envName string, defaultValue ...string) string
	// EnvBool get bool value from env with optional default.
	EnvBool(envName string, defaultValue ...bool) bool
	// Add config to application.
	Add(name string, configuration any)
	// Get config from application.
	Get(path string, defaultValue ...any) any
	// GetString get string type config from application.
	GetString(path string, defaultValue ...string) string
	// GetInt get int type config from application.
	GetInt(path string, defaultValue ...int) int
	// GetBool get bool type config from application.
	GetBool(path string, defaultValue ...bool) bool
	// GetDuration get duration type config from application
	GetDuration(path string, defaultValue ...time.Duration) time.Duration
	// UnmarshalKey unmarshal a specific key from config into a struct.
	UnmarshalKey(key string, rawVal any) error
}
