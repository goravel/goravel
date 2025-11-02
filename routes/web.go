package routes

import (
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		time.Sleep(5 * time.Second)
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})
}
