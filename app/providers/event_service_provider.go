package providers

import (
	"github.com/goravel/framework/contracts/foundation"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Register(app foundation.Application) {
	// Events can be registered when needed
	// or auto-discovered by the framework
}

func (receiver *EventServiceProvider) Boot(app foundation.Application) {

}
