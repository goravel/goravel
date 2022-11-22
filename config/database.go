package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("database", map[string]interface{}{
		// Default database connection name, only support Mysql now.
		"default": config.Env("DB_CONNECTION", "mysql"),

		// Database connections
		"connections": map[string]interface{}{
			"mysql": map[string]interface{}{
				"driver":   "mysql",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
				"loc":      "Local",
			},
			"postgresql": map[string]interface{}{
				"driver":   "postgresql",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"sslmode":  "disable",
				"timezone": "UTC", //Asia/Shanghai
			},
			"sqlite": map[string]interface{}{
				"driver":   "sqlite",
				"database": config.Env("DB_DATABASE", "forge"),
			},
			"sqlserver": map[string]interface{}{
				"driver":   "sqlserver",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
			},
		},

		// Migration Repository Table
		//
		// This table keeps track of all the migrations that have already run for
		// your application. Using this information, we can determine which of
		// the migrations on disk haven't actually been run in the database.
		"migrations": "migrations",

		// Redis Databases
		//
		// Redis is an open source, fast, and advanced key-value store that also
		// provides a richer body of commands than a typical key-value system
		// such as APC or Memcached.
		"redis": map[string]interface{}{
			"default": map[string]interface{}{
				"host":     config.Env("REDIS_HOST", ""),
				"password": config.Env("REDIS_PASSWORD", ""),
				"port":     config.Env("REDIS_PORT", 6379),
				"database": config.Env("REDIS_DB", 0),
			},
		},
	})
}
