package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
	"goravel/app/http/controllers"
)

func Web() {
	facades.Route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "Goravel",
		})
	})

	facades.Route.GET("/user", controllers.UserController{}.Show)
}
