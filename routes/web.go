package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	middlewareThrottle "github.com/goravel/framework/http/middleware"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Middleware(middlewareThrottle.Throttle("user_throttle")).Get("/", func(ctx http.Context) http.Response {
		// 测试是否写入成功
		facades.Cache().Put("test", "test", 0)
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})
}
