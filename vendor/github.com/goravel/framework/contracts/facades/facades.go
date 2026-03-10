package facades

import "github.com/goravel/framework/contracts/binding"

const (
	Artisan     = "Artisan"
	Auth        = "Auth"
	Cache       = "Cache"
	Config      = "Config"
	Crypt       = "Crypt"
	DB          = "DB"
	Event       = "Event"
	Gate        = "Gate"
	Grpc        = "Grpc"
	Hash        = "Hash"
	Http        = "Http"
	Lang        = "Lang"
	Log         = "Log"
	Mail        = "Mail"
	Process     = "Process"
	Orm         = "Orm"
	Queue       = "Queue"
	RateLimiter = "RateLimiter"
	Route       = "Route"
	Schedule    = "Schedule"
	Schema      = "Schema"
	Seeder      = "Seeder"
	Session     = "Session"
	Storage     = "Storage"
	Telemetry   = "Telemetry"
	Testing     = "Testing"
	Validation  = "Validation"
	View        = "View"
)

var FacadeToBinding = map[string]string{
	Artisan:     binding.Artisan,
	Auth:        binding.Auth,
	Cache:       binding.Cache,
	Config:      binding.Config,
	Crypt:       binding.Crypt,
	DB:          binding.DB,
	Event:       binding.Event,
	Gate:        binding.Gate,
	Grpc:        binding.Grpc,
	Hash:        binding.Hash,
	Http:        binding.Http,
	Lang:        binding.Lang,
	Log:         binding.Log,
	Mail:        binding.Mail,
	Orm:         binding.Orm,
	Process:     binding.Process,
	Queue:       binding.Queue,
	RateLimiter: binding.RateLimiter,
	Route:       binding.Route,
	Schedule:    binding.Schedule,
	Schema:      binding.Schema,
	Seeder:      binding.Seeder,
	Session:     binding.Session,
	Storage:     binding.Storage,
	Telemetry:   binding.Telemetry,
	Testing:     binding.Testing,
	Validation:  binding.Validation,
	View:        binding.View,
}
