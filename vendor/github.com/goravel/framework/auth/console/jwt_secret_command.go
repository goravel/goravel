package console

import (
	"errors"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/str"
)

type JwtSecretCommand struct {
	config config.Config
}

func NewJwtSecretCommand(config config.Config) *JwtSecretCommand {
	return &JwtSecretCommand{config: config}
}

// Signature The name and signature of the console command.
func (r *JwtSecretCommand) Signature() string {
	return "jwt:secret"
}

// Description The console command description.
func (r *JwtSecretCommand) Description() string {
	return "Set the JWTAuth secret key used to sign the tokens"
}

// Extend The console command extend.
func (r *JwtSecretCommand) Extend() command.Extend {
	return command.Extend{
		Category: "jwt",
	}
}

// Handle Execute the console command.
func (r *JwtSecretCommand) Handle(ctx console.Context) error {
	key := r.generateRandomKey()

	if err := r.setSecretInEnvironmentFile(key); err != nil {
		ctx.Error(err.Error())

		return nil
	}

	ctx.Success("Jwt Secret set successfully")

	return nil
}

// generateRandomKey Generate a random key for the application.
func (r *JwtSecretCommand) generateRandomKey() string {
	return str.Random(32)
}

// setSecretInEnvironmentFile Set the application key in the environment file.
func (r *JwtSecretCommand) setSecretInEnvironmentFile(key string) error {
	currentKey := r.config.GetString("jwt.secret")

	if currentKey != "" {
		return errors.New("exist jwt secret")
	}

	err := r.writeNewEnvironmentFileWith(key)

	if err != nil {
		return err
	}

	return nil
}

// writeNewEnvironmentFileWith Write a new environment file with the given key.
func (r *JwtSecretCommand) writeNewEnvironmentFileWith(key string) error {
	content, err := os.ReadFile(support.EnvFilePath)
	if err != nil {
		return err
	}

	newContent := strings.Replace(string(content), "JWT_SECRET="+r.config.GetString("jwt.secret"), "JWT_SECRET="+key, 1)

	err = os.WriteFile(support.EnvFilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
