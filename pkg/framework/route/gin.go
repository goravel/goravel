package route

import (
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
)

type Gin struct {
}

func (g *Gin) Init() *gin.Engine {
	if facades.Config.Env("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	route := gin.New()

	return route
}
