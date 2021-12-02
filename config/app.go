package config

import (
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/facades"
	"goravel/app/providers"
)

//Boot Start all init methods of the current folder
func Boot() {}

func init() {
	config := facades.Config
	config.Add("app", map[string]interface{}{
		"name":  config.Env("APP_NAME", "nft"),
		"env":   config.Env("APP_ENV", "production"),
		"debug": config.Env("APP_DEBUG", false),
		"key":   config.Env("APP_KEY", ""),
		"url":   config.Env("APP_URL", "http://localhost:3000"),
		"host":  config.Env("APP_HOST", "http://localhost:3000"),
		"providers": []support.ServiceProvider{
			&providers.RouteServiceProvider{},
		},
	})
}
