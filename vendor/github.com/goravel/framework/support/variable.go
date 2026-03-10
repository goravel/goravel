package support

type Paths struct {
	// The base directory path, default is "app".
	App string
	// The bootstrap directory path, default is "bootstrap".
	Bootstrap string
	// The commands directory path, default is "app/console/commands".
	Commands string
	// The config directory path, default is "config".
	Config string
	// The controllers directory path, default is "app/http/controllers".
	Controllers string
	// The database directory path, default is "database".
	Database string
	// The events directory path, default is "app/events".
	Events string
	// The facades directory path, default is "app/facades".
	Facades string
	// The factories directory path, default is "database/factories".
	Factories string
	// The filters directory path, default is "app/filters".
	Filters string
	// The jobs directory path, default is "app/jobs".
	Jobs string
	// The language files directory path, default is "lang".
	Lang string
	// The listeners directory path, default is "app/listeners".
	Listeners string
	// The mails directory path, default is "app/mails".
	Mails string
	// The middleware directory path, default is "app/http/middleware".
	Middleware string
	// The migrations directory path, default is "database/migrations".
	Migrations string
	// The models directory path, default is "app/models".
	Models string
	// The observers directory path, default is "app/observers".
	Observers string
	// The packages directory path, default is "packages".
	Packages string
	// The policies directory path, default is "app/policies".
	Policies string
	// The providers directory path, default is "app/providers".
	Providers string
	// The public directory path, default is "public".
	Public string
	// The requests directory path, default is "app/http/requests".
	Requests string
	// The resources directory path, default is "resources".
	Resources string
	// The routes directory path, default is "routes".
	Routes string
	// The rules directory path, default is "app/rules".
	Rules string
	// The seeders directory path, default is "database/seeders".
	Seeders string
	// The storage directory path, default is "storage".
	Storage string
	// The tests directory path, default is "tests".
	Tests string
	// The view directory path, default is "resources/views".
	Views string
}

type Configuration struct {
	Paths Paths
}

var (
	RelativePath = ""
	RootPath     = ""

	RuntimeMode = ""

	EnvFilePath          = ".env"
	EnvFileEncryptPath   = ".env.encrypted"
	EnvFileEncryptCipher = "AES-256-CBC"

	DontVerifyAppKey          = false
	DontVerifyAppKeyWhitelist = []string{"about", "list", "key:generate", "jwt:secret", "env:decrypt", "package:install"}

	Config = Configuration{
		Paths: Paths{
			App:         "app",
			Bootstrap:   "bootstrap",
			Commands:    "app/console/commands",
			Config:      "config",
			Controllers: "app/http/controllers",
			Database:    "database",
			Events:      "app/events",
			Facades:     "app/facades",
			Factories:   "database/factories",
			Filters:     "app/filters",
			Jobs:        "app/jobs",
			Lang:        "lang",
			Listeners:   "app/listeners",
			Mails:       "app/mails",
			Middleware:  "app/http/middleware",
			Migrations:  "database/migrations",
			Models:      "app/models",
			Observers:   "app/observers",
			Packages:    "packages",
			Policies:    "app/policies",
			Providers:   "app/providers",
			Public:      "public",
			Requests:    "app/http/requests",
			Resources:   "resources",
			Routes:      "routes",
			Rules:       "app/rules",
			Seeders:     "database/seeders",
			Storage:     "storage",
			Tests:       "tests",
			Views:       "resources/views",
		},
	}
)
