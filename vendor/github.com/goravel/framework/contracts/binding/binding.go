package binding

const (
	Artisan     = "goravel.artisan"
	Auth        = "goravel.auth"
	Cache       = "goravel.cache"
	Config      = "goravel.config"
	Crypt       = "goravel.crypt"
	DB          = "goravel.db"
	Event       = "goravel.event"
	Gate        = "goravel.gate"
	Grpc        = "goravel.grpc"
	Hash        = "goravel.hash"
	Http        = "goravel.http"
	Lang        = "goravel.lang"
	Log         = "goravel.log"
	Mail        = "goravel.mail"
	Orm         = "goravel.orm"
	Process     = "goravel.process"
	Queue       = "goravel.queue"
	RateLimiter = "goravel.rate_limiter"
	Route       = "goravel.route"
	Schedule    = "goravel.schedule"
	Schema      = "goravel.schema"
	Seeder      = "goravel.seeder"
	Session     = "goravel.session"
	Storage     = "goravel.storage"
	Telemetry   = "goravel.telemetry"
	Testing     = "goravel.testing"
	Validation  = "goravel.validation"
	View        = "goravel.view"
)

type Relationship struct {
	// The bindings that are binded in the service provider.
	Bindings []string
	// The dependencies required by the service provider.
	Dependencies []string
	// The bindings that the service provider can provide for.
	ProvideFor []string
}

type Driver struct {
	// The name of the driver.
	Name string
	// A brief description of the driver.
	Description string
	// The package address of the driver.
	Package string
}

type Info struct {
	// The package path of the binding's service provider.
	PkgPath string
	// The dependencies required by the binding.
	Dependencies []string
	// A brief description of the binding.
	Description string
	// The drivers supported for the binding, some bindings cannot be used without specific drivers.
	// Eg: The Route facade needs goravel/gin or goravel/fiber driver.
	Drivers []Driver
	// Other bindings that should be installed together with this binding.
	// They do not have to be dependencies of this binding, but we want to install them together for better developer experience.
	// Eg: The Schema facade can be installed together with the Orm facade.
	InstallTogether []string
	// Indicates whether the binding is a base binding that should be registered by default.
	IsBase bool
}

var (
	Bindings = map[string]Info{
		Artisan: {
			Description: "The CLI tool that comes with Goravel for interacting with the command line.",
			PkgPath:     "github.com/goravel/framework/console",
			IsBase:      true,
		},
		Config: {
			Description: "Gets and sets configuration values.",
			PkgPath:     "github.com/goravel/framework/config",
			IsBase:      true,
		},
		Process: {
			Description: "Executes and manages external processes with concurrency support.",
			PkgPath:     "github.com/goravel/framework/process",
			IsBase:      true,
		},
		Auth: {
			Description: "Provides support for JWT and Session drivers.",
			PkgPath:     "github.com/goravel/framework/auth",
			Dependencies: []string{
				Cache,
				Config,
				Log,
				Orm,
			},
		},
		Cache: {
			Description: "Gets and sets cached items.",
			PkgPath:     "github.com/goravel/framework/cache",
			Dependencies: []string{
				Config,
				Log,
			},
			Drivers: []Driver{
				{
					Name:        "Memory",
					Description: "default",
					Package:     "memory",
				},
				{
					Name:    "Redis",
					Package: "github.com/goravel/redis",
				},
			},
		},
		Crypt: {
			Description: "Provides encryption and decryption services.",
			PkgPath:     "github.com/goravel/framework/crypt",
			Dependencies: []string{
				Config,
			},
		},
		DB: {
			Description: "Database management and query builder.",
			PkgPath:     "github.com/goravel/framework/database",
			Dependencies: []string{
				Config,
				Log,
			},
			Drivers: []Driver{
				{
					Name:    "Postgres",
					Package: "github.com/goravel/postgres",
				},
				{
					Name:    "MySQL",
					Package: "github.com/goravel/mysql",
				},
				{
					Name:    "SQLServer",
					Package: "github.com/goravel/sqlserver",
				},
				{
					Name:    "SQLite",
					Package: "github.com/goravel/sqlite",
				},
			},
			InstallTogether: []string{
				Schema,
			},
		},
		Event: {
			Description: "Provides a simple observer pattern implementation.",
			PkgPath:     "github.com/goravel/framework/event",
			Dependencies: []string{
				Queue,
			},
		},
		Gate: {
			Description: "An easy-to-use authorization feature to manage user actions on resources.",
			PkgPath:     "github.com/goravel/framework/auth",
			Dependencies: []string{
				Cache,
				Orm,
			},
		},
		Grpc: {
			Description: "Provides gRPC server and client support.",
			PkgPath:     "github.com/goravel/framework/grpc",
			Dependencies: []string{
				Config,
			},
		},
		Hash: {
			Description: "Provides secure Argon2id and Bcrypt hashing for storing user passwords.",
			PkgPath:     "github.com/goravel/framework/hash",
			Dependencies: []string{
				Config,
			},
		},
		Http: {
			Description: "An easy-to-use, expressive, and minimalist API built on the standard net/http library.",
			PkgPath:     "github.com/goravel/framework/http",
			Dependencies: []string{
				Config,
			},
		},
		Lang: {
			Description: "Provides localization support for multiple languages.",
			PkgPath:     "github.com/goravel/framework/translation",
			Dependencies: []string{
				Log,
			},
		},
		Log: {
			Description: "Provides logging capabilities with support for multiple channels and formats.",
			PkgPath:     "github.com/goravel/framework/log",
			Dependencies: []string{
				Config,
			},
		},
		Mail: {
			Description: "A clean, simple API over popular email services.",
			PkgPath:     "github.com/goravel/framework/mail",
			Dependencies: []string{
				Config,
				Queue,
			},
		},
		Orm: {
			Description: "An elegant ORM for Go, inspired by Eloquent.",
			PkgPath:     "github.com/goravel/framework/database",
			Dependencies: []string{
				Config,
				Log,
			},
			Drivers: []Driver{
				{
					Name:    "Postgres",
					Package: "github.com/goravel/postgres",
				},
				{
					Name:    "MySQL",
					Package: "github.com/goravel/mysql",
				},
				{
					Name:    "SQLServer",
					Package: "github.com/goravel/sqlserver",
				},
				{
					Name:    "SQLite",
					Package: "github.com/goravel/sqlite",
				},
			},
			InstallTogether: []string{
				Schema,
			},
		},
		Queue: {
			Description: "A solution by allowing you to create queued jobs that can run in the background.",
			PkgPath:     "github.com/goravel/framework/queue",
			Dependencies: []string{
				Config,
				Cache,
				DB,
				Log,
				Schema,
			},
			Drivers: []Driver{
				{
					Name:        "Sync",
					Description: "default",
					Package:     "sync",
				},
				{
					Name:    "Database",
					Package: "database",
				},
				{
					Name:    "Redis",
					Package: "github.com/goravel/redis",
				},
			},
		},
		RateLimiter: {
			Description: "Provides a simple and efficient way to limit the rate of incoming requests.",
			PkgPath:     "github.com/goravel/framework/http",
			Dependencies: []string{
				Cache,
				Log,
			},
		},
		Route: {
			Description: "Routing system, which supports multiple web frameworks.",
			PkgPath:     "github.com/goravel/framework/route",
			Dependencies: []string{
				Config,
				Http,
				RateLimiter,
				Session,
				Validation,
				View,
			},
			Drivers: []Driver{
				{
					Name:        "Gin",
					Description: "Gin is a high-performance HTTP web framework written in Go.",
					Package:     "github.com/goravel/gin",
				},
				{
					Name:        "Fiber",
					Description: "Fiber is an Express inspired web framework built on top of Fasthttp.",
					Package:     "github.com/goravel/fiber",
				},
			},
		},
		Schedule: {
			Description: "A fresh approach to managing scheduled tasks on your server.",
			PkgPath:     "github.com/goravel/framework/schedule",
			Dependencies: []string{
				Artisan,
				Cache,
				Config,
				Log,
			},
		},
		Schema: {
			Description: "Database schema builder and migration system.",
			PkgPath:     "github.com/goravel/framework/database",
			Dependencies: []string{
				Config,
				Log,
				Orm,
			},
		},
		Seeder: {
			Description: "Database seeding system to populate your database with test data.",
			PkgPath:     "github.com/goravel/framework/database",
		},
		Session: {
			Description: "Enables you to store user information across multiple requests.",
			PkgPath:     "github.com/goravel/framework/session",
			Dependencies: []string{
				Config,
			},
			Drivers: []Driver{
				{
					Name:        "File",
					Description: "default",
					Package:     "file",
				},
				{
					Name:    "Redis",
					Package: "github.com/goravel/redis",
				},
			},
		},
		Storage: {
			Description: "Provides a unified API for interacting with various file storage systems.",
			PkgPath:     "github.com/goravel/framework/filesystem",
			Dependencies: []string{
				Config,
			},
			Drivers: []Driver{
				{
					Name:        "Local",
					Description: "default",
					Package:     "local",
				},
				{
					Name:        "S3",
					Description: "power by Amazon",
					Package:     "github.com/goravel/s3",
				},
				{
					Name:        "OSS",
					Description: "power by Alibaba Cloud",
					Package:     "github.com/goravel/oss",
				},
				{
					Name:        "cos",
					Description: "power by Tencent Cloud",
					Package:     "github.com/goravel/cos",
				},
				{
					Name:        "MinIO",
					Description: "a high-performance, S3-compatible object storage solution",
					Package:     "github.com/goravel/minio",
				},
			},
		},
		// Telemetry: {
		// 	Description: "Provides distributed tracing with OpenTelemetry.",
		// 	PkgPath:     "github.com/goravel/framework/telemetry",
		// 	Dependencies: []string{
		// 		Config,
		// 		Log,
		// 	},
		// },
		Testing: {
			Description: "Provides tools for testing your application.",
			PkgPath:     "github.com/goravel/framework/testing",
			Dependencies: []string{
				Artisan,
				Cache,
				Config,
				Orm,
				Process,
				Route,
				Session,
			},
		},
		Validation: {
			Description: "Provides validation services for incoming data.",
			PkgPath:     "github.com/goravel/framework/validation",
		},
		View: {
			Description: "Provides a simple yet powerful templating engine.",
			PkgPath:     "github.com/goravel/framework/view",
		},
	}
)
