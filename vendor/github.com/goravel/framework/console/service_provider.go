package console

import (
	"github.com/goravel/framework/console/console"
	"github.com/goravel/framework/contracts/binding"
	consolecontract "github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Artisan,
		},
		Dependencies: binding.Bindings[binding.Artisan].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Artisan, func(app foundation.Application) (any, error) {
		name := "artisan"
		usage := "Goravel Framework"
		usageText := "artisan command [options] [arguments...]"

		return NewApplication(name, usage, usageText, app.Version(), true), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	r.registerCommands(app)
}

func (r *ServiceProvider) registerCommands(app foundation.Application) {
	artisanFacade := app.MakeArtisan()
	configFacade := app.MakeConfig()
	processFacade := app.MakeProcess()
	artisanFacade.Register([]consolecontract.Command{
		console.NewListCommand(),
		console.NewKeyGenerateCommand(configFacade),
		console.NewMakeCommand(),
		console.NewBuildCommand(configFacade, processFacade),
		// The deploy command is not completely ready yet, comment it out for now.
		// https://github.com/goravel/goravel/issues/778
		// https://github.com/goravel/goravel/issues/804
		// Add the configuration when the command is enabled.
		// "deploy": map[string]any{
		// 	"base_dir":                  config.Env("DEPLOY_BASE_DIR", "/var/www/"),
		// 	"ssh_ip":                    config.Env("DEPLOY_SSH_IP", "127.0.0.1"),
		// 	"reverse_proxy_port":        config.Env("DEPLOY_REVERSE_PROXY_PORT", "9000"),
		// 	"ssh_port":                  config.Env("DEPLOY_SSH_PORT", "22"),
		// 	"ssh_user":                  config.Env("DEPLOY_SSH_USER", "root"),
		// 	"ssh_key_path":              config.Env("DEPLOY_SSH_KEY_PATH", "~/.ssh/id_rsa"),
		// 	"prod_env_file_path":        config.Env("DEPLOY_PROD_ENV_FILE_PATH", ".env.production"),
		// 	"domain":                    config.Env("DEPLOY_DOMAIN", ""),
		// 	"env_decrypt_key":           config.Env("DEPLOY_ENV_DECRYPT_KEY", ""),
		// 	"reverse_proxy_enabled":     config.Env("DEPLOY_REVERSE_PROXY_ENABLED", true),
		// 	"reverse_proxy_tls_enabled": config.Env("DEPLOY_REVERSE_PROXY_TLS_ENABLED", true),
		// 	"remote_env_decrypt":        config.Env("DEPLOY_REMOTE_ENV_DECRYPT", false),
		// },

		// "build": map[string]any{
		// 	"os":     config.Env("DEPLOY_OS", "linux"),
		// 	"arch":   config.Env("DEPLOY_ARCH", "amd64"),
		// 	"static": config.Env("DEPLOY_STATIC", true),
		// },
		// console.NewDeployCommand(configFacade, artisanFacade),
	})
}
