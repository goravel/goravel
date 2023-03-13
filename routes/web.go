package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Web() {
	facades.Route.Get("/", func(ctx http.Context) {
		ctx.Response().Json(http.StatusOK, http.Json{
			"Hello": "Goravel",
		})
	})

	userController := controllers.NewUserController()
	facades.Route.Get("/users/{id}", userController.Show)
}
