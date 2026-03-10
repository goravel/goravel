package foundation

import (
	"context"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	"github.com/goravel/framework/config"
	frameworkconsole "github.com/goravel/framework/console"
	"github.com/goravel/framework/contracts/binding"
	contractsconsole "github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/foundation/configuration"
	"github.com/goravel/framework/foundation/console"
	"github.com/goravel/framework/foundation/json"
	"github.com/goravel/framework/process"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/path"
)

var App foundation.Application
var _ = flag.String("env", support.EnvFilePath, "custom .env path")

type RunnerWithInfo struct {
	signature string
	runner    foundation.Runner
	running   atomic.Bool
	doneOnce  sync.Once
}

func init() {
	setEnv()
	setRootPath()

	app := &Application{
		Container:     NewContainer(),
		publishes:     make(map[string]map[string]string),
		publishGroups: make(map[string]map[string]string),
		runnerWg:      sync.WaitGroup{},
	}

	app.providerRepository = NewProviderRepository()
	App = app

	app.RegisterBaseServiceProviders()
	app.SetJson(json.New())
}

type Application struct {
	*Container
	ctx                context.Context
	cancel             context.CancelFunc
	builder            *ApplicationBuilder
	providerRepository foundation.ProviderRepository
	publishes          map[string]map[string]string
	publishGroups      map[string]map[string]string
	json               foundation.Json
	bootedRunners      []string
	runnerWg           sync.WaitGroup
	runnersToRun       []*RunnerWithInfo
}

func NewApplication() foundation.Application {
	return App
}

func (r *Application) About(section string, items []foundation.AboutItem) {
	console.AddAboutInformation(section, items...)
}

func (r *Application) Boot() {
	r.providerRepository.LoadFromConfig(r.MakeConfig())
	clear(r.publishes)
	clear(r.publishGroups)

	r.setTimezone()

	r.providerRepository.Register(r)
	r.providerRepository.Boot(r)

	r.registerCommands([]contractsconsole.Command{
		console.NewAboutCommand(r),
		console.NewEnvEncryptCommand(),
		console.NewEnvDecryptCommand(),
		console.NewTestMakeCommand(),
		console.NewPackageMakeCommand(),
		console.NewProviderMakeCommand(),
		console.NewPackageInstallCommand(binding.Bindings, r.MakeProcess(), r.Json()),
		console.NewPackageUninstallCommand(binding.Bindings, r.MakeProcess(), r.Json()),
		console.NewVendorPublishCommand(r.publishes, r.publishGroups),
	})
	r.bootArtisan()
}

func (r *Application) Build() foundation.Application {
	clear(r.publishes)
	clear(r.publishGroups)

	r.bootedRunners = nil
	r.runnersToRun = nil
	r.ctx, r.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	r.setTimezone()
	r.configurePaths()
	r.configureCustomConfig()
	r.configureServiceProviders()
	r.providerRepository.Register(r)
	r.providerRepository.Boot(r)
	r.configureMiddleware()
	r.configureEventListeners()
	r.configureCommands()
	r.configureSchedule()
	r.configureMigrations()
	r.configureSeeders()
	r.configureGrpc()
	r.configureJobs()
	r.configureValidation()
	r.configureRoutes()
	r.configureRunners()
	r.registerCommands([]contractsconsole.Command{
		console.NewAboutCommand(r),
		console.NewEnvEncryptCommand(),
		console.NewEnvDecryptCommand(),
		console.NewTestMakeCommand(),
		console.NewPackageMakeCommand(),
		console.NewProviderMakeCommand(),
		console.NewPackageInstallCommand(binding.Bindings, r.MakeProcess(), r.Json()),
		console.NewPackageUninstallCommand(binding.Bindings, r.MakeProcess(), r.Json()),
		console.NewVendorPublishCommand(r.publishes, r.publishGroups),
	})
	r.configureCallback()
	r.bootArtisan()

	return r
}

func (r *Application) Commands(commands []contractsconsole.Command) {
	r.registerCommands(commands)
}

func (r *Application) Context() context.Context {
	return r.ctx
}

// GetJson get the JSON implementation.
// DEPRECATED, use Json instead.
func (r *Application) GetJson() foundation.Json {
	return r.json
}

func (r *Application) IsLocale(ctx context.Context, locale string) bool {
	return r.CurrentLocale(ctx) == locale
}

func (r *Application) Json() foundation.Json {
	return r.json
}

func (r *Application) Publishes(packageName string, paths map[string]string, groups ...string) {
	if _, exist := r.publishes[packageName]; !exist {
		r.publishes[packageName] = make(map[string]string)
	}
	maps.Copy(r.publishes[packageName], paths)
	for _, group := range groups {
		r.addPublishGroup(group, paths)
	}
}

func (r *Application) Refresh() {
	// clear Container.instances except the config facade to keep the config values
	r.Fresh()

	// reset provider repository
	r.providerRepository.Reset()

	// re-register base service providers
	r.RegisterBaseServiceProviders()

	// rebuild the application
	r.Build()
}

func (r *Application) RegisterBaseServiceProviders() {
	baseProviders := r.getBaseServiceProviders()
	r.providerRepository.Add(baseProviders)
	r.providerRepository.Register(r)
}

func (r *Application) Restart() error {
	if err := r.Shutdown(); err != nil {
		return err
	}

	r.Refresh()

	go r.Start()

	i := 0
	for {
		failed := ""
		for _, runner := range r.runnersToRun {
			if !runner.running.Load() {
				failed = runner.signature
				break
			}
		}

		if failed == "" {
			break
		}
		if i > 100 {
			return fmt.Errorf("timeout waiting for runner %s to run", failed)
		}

		i++

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (r *Application) Start() {
	var (
		errsMu sync.Mutex
		errs   []error
	)

	run := func(runner *RunnerWithInfo) {
		r.runnerWg.Add(1)

		go func() {
			runner.running.Store(true)
			if err := runner.runner.Run(); err != nil {
				runner.doneOnce.Do(func() {
					r.runnerWg.Done()
				})
				runner.running.Store(false)
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("failed to run %s: %w", runner.signature, err))
				errsMu.Unlock()
				if log := r.MakeLog(); log != nil {
					log.Errorf("failed to run %s: %v\n", runner.signature, err)
				}
				r.cancel()
			}
			// Run may be a blocking call, so don't write anything after it.
		}()

		go func() {
			<-r.ctx.Done()
			if !runner.running.Load() {
				return
			}

			if err := runner.runner.Shutdown(); err != nil {
				if log := r.MakeLog(); log != nil {
					log.Errorf("failed to shutdown %s: %v\n", runner.signature, err)
				}
			}
			runner.running.Store(false)
			runner.doneOnce.Do(func() {
				r.runnerWg.Done()
			})
		}()
	}

	for _, runner := range r.runnersToRun {
		run(runner)
	}

	r.runnerWg.Wait()

	if len(errs) > 0 {
		panic(errors.Join(errs...))
	}
}

func (r *Application) SetBuilder(builder foundation.ApplicationBuilder) foundation.Application {
	r.builder = builder.(*ApplicationBuilder)

	return r
}

func (r *Application) SetJson(j foundation.Json) {
	if j != nil {
		r.json = j
	}
}

func (r *Application) SetLocale(ctx context.Context, locale string) context.Context {
	lang := r.MakeLang(ctx)
	if lang == nil {
		color.Errorln("Lang facade not initialized.")
		return ctx
	}

	return lang.SetLocale(locale)
}

func (r *Application) Shutdown() error {
	r.cancel()

	i := 0
	for {
		running := ""
		for _, runner := range r.runnersToRun {
			if runner.running.Load() {
				running = runner.signature
				break
			}
		}

		if running == "" {
			break
		}
		if i > 100 {
			return fmt.Errorf("timeout waiting for runner %s to stop", running)
		}

		i++

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (r *Application) Version() string {
	return support.Version
}

func (r *Application) BasePath(paths ...string) string {
	return path.Base(paths...)
}

func (r *Application) BootstrapPath(paths ...string) string {
	return path.Bootstrap(paths...)
}

func (r *Application) ConfigPath(paths ...string) string {
	return path.Config(paths...)
}

func (r *Application) ModelPath(paths ...string) string {
	return path.Model(paths...)
}

func (r *Application) DatabasePath(paths ...string) string {
	return path.Database(paths...)
}

func (r *Application) CurrentLocale(ctx context.Context) string {
	lang := r.MakeLang(ctx)
	if lang == nil {
		color.Errorln("Lang facade not initialized.")
		return ""
	}

	return lang.CurrentLocale()
}

func (r *Application) ExecutablePath(paths ...string) string {
	return path.Executable(paths...)
}

func (r *Application) FacadesPath(paths ...string) string {
	return path.Facade(paths...)
}

func (r *Application) LangPath(paths ...string) string {
	return path.Lang(paths...)
}

func (r *Application) Path(paths ...string) string {
	return path.App(paths...)
}

func (r *Application) PublicPath(paths ...string) string {
	return path.Public(paths...)
}

func (r *Application) ResourcePath(paths ...string) string {
	return path.Resource(paths...)
}

func (r *Application) StoragePath(paths ...string) string {
	return path.Storage(paths...)
}

func (r *Application) addPublishGroup(group string, paths map[string]string) {
	if _, exist := r.publishGroups[group]; !exist {
		r.publishGroups[group] = make(map[string]string)
	}

	maps.Copy(r.publishGroups[group], paths)
}

func (r *Application) bootArtisan() {
	artisanFacade := r.MakeArtisan()
	if artisanFacade == nil {
		color.Warningln(errors.ConsoleFacadeNotSet.Error())
		return
	}

	_ = artisanFacade.Run(os.Args, true)
}

func (r *Application) configureCallback() {
	if r.builder.callback != nil {
		r.builder.callback()
	}
}

func (r *Application) configureCommands() {
	if r.builder.commands != nil {
		if commands := r.builder.commands(); len(commands) > 0 {
			artisanFacade := r.MakeArtisan()
			if artisanFacade == nil {
				color.Errorln("Artisan facade not found, please install it first: ./artisan package:install Artisan")
			} else {
				artisanFacade.Register(commands)
			}
		}
	}
}

func (r *Application) configureCustomConfig() {
	if r.builder.config != nil {
		r.builder.config()
	}
}

func (r *Application) configureEventListeners() {
	if r.builder.eventToListeners != nil {
		if eventToListeners := r.builder.eventToListeners(); len(eventToListeners) > 0 {
			eventFacade := r.MakeEvent()
			if eventFacade == nil {
				color.Errorln("Event facade not found, please install it first: ./artisan package:install Event")
			} else {
				eventFacade.Register(eventToListeners)
			}
		}
	}
}

func (r *Application) configureGrpc() {
	var (
		grpcClientInterceptors  map[string][]grpc.UnaryClientInterceptor
		grpcServerInterceptors  []grpc.UnaryServerInterceptor
		grpcClientStatsHandlers map[string][]stats.Handler
		grpcServerStatsHandlers []stats.Handler
	)

	if r.builder.grpcClientInterceptors != nil {
		grpcClientInterceptors = r.builder.grpcClientInterceptors()
	}

	if r.builder.grpcServerInterceptors != nil {
		grpcServerInterceptors = r.builder.grpcServerInterceptors()
	}

	if r.builder.grpcClientStatsHandlers != nil {
		grpcClientStatsHandlers = r.builder.grpcClientStatsHandlers()
	}

	if r.builder.grpcServerStatsHandlers != nil {
		grpcServerStatsHandlers = r.builder.grpcServerStatsHandlers()
	}

	if len(grpcClientInterceptors) > 0 || len(grpcServerInterceptors) > 0 ||
		len(grpcClientStatsHandlers) > 0 || len(grpcServerStatsHandlers) > 0 {
		grpcFacade := r.MakeGrpc()
		if grpcFacade == nil {
			color.Errorln("gRPC facade not found, please install it first: ./artisan package:install Grpc")
		} else {
			if len(grpcClientInterceptors) > 0 {
				grpcFacade.UnaryClientInterceptorGroups(grpcClientInterceptors)
			}
			if len(grpcServerInterceptors) > 0 {
				grpcFacade.UnaryServerInterceptors(grpcServerInterceptors)
			}
			if len(grpcClientStatsHandlers) > 0 {
				grpcFacade.ClientStatsHandlerGroups(grpcClientStatsHandlers)
			}
			if len(grpcServerStatsHandlers) > 0 {
				grpcFacade.ServerStatsHandlers(grpcServerStatsHandlers)
			}
		}
	}
}

func (r *Application) configureJobs() {
	if r.builder.jobs != nil {
		jobs := r.builder.jobs()

		if len(jobs) > 0 {
			queueFacade := r.MakeQueue()
			if queueFacade == nil {
				color.Errorln("Queue facade not found, please install it first: ./artisan package:install Queue")
			} else {
				queueFacade.Register(jobs)
			}
		}
	}
}

func (r *Application) configureMiddleware() {
	if r.builder.middleware != nil {
		routeFacade := r.MakeRoute()
		if routeFacade == nil {
			color.Errorln("Route facade not found, please install it first: ./artisan package:install Route")
		} else {
			defaultGlobalMiddleware := routeFacade.GetGlobalMiddleware()
			middleware := configuration.NewMiddleware(defaultGlobalMiddleware)
			r.builder.middleware(middleware)
			routeFacade.SetGlobalMiddleware(middleware.GetGlobalMiddleware())

			if recoveryHandler := middleware.GetRecover(); recoveryHandler != nil {
				routeFacade.Recover(recoveryHandler)
			}
		}
	}
}

func (r *Application) configureMigrations() {
	if r.builder.migrations != nil {
		if migrations := r.builder.migrations(); len(migrations) > 0 {
			schemaFacade := r.MakeSchema()
			if schemaFacade == nil {
				color.Errorln("Schema facade not found, please install it first: ./artisan package:install Schema")
			} else {
				schemaFacade.Register(migrations)
			}
		}
	}
}

func (r *Application) configurePaths() {
	if r.builder.paths != nil {
		paths := configuration.NewPaths()
		r.builder.paths(paths)
	}
}

func (r *Application) configureRoutes() {
	if r.builder.routes != nil {
		r.builder.routes()
	}
}

func (r *Application) configureRunners() {
	for _, serviceProvider := range r.providerRepository.GetBooted() {
		if serviceProviderWithRunners, ok := serviceProvider.(foundation.ServiceProviderWithRunners); ok {
			for _, runner := range serviceProviderWithRunners.Runners(r) {
				signature := runner.Signature()
				if slices.Contains(r.bootedRunners, signature) {
					continue
				}

				r.bootedRunners = append(r.bootedRunners, signature)

				if runner.ShouldRun() {
					r.runnersToRun = append(r.runnersToRun, &RunnerWithInfo{signature: signature, runner: runner})
				}
			}
		}
	}

	if r.builder.runners != nil {
		for _, runner := range r.builder.runners() {
			signature := runner.Signature()
			if slices.Contains(r.bootedRunners, signature) {
				continue
			}

			r.bootedRunners = append(r.bootedRunners, signature)

			if runner.ShouldRun() {
				r.runnersToRun = append(r.runnersToRun, &RunnerWithInfo{signature: signature, runner: runner})
			}
		}
	}
}

func (r *Application) configureSchedule() {
	if r.builder.schedule != nil {
		if events := r.builder.schedule(); len(events) > 0 {
			scheduleFacade := r.MakeSchedule()
			if scheduleFacade == nil {
				color.Errorln("Schedule facade not found, please install it first: ./artisan package:install Schedule")
			} else {
				scheduleFacade.Register(events)
			}
		}
	}
}

func (r *Application) configureSeeders() {
	if r.builder.seeders != nil {
		if seeders := r.builder.seeders(); len(seeders) > 0 {
			seederFacade := r.MakeSeeder()
			if seederFacade == nil {
				color.Errorln("Seeder facade not found, please install it first: ./artisan package:install Seeder")
			} else {
				seederFacade.Register(seeders)
			}
		}
	}
}

func (r *Application) configureServiceProviders() {
	if r.builder.configuredServiceProviders != nil {
		configuredServiceProviders := r.builder.configuredServiceProviders()
		if len(configuredServiceProviders) > 0 {
			r.providerRepository.Add(configuredServiceProviders)
		}
	}
}

func (r *Application) configureValidation() {
	var (
		rules   []validation.Rule
		filters []validation.Filter
	)

	if r.builder.rules != nil {
		rules = r.builder.rules()
	}

	if r.builder.filters != nil {
		filters = r.builder.filters()
	}

	if len(rules) > 0 || len(filters) > 0 {
		validationFacade := r.MakeValidation()
		if validationFacade == nil {
			color.Errorln("Validation facade not found, please install it first: ./artisan package:install Validation")
		} else {
			if len(rules) > 0 {
				if err := validationFacade.AddRules(rules); err != nil {
					color.Errorf("add validation rules error: %+v", err)
				}
			}
			if len(filters) > 0 {
				if err := validationFacade.AddFilters(filters); err != nil {
					color.Errorf("add validation filters error: %+v", err)
				}
			}
		}
	}
}

func (r *Application) getBaseServiceProviders() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&config.ServiceProvider{},
		&frameworkconsole.ServiceProvider{},
		&process.ServiceProvider{},
	}
}

func (r *Application) registerCommands(commands []contractsconsole.Command) {
	artisanFacade := r.MakeArtisan()
	if artisanFacade == nil {
		color.Warningln(errors.ConsoleFacadeNotSet.Error())
		return
	}

	artisanFacade.Register(commands)
}

func (r *Application) setTimezone() {
	configFacade := r.MakeConfig()
	if configFacade == nil {
		color.Warningln(errors.ConfigFacadeNotSet.Error())
		carbon.SetTimezone(carbon.UTC)
		return
	}

	carbon.SetTimezone(configFacade.GetString("app.timezone", carbon.UTC))
}

func setEnv() {
	args := os.Args

	if strings.HasSuffix(args[0], ".test") ||
		strings.HasSuffix(args[0], ".test.exe") ||
		strings.Contains(args[0], "__debug") {
		support.RuntimeMode = support.RuntimeTest
		support.DontVerifyAppKey = true
	} else {
		if len(args) >= 2 {
			for _, arg := range args[1:] {
				if arg == "artisan" {
					support.RuntimeMode = support.RuntimeArtisan

					if len(args) == 2 {
						// Run go run . artisan without any command
						support.DontVerifyAppKey = true
					}
				}

				support.DontVerifyAppKey = support.DontVerifyAppKey || slices.Contains(support.DontVerifyAppKeyWhitelist, arg)
			}
		}
	}

	envFilePath := getEnvFilePath()
	if support.RuntimeMode == support.RuntimeTest {
		var (
			relativePath string
			envExist     bool
			testEnv      = envFilePath
		)

		for range 50 {
			if _, err := os.Stat(testEnv); err == nil {
				envExist = true

				break
			} else {
				testEnv = filepath.Join("..", testEnv)
				relativePath = filepath.Join("..", relativePath)
			}
		}

		if envExist {
			envFilePath = testEnv
			support.RelativePath = relativePath
		}
	}

	support.EnvFilePath = envFilePath
}

func setRootPath() {
	support.RootPath = env.CurrentAbsolutePath()
}

func getEnvFilePath() string {
	envFilePath := ".env"
	args := os.Args
	for index, arg := range args {
		if path, ok := strings.CutPrefix(arg, "--env="); ok && len(path) > 0 {
			envFilePath = path
			break
		}

		if path, ok := strings.CutPrefix(arg, "-env="); ok && len(path) > 0 {
			envFilePath = path
			break
		}

		if path, ok := strings.CutPrefix(arg, "-e="); ok && len(path) > 0 {
			envFilePath = path
			break
		}

		if arg == "--env" || arg == "-env" || arg == "-e" {
			if len(args) >= index+1 && !strings.HasPrefix(args[index+1], "-") {
				envFilePath = args[index+1]
				break
			}
		}
	}

	return envFilePath
}
