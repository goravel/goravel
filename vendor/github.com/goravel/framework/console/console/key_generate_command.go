package console

import (
	"os"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/str"
)

type KeyGenerateCommand struct {
	config config.Config
}

func NewKeyGenerateCommand(config config.Config) *KeyGenerateCommand {
	return &KeyGenerateCommand{
		config: config,
	}
}

// Signature The name and signature of the console command.
func (r *KeyGenerateCommand) Signature() string {
	return "key:generate"
}

// Description The console command description.
func (r *KeyGenerateCommand) Description() string {
	return "Set the application key"
}

// Extend The console command extend.
func (r *KeyGenerateCommand) Extend() command.Extend {
	return command.Extend{
		Category: "key",
	}
}

// Handle Execute the console command.
func (r *KeyGenerateCommand) Handle(ctx console.Context) error {
	if r.config.GetString("app.env") == "production" {
		color.Warningln("**************************************")
		color.Warningln("*     Application In Production!     *")
		color.Warningln("**************************************")

		if !ctx.Confirm("Do you really wish to run this command?") {
			ctx.Warning("Command cancelled!")
			return nil
		}
	}

	key := r.generateRandomKey()
	if err := r.writeNewEnvironmentFileWith(key); err != nil {
		ctx.Error(err.Error())

		return nil
	}

	ctx.Success("Application key set successfully")

	return nil
}

// generateRandomKey Generate a random key for the application.
func (r *KeyGenerateCommand) generateRandomKey() string {
	return str.Random(32)
}

// writeNewEnvironmentFileWith Write a new environment file with the given key.
func (r *KeyGenerateCommand) writeNewEnvironmentFileWith(key string) error {
	content, err := os.ReadFile(support.EnvFilePath)
	if err != nil {
		return err
	}

	newContent := strings.Replace(string(content), "APP_KEY="+r.config.GetString("app.key"), "APP_KEY="+key, 1)

	err = os.WriteFile(support.EnvFilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
