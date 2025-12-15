package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	"goravel/app/facades"
)

func Web() {
	facades.Route().Static("/public", "./public")
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})
}
