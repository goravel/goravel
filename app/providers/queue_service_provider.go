package providers

import (
	"github.com/goravel/framework/contracts/foundation"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Register(app foundation.Application) {
	// Queue jobs can be registered when needed
	// or auto-discovered by the framework
}

func (receiver *QueueServiceProvider) Boot(app foundation.Application) {

}
