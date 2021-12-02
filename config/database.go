package config

import (
	"github.com/goravel/framework/support/facades"
)

func init() {
	config := facades.Config
	config.Add("database", map[string]interface{}{
		"default": config.Env("DB_CONNECTION", "mysql"),
		"connections": map[string]interface{}{
			"mysql": map[string]interface{}{
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "nft"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
			},
		},
	})
}
