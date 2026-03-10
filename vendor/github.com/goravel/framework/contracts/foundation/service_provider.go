package foundation

import "github.com/goravel/framework/contracts/binding"

type ServiceProvider interface {
	// Register any application services.
	Register(app Application)
	// Boot any application services after register.
	Boot(app Application)
}

type ServiceProviderWithRelations interface {
	ServiceProvider

	// Relationship returns the service provider's relationship.
	Relationship() binding.Relationship
}

type ServiceProviderWithRunners interface {
	ServiceProvider

	// Runners returns the service provider's runners.
	Runners(app Application) []Runner
}
