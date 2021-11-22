package main

import (
	"github.com/goravel/framework/support/facades"
	"goravel/bootstrap"
)

func main() {
	bootstrap.Boot()

	facades.Route.Run(facades.Config.GetString("app.host"))
}
