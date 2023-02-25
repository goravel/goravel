package main

import (
	"github.com/goravel/framework/facades"

	"goravel/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Start http server by facades.Route.
	go func() {
		app_host := facades.Config.GetString("app.host") + ":" + facades.Config.GetString("app.port")
		if err := facades.Route.Run(app_host); err != nil {
			facades.Log.Errorf("Route run error: %v", err)
		}
	}()

	select {}
}
