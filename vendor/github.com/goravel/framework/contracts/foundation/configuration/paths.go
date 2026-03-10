package configuration

type Paths interface {
	// App sets the path for the application directory, default is "app".
	App(path string) Paths
	// Bootstrap sets the path for the bootstrap directory, default is "bootstrap".
	Bootstrap(path string) Paths
	// Commands sets the path for the commands directory, default is "app/console/commands".
	Commands(path string) Paths
	// Config sets the path for the configuration directory, default is "config".
	Config(path string) Paths
	// Controllers sets the path for the controllers directory, default is "app/http/controllers".
	Controllers(path string) Paths
	// Database sets the path for the database directory, default is "database".
	Database(path string) Paths
	// Events sets the path for the events directory, default is "app/events".
	Events(path string) Paths
	// Facades sets the path for the facades directory, default is "app/facades".
	Facades(path string) Paths
	// Factories sets the path for the factories directory, default is "database/factories".
	Factories(path string) Paths
	// Filters sets the path for the filters directory, default is "app/http/filters".
	Filters(path string) Paths
	// Jobs sets the path for the jobs directory, default is "app/jobs".
	Jobs(path string) Paths
	// Lang sets the path for the language files directory, default is "lang".
	Lang(path string) Paths
	// Listeners sets the path for the listeners directory, default is "app/listeners".
	Listeners(path string) Paths
	// Mails sets the path for the mails directory, default is "app/mails".
	Mails(path string) Paths
	// Middleware sets the path for the middleware directory, default is "app/http/middleware".
	Middleware(path string) Paths
	// Migrations sets the path for the migrations directory, default is "database/migrations".
	Migrations(path string) Paths
	// Models sets the path for the models directory, default is "app/models".
	Models(path string) Paths
	// Observers sets the path for the observers directory, default is "app/observers".
	Observers(path string) Paths
	// Packages sets the path for the packages directory, default is "packages".
	Packages(path string) Paths
	// Policies sets the path for the policies directory, default is "app/policies".
	Policies(path string) Paths
	// Providers sets the path for the providers directory, default is "app/providers".
	Providers(path string) Paths
	// Public sets the path for the public directory, default is "public".
	Public(path string) Paths
	// Requests sets the path for the requests directory, default is "app/http/requests".
	Requests(path string) Paths
	// Resources sets the path for the resources directory, default is "resources".
	Resources(path string) Paths
	// Routes sets the path for the routes directory, default is "routes".
	Routes(path string) Paths
	// Rules sets the path for the rules directory, default is "app/rules".
	Rules(path string) Paths
	// Seeders sets the path for the seeders directory, default is "database/seeders".
	Seeders(path string) Paths
	// Storage sets the path for the storage directory, default is "storage".
	Storage(path string) Paths
	// Tests sets the path for the tests directory, default is "tests".
	Tests(path string) Paths
	// Views sets the path for the views directory, default is "resources/views".
	Views(path string) Paths
}
