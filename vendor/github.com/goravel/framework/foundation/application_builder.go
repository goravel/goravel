package foundation

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	contractsconfiguration "github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/contracts/validation"
)

func Setup() foundation.ApplicationBuilder {
	return NewApplicationBuilder(App)
}

type ApplicationBuilder struct {
	app                        foundation.Application
	callback                   func()
	commands                   func() []console.Command
	config                     func()
	configuredServiceProviders func() []foundation.ServiceProvider
	eventToListeners           func() map[event.Event][]event.Listener
	filters                    func() []validation.Filter
	grpcClientInterceptors     func() map[string][]grpc.UnaryClientInterceptor
	grpcClientStatsHandlers    func() map[string][]stats.Handler
	grpcServerInterceptors     func() []grpc.UnaryServerInterceptor
	grpcServerStatsHandlers    func() []stats.Handler
	jobs                       func() []queue.Job
	middleware                 func(middleware contractsconfiguration.Middleware)
	migrations                 func() []schema.Migration
	paths                      func(paths contractsconfiguration.Paths)
	routes                     func()
	rules                      func() []validation.Rule
	runners                    func() []foundation.Runner
	schedule                   func() []schedule.Event
	seeders                    func() []seeder.Seeder
}

func NewApplicationBuilder(app foundation.Application) *ApplicationBuilder {
	return &ApplicationBuilder{
		app: app,
	}
}

func (r *ApplicationBuilder) Create() foundation.Application {
	return r.app.SetBuilder(r).Build()
}

func (r *ApplicationBuilder) WithCallback(callback func()) foundation.ApplicationBuilder {
	r.callback = callback

	return r
}

func (r *ApplicationBuilder) WithCommands(fn func() []console.Command) foundation.ApplicationBuilder {
	r.commands = fn

	return r
}

func (r *ApplicationBuilder) WithConfig(fn func()) foundation.ApplicationBuilder {
	r.config = fn

	return r
}

func (r *ApplicationBuilder) WithEvents(fn func() map[event.Event][]event.Listener) foundation.ApplicationBuilder {
	r.eventToListeners = fn

	return r
}

func (r *ApplicationBuilder) WithFilters(fn func() []validation.Filter) foundation.ApplicationBuilder {
	r.filters = fn

	return r
}

func (r *ApplicationBuilder) WithGrpcClientInterceptors(fn func() map[string][]grpc.UnaryClientInterceptor) foundation.ApplicationBuilder {
	r.grpcClientInterceptors = fn

	return r
}

func (r *ApplicationBuilder) WithGrpcClientStatsHandlers(fn func() map[string][]stats.Handler) foundation.ApplicationBuilder {
	r.grpcClientStatsHandlers = fn

	return r
}

func (r *ApplicationBuilder) WithGrpcServerInterceptors(fn func() []grpc.UnaryServerInterceptor) foundation.ApplicationBuilder {
	r.grpcServerInterceptors = fn

	return r
}

func (r *ApplicationBuilder) WithGrpcServerStatsHandlers(fn func() []stats.Handler) foundation.ApplicationBuilder {
	r.grpcServerStatsHandlers = fn

	return r
}

func (r *ApplicationBuilder) WithJobs(fn func() []queue.Job) foundation.ApplicationBuilder {
	r.jobs = fn

	return r
}

func (r *ApplicationBuilder) WithMiddleware(fn func(handler contractsconfiguration.Middleware)) foundation.ApplicationBuilder {
	r.middleware = fn

	return r
}

func (r *ApplicationBuilder) WithMigrations(fn func() []schema.Migration) foundation.ApplicationBuilder {
	r.migrations = fn

	return r
}

func (r *ApplicationBuilder) WithPaths(fn func(paths contractsconfiguration.Paths)) foundation.ApplicationBuilder {
	r.paths = fn

	return r
}

func (r *ApplicationBuilder) WithProviders(fn func() []foundation.ServiceProvider) foundation.ApplicationBuilder {
	r.configuredServiceProviders = fn

	return r
}

func (r *ApplicationBuilder) WithRouting(fn func()) foundation.ApplicationBuilder {
	r.routes = fn

	return r
}

func (r *ApplicationBuilder) WithRules(fn func() []validation.Rule) foundation.ApplicationBuilder {
	r.rules = fn

	return r
}

func (r *ApplicationBuilder) WithRunners(fn func() []foundation.Runner) foundation.ApplicationBuilder {
	r.runners = fn

	return r
}

func (r *ApplicationBuilder) WithSchedule(fn func() []schedule.Event) foundation.ApplicationBuilder {
	r.schedule = fn

	return r
}

func (r *ApplicationBuilder) WithSeeders(fn func() []seeder.Seeder) foundation.ApplicationBuilder {
	r.seeders = fn

	return r
}
