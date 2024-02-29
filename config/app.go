package config

import (
	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/cache"
	"github.com/goravel/framework/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/crypt"
	"github.com/goravel/framework/database"
	"github.com/goravel/framework/event"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/filesystem"
	"github.com/goravel/framework/grpc"
	"github.com/goravel/framework/hash"
	"github.com/goravel/framework/http"
	"github.com/goravel/framework/log"
	"github.com/goravel/framework/mail"
	"github.com/goravel/framework/queue"
	"github.com/goravel/framework/route"
	"github.com/goravel/framework/schedule"
	"github.com/goravel/framework/session"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/testing"
	"github.com/goravel/framework/translation"
	"github.com/goravel/framework/validation"
	"github.com/goravel/gin"

	"goravel/app/providers"
)

// Boot Start all init methods of the current folder to bootstrap all config.
func Boot() {}

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		// Application Name
		//
		// This value is the name of your application. This value is used when the
		// framework needs to place the application's name in a notification or
		// any other location as required by the application or its packages.
		"name": config.Env("APP_NAME", "Goravel"),

		// Application Environment
		//
		// This value determines the "environment" your application is currently
		// running in. This may determine how you prefer to configure various
		// services the application utilizes. Set this in your ".env" file.
		"env": config.Env("APP_ENV", "production"),

		// Application Debug Mode
		"debug": config.Env("APP_DEBUG", false),

		// Application Timezone
		//
		// Here you may specify the default timezone for your application.
		// Example: UTC, Asia/Shanghai
		// More: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
		"timezone": carbon.UTC,

		// Application Locale Configuration
		//
		// The application locale determines the default locale that will be used
		// by the translation service provider.You are free to set this value
		// to any of the locales which will be supported by the application.
		"locale": "en",

		// Application Fallback Locale
		//
		// The fallback locale determines the locale to use when the current one
		// is not available.You may change the value to correspond to any of
		// the language folders that are provided through your application.
		"fallback_locale": "en",

		// Encryption Key
		//
		// 32 character string, otherwise these encrypted strings
		// will not be safe. Please do this before deploying an application!
		"key": config.Env("APP_KEY", ""),

		// Autoload service providers
		//
		// The service providers listed here will be automatically loaded on the
		// request to your application. Feel free to add your own services to
		// this array to grant expanded functionality to your applications.
		"providers": []foundation.ServiceProvider{
			&log.ServiceProvider{},
			&console.ServiceProvider{},
			&database.ServiceProvider{},
			&cache.ServiceProvider{},
			&http.ServiceProvider{},
			&route.ServiceProvider{},
			&schedule.ServiceProvider{},
			&event.ServiceProvider{},
			&queue.ServiceProvider{},
			&grpc.ServiceProvider{},
			&mail.ServiceProvider{},
			&auth.ServiceProvider{},
			&hash.ServiceProvider{},
			&crypt.ServiceProvider{},
			&filesystem.ServiceProvider{},
			&validation.ServiceProvider{},
			&session.ServiceProvider{},
			&translation.ServiceProvider{},
			&testing.ServiceProvider{},
			&providers.AppServiceProvider{},
			&providers.AuthServiceProvider{},
			&providers.RouteServiceProvider{},
			&providers.GrpcServiceProvider{},
			&providers.ConsoleServiceProvider{},
			&providers.QueueServiceProvider{},
			&providers.EventServiceProvider{},
			&providers.ValidationServiceProvider{},
			&providers.DatabaseServiceProvider{},
			&gin.ServiceProvider{},
		},
	})
}
