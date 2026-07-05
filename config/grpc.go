package config

import (
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("grpc", map[string]any{
		// Configure your server host
		"host": config.Env("GRPC_HOST"),

		// Configure your server port
		"port": config.Env("GRPC_PORT"),

		// Configure clients which the framework will connect to
		"clients": map[string]any{
			//"user": map[string]any{
			//	"host":           config.Env("GRPC_USER_HOST"),
			//	"port":           config.Env("GRPC_USER_PORT"),
			//  "credentials":    config.Env("GRPC_USER_CREDENTIALS"),
			//  // the group name of UnaryClientInterceptorGroups
			//	"interceptors":   []string{},
            //  "stats_handlers": []string{},
			//},
		},
	})
}
