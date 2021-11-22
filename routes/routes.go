package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
)

func V1() {
	facades.Route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
