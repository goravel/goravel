package http

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/http/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel *Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		middleware.Cors(),
	}
}
