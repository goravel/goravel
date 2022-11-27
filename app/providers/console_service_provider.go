package providers

import (
	"goravel/app/console"

	"github.com/goravel/framework/facades"
)

type ConsoleServiceProvider struct {
}

func (receiver *ConsoleServiceProvider) Register() {
	kernel := console.Kernel{}
	facades.Schedule.Register(kernel.Schedule())
	facades.Artisan.Register(kernel.Commands())
}

func (receiver *ConsoleServiceProvider) Boot() {

}
