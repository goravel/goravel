package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("database", map[string]any{
		// Default database connection name
		"default": config.Env("DB_CONNECTION", "mysql"),

		// Database connections
		"connections": map[string]any{
			"mysql": map[string]any{
				"driver":   "mysql",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", 3306),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
				"loc":      "Local",
				"prefix":   "",
				"singular": false, // Table name is singular
			},
			"postgresql": map[string]any{
				"driver":   "postgresql",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", 5432),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"sslmode":  "disable",
				"timezone": "UTC", //Asia/Shanghai
				"prefix":   "",
				"singular": false, // Table name is singular
			},
			"sqlite": map[string]any{
				"driver":   "sqlite",
				"database": config.Env("DB_DATABASE", "forge"),
				"prefix":   "",
				"singular": false, // Table name is singular
			},
			"sqlserver": map[string]any{
				"driver":   "sqlserver",
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", 1433),
				"database": config.Env("DB_DATABASE", "forge"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
				"prefix":   "",
				"singular": false, // Table name is singular
			},
		},

		// Set pool configuration
		"pool": map[string]any{
			// Sets the maximum number of connections in the idle
			// connection pool.
			//
			// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
			// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
			//
			// If n <= 0, no idle connections are retained.
			"max_idle_conns": 10,
			// Sets the maximum number of open connections to the database.
			//
			// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
			// MaxIdleConns, then MaxIdleConns will be reduced to match the new
			// MaxOpenConns limit.
			//
			// If n <= 0, then there is no limit on the number of open connections.
			"max_open_conns": 100,
			// Sets the maximum amount of time a connection may be idle.
			//
			// Expired connections may be closed lazily before reuse.
			//
			// If d <= 0, connections are not closed due to a connection's idle time.
			// Unit: Second
			"conn_max_idletime": 3600,
			// Sets the maximum amount of time a connection may be reused.
			//
			// Expired connections may be closed lazily before reuse.
			//
			// If d <= 0, connections are not closed due to a connection's age.
			// Unit: Second
			"conn_max_lifetime": 3600,
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
		"redis": map[string]any{
			"default": map[string]any{
				"host":     config.Env("REDIS_HOST", ""),
				"password": config.Env("REDIS_PASSWORD", ""),
				"port":     config.Env("REDIS_PORT", 6379),
				"database": config.Env("REDIS_DB", 0),
			},
		},
	})
}
