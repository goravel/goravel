package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("limiter", map[string]any{
		// Limiter Store Deiver
		//
		// This option controls the default store diver that gets used
		// while using limit middleware.
		//
		// Supported Drivers: "memory", "redis"
		"store": "memory",

		// If you are using the "redis" limiter driver, you may specify a
		// connection that should be used to connect to the redis cache.
		// This connection should be defined in the "database" config file.
		"redis": "default",
	})
}
