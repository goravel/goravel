package routes

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Web() {
	facades.Route.Get("/", func(request http.Request) {
		facades.Response.Json(200, http.Json{
			"Hello": "Goravel",
		})
	})

	userController := controllers.NewUserController()
	facades.Route.Get("/{id}", func(request http.Request) {
		fmt.Printf("hwb---- %+v\n", request.Input("id"))
	})

	facades.Route.Prefix("dd").Get("/{id}/{name}", userController.Show)

	facades.Route.Group(func(route route.Route) {
		route.Get("aa/{id}", userController.Show)
	})

	facades.Route.Prefix("abc").Group(func(route route.Route) {
		route.Get("hwb/{id}", userController.Show)
	})

	facades.Route.Any("any", userController.Show)
}
