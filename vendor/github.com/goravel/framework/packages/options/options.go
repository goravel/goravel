package options

import "github.com/goravel/framework/contracts/packages/modify"

// Driver sets the driver option for the modify.Apply function.
func Driver(driver string) modify.Option {
	return func(options map[string]any) {
		options["driver"] = driver
	}
}

// Facade sets the facade option for the modify.Apply function.
func Facade(facade string) modify.Option {
	return func(options map[string]any) {
		options["facade"] = facade
	}
}

// Force sets the force option for the modify.Apply function.
func Force(force bool) modify.Option {
	return func(options map[string]any) {
		options["force"] = force
	}
}
