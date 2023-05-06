package main

import (
	"goravel/bootstrap"
	smsfacades "goravel/packages/sms/facades"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Start http server by facades.Route.
	//go func() {
	//	if err := facades.Route.Run(); err != nil {
	//		facades.Log.Errorf("Route run error: %v", err)
	//	}
	//}()

	smsfacades.Sms().Send()

	select {}
}
