package config

import (
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/convert"
	"github.com/goravel/framework/support/file"
)

var _ config.Config = &Application{}

type Application struct {
	vip *viper.Viper
}

func NewApplication(envFilePath string) *Application {
	app := &Application{}
	app.vip = viper.New()
	app.vip.AutomaticEnv()

	if file.Exists(envFilePath) {
		app.vip.SetConfigType("env")
		app.vip.SetConfigFile(envFilePath)

		if err := app.vip.ReadInConfig(); err != nil {
			color.Errorln("Invalid Config error: " + err.Error())
			os.Exit(0)
		}
	}

	appKey := app.Env("APP_KEY")
	if !support.DontVerifyAppKey {
		if appKey == nil {
			color.Errorln("Please initialize APP_KEY first.")
			color.Default().Println("Create a .env file and run command: go run . artisan key:generate")
			color.Default().Println("Or set a system variable: APP_KEY={32-bit number} go run .")
			os.Exit(0)
		}

		if len(appKey.(string)) != 32 {
			color.Errorln("Invalid APP_KEY, the length must be 32, please reset it.")
			color.Warningln("Example command: \ngo run . artisan key:generate")
			os.Exit(0)
		}
	}

	return app
}

// Env Get config from env.
func (app *Application) Env(envName string, defaultValue ...any) any {
	value := app.Get(envName, defaultValue...)
	if cast.ToString(value) == "" {
		return convert.Default(defaultValue...)
	}

	return value
}

// EnvString get string value from env with optional default.
func (app *Application) EnvString(envName string, defaultValue ...string) string {
	value := app.Env(envName)
	if cast.ToString(value) == "" {
		return convert.Default(defaultValue...)
	}
	return cast.ToString(value)
}

// EnvBool get bool value from env with optional default.
func (app *Application) EnvBool(envName string, defaultValue ...bool) bool {
	value := app.Env(envName)
	// If no value and a default provided, return default
	if cast.ToString(value) == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return cast.ToBool(value)
}

// Add config to application.
func (app *Application) Add(name string, configuration any) {
	app.vip.Set(name, configuration)
}

// Get config from application.
func (app *Application) Get(path string, defaultValue ...any) any {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.Get(path)
}

// GetString get string type config from application.
func (app *Application) GetString(path string, defaultValue ...string) string {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.GetString(path)
}

// GetInt get int type config from application.
func (app *Application) GetInt(path string, defaultValue ...int) int {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.GetInt(path)
}

// GetBool get bool type config from application.
func (app *Application) GetBool(path string, defaultValue ...bool) bool {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.GetBool(path)
}

// GetDuration get time.Duration type config from application
func (app *Application) GetDuration(path string, defaultValue ...time.Duration) time.Duration {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.GetDuration(path)
}

// UnmarshalKey unmarshal a specific key from config into a struct.
func (app *Application) UnmarshalKey(key string, rawVal any) error {
	return app.vip.UnmarshalKey(key, rawVal, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
}
