package packages

import (
	"github.com/goravel/framework/contracts/packages/modify"
)

type Setup interface {
	// Execute runs the setup command based on the provided arguments.
	Execute()
	// Package returns the package information of application.
	Paths() Paths
	// Install adds the provided modifiers to be executed during installation.
	Install(modifiers ...modify.Apply) Setup
	// Uninstall adds the provided modifiers to be executed during uninstallation.
	Uninstall(modifiers ...modify.Apply) Setup
}

type Paths interface {
	// App returns the path for the app package, eg: goravel/app.
	App() Path
	// Bootstrap returns the path for the bootstrap package, eg: goravel/bootstrap.
	Bootstrap() Path
	// Config returns the path for the config package, eg: goravel/config.
	Config() Path
	// Database returns the path for the database package, eg: goravel/database.
	Database() Path
	// Facades returns the path for the facades package, eg: goravel/app/facades.
	Facades() Path
	// Lang returns the path for the lang package, eg: goravel/lang.
	Lang() Path
	// Main returns the path for the main package, eg: github.com/goravel/goravel.
	Main() Path
	// Migrations returns the path for the migrations package, eg: goravel/database/migrations.
	Migrations() Path
	// Models returns the path for the models package, eg: goravel/app/models.
	Models() Path
	// Module returns the path for the module package, eg: github.com/goravel/framework/auth.
	Module() Path
	// Public returns the path for the public package, eg: goravel/public.
	Public() Path
	// Resources returns the path for the resources package, eg: goravel/resources.
	Resources() Path
	// Routes returns the path for the routes package, eg: goravel/routes.
	Routes() Path
	// Storage returns the path for the storage package, eg: goravel/storage.
	Storage() Path
	// Tests returns the path for the tests package, eg: goravel/tests.
	Tests() Path
	// Views returns the path for the views package, eg: goravel/resources/views.
	Views() Path
}

type Path interface {
	// Abs returns the absolute path of the package.
	Abs(paths ...string) string
	// Package returns the sub-package name, or the main package name if no sub-package path is specified.
	Package() string
	// Import returns the sub-package import path, or the main package import path if no sub-package path is specified.
	Import() string
	// String returns the setting path of the package.
	String(paths ...string) string
}
