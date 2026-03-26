package bootstrap

import (
	"github.com/goravel/framework/cache"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/http"
	"github.com/goravel/framework/log"
	"github.com/goravel/framework/route"
	"github.com/goravel/framework/session"
	"github.com/goravel/framework/validation"
	"github.com/goravel/framework/view"
	"github.com/goravel/gin"
)

func Providers() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&http.ServiceProvider{},
		&log.ServiceProvider{},
		&cache.ServiceProvider{},
		&session.ServiceProvider{},
		&validation.ServiceProvider{},
		&view.ServiceProvider{},
		&route.ServiceProvider{},
		&gin.ServiceProvider{},
	}
}
