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
		if facades.Config.GetBool("tls.enabled") {
			certFile := facades.Storage.Disk("local").Path(facades.Config.GetString("tls.ssl.cert_file"))
			keyFile := facades.Storage.Disk("local").Path(facades.Config.GetString("tls.ssl.key_file"))
			if err := facades.Route.RunTLS(facades.Config.GetString("app.host"), certFile, keyFile); err != nil {
				facades.Log.Errorf("Route runTLS error: %v", err)
			}
		} else {
			if err := facades.Route.Run(facades.Config.GetString("app.host")); err != nil {
				facades.Log.Errorf("Route run error: %v", err)
			}
		}
	}()

	select {}
}
