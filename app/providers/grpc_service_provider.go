package providers

import (
	"github.com/goravel/framework/contracts/foundation"

	"goravel/routes"
)

type GrpcServiceProvider struct {
}

func (receiver *GrpcServiceProvider) Register(app foundation.Application) {
	// GRPC interceptors can be configured in bootstrap or auto-discovered
}

func (receiver *GrpcServiceProvider) Boot(app foundation.Application) {
	// Add routes
	routes.Grpc()
}
