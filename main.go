package main

import (
	"github.com/goravel/framework/support/facades"
	"goravel/bootstrap"
)

func main() {
	//This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	//Start http server by facades.Route.
	facades.Route.Run(facades.Config.GetString("app.host"))
}
