package route

import (
	"github.com/goravel/framework/support/facades"
)

type ServiceProvider struct {
}

func (router *ServiceProvider) Boot() {
	gin := Gin{}
	facades.Route = gin.Init()
}

func (router *ServiceProvider) Register() {

}
