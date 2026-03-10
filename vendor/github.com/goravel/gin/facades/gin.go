package facades

import (
	"log"

	"github.com/goravel/gin"

	"github.com/goravel/framework/contracts/route"
)

func Route(driver string) route.Route {
	instance, err := gin.App.MakeWith(gin.BindingRoute, map[string]any{
		"driver": driver,
	})

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*gin.Route)
}
