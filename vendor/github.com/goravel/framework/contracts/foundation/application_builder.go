package foundation

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/contracts/validation"
)

type ApplicationBuilder interface {
	// Create a new application instance after configuring.
	Create() Application
	// WithCallback sets a callback function to be called during application creation.
	WithCallback(func()) ApplicationBuilder
	// WithCommands sets the application's commands.
	WithCommands(func() []console.Command) ApplicationBuilder
	// WithConfig sets a callback function to configure the application.
	WithConfig(func()) ApplicationBuilder
	// WithEvents sets event listeners for the application.
	WithEvents(func() map[event.Event][]event.Listener) ApplicationBuilder
	// WithFilters sets the application's validation filters.
	WithFilters(func() []validation.Filter) ApplicationBuilder
	// WithGrpcClientInterceptors sets the grouped gRPC client interceptors.
	WithGrpcClientInterceptors(func() map[string][]grpc.UnaryClientInterceptor) ApplicationBuilder
	// WithGrpcClientStatsHandlers sets the grouped gRPC client stats handlers.
	WithGrpcClientStatsHandlers(func() map[string][]stats.Handler) ApplicationBuilder
	// WithGrpcServerInterceptors sets the list of gRPC server interceptors.
	WithGrpcServerInterceptors(func() []grpc.UnaryServerInterceptor) ApplicationBuilder
	// WithGrpcServerStatsHandlers sets the list of gRPC server stats handlers.
	WithGrpcServerStatsHandlers(func() []stats.Handler) ApplicationBuilder
	// WithJobs registers the application's jobs.
	WithJobs(func() []queue.Job) ApplicationBuilder
	// WithMiddleware registers the http's middleware.
	WithMiddleware(func(handler configuration.Middleware)) ApplicationBuilder
	// WithMigrations registers the database migrations.
	WithMigrations(func() []schema.Migration) ApplicationBuilder
	// WithPaths sets custom paths for the application.
	WithPaths(func(paths configuration.Paths)) ApplicationBuilder
	// WithProviders registers and boots custom service providers.
	WithProviders(func() []ServiceProvider) ApplicationBuilder
	// WithRouting registers the application's routes.
	WithRouting(func()) ApplicationBuilder
	// WithRules registers the custom validation rules.
	WithRules(func() []validation.Rule) ApplicationBuilder
	// WithRunners registers the application's runners.
	WithRunners(func() []Runner) ApplicationBuilder
	// WithSchedule sets scheduled events for the application.
	WithSchedule(func() []schedule.Event) ApplicationBuilder
	// WithSeeders registers the database seeders.
	WithSeeders(func() []seeder.Seeder) ApplicationBuilder
}
