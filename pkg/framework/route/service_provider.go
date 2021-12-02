package route

import (
	"github.com/goravel/framework/support/facades"
)

type ServiceProvider struct {
}

func (route *ServiceProvider) Boot() {
	app := Application{}
	facades.Route = app.Init()
}

func (route *ServiceProvider) Register() {

}
