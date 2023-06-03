package sms

import (
	"github.com/goravel/framework/contracts/console"
	foundationcontract "github.com/goravel/framework/contracts/foundation"

	"goravel/packages/sms/commands"
)

const Binding = "sms"

var App foundationcontract.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundationcontract.Application) {
	App = app

	app.Bind(Binding, func(app foundationcontract.Application) (any, error) {
		return NewSms(app.MakeConfig()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundationcontract.Application) {
	app.Publishes("./packages/sms", map[string]string{
		"config/sms.go": app.ConfigPath("sms.go"),
	})
	app.Commands([]console.Command{
		commands.NewListCommand(app.MakeArtisan()),
	})
}
