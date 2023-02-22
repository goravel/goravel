package http

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/http/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel *Kernel) Middleware() []http.Middleware {
	middlewares := []http.Middleware{
		middleware.Cors(),
	}
	if facades.Config.GetBool("tls.enabled") {
		middlewares = append(middlewares, middleware.Tls(facades.Config.GetString("app.host")))
	}
	return middlewares
}
