package main

import (
	"fmt"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/time"

	"goravel/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Start http server by facades.Route.
	go func() {
		if err := facades.Route.Run(); err != nil {
			facades.Log.Errorf("Route run error: %v", err)
		}
	}()

	fmt.Println("hwb---", time.Now())
	facades.Hash.Make("123")
	fmt.Println("hwb---", time.Now())

	select {}
}
