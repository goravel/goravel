package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("grpc", map[string]any{
		// Grpc Configuration
		//
		// Configure your server host
		"host": config.Env("GRPC_HOST", ""),

		// Configure your client host and interceptors.
		// Interceptors can be the group name of UnaryClientInterceptorGroups in app/grpc/kernel.go.
		"clients": map[string]any{
			//"user": map[string]any{
			//	"host":         config.Env("GRPC_USER_HOST", ""),
			//	"interceptors": []string{},
			//},
		},
	})
}
