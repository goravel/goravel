package main

import (
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/facades"

	"goravel/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	//// Create a channel to listen for OS signals
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//
	//// Start http server by facades.Route().
	//go func() {
	//	if err := facades.Route().Run(); err != nil {
	//		facades.Log().Errorf("Route Run error: %v", err)
	//	}
	//}()
	//
	//// Listen for the OS signal
	//go func() {
	//	<-quit
	//	if err := facades.Route().Shutdown(); err != nil {
	//		facades.Log().Errorf("Route Shutdown error: %v", err)
	//	}
	//
	//	os.Exit(0)
	//}()
	//
	//select {}

	facades.Schema().Create("users", func(table migration.Blueprint) {
		table.ID("id")
		table.String("migration")
		table.Integer("batch")
	})
}
