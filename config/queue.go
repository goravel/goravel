package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("queue", map[string]interface{}{
		// Default Queue Connection Name
		"default": config.Env("QUEUE_CONNECTION", "sync"),

		// Queue Connections
		//
		// Here you may configure the connection information for each server that is used by your application.
		// Drivers: "sync", "redis"
		"connections": map[string]interface{}{
			"sync": map[string]interface{}{
				"driver": "sync",
			},
			"redis": map[string]interface{}{
				"driver":     "redis",
				"connection": "default",
				"queue":      config.Env("REDIS_QUEUE", "default"),
			},
		},
	})
}
