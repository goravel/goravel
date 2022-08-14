package providers

import (
	"github.com/goravel/framework/support/facades"
	"goravel/app/console"
)

type ConsoleServiceProvider struct {
}

func (receiver *ConsoleServiceProvider) Boot() {

}

func (receiver *ConsoleServiceProvider) Register() {
	facades.Schedule.Register(console.Kernel{}.Schedule())
	facades.Artisan.Register(console.Kernel{}.Commands())
}
