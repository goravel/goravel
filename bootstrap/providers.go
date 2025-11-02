package bootstrap

import (
	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/cache"
	"github.com/goravel/framework/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/crypt"
	"github.com/goravel/framework/database"
	"github.com/goravel/framework/event"
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
	"github.com/goravel/framework/testing"
	"github.com/goravel/framework/translation"
	"github.com/goravel/framework/validation"
	"github.com/goravel/framework/view"
	"github.com/goravel/gin"
	"github.com/goravel/postgres"

	"goravel/app/providers"
)

func Providers() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&log.ServiceProvider{},
		&console.ServiceProvider{},
		&postgres.ServiceProvider{},
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
		&view.ServiceProvider{},
		&providers.AppServiceProvider{},
		&providers.AuthServiceProvider{},
		&providers.ConsoleServiceProvider{},
		&providers.QueueServiceProvider{},
		&providers.EventServiceProvider{},
		&providers.ValidationServiceProvider{},
		&providers.DatabaseServiceProvider{},
		&gin.ServiceProvider{},
	}
}
