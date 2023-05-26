package sms

import (
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"
)

const Binding = "sms"

var App foundationcontract.Application

type ServiceProvider struct {
	foundation.ServiceProvider
}

func (receiver *ServiceProvider) Register(app foundationcontract.Application) {
	//receiver.Publishes(map[string]string{"a": "b"})
	App = app

	app.Bind(Binding, func() (any, error) {
		return NewSms(app.MakeConfig()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundationcontract.Application) {

}
