package console

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
)

type ProviderMakeCommand struct {
}

func NewProviderMakeCommand() *ProviderMakeCommand {
	return &ProviderMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *ProviderMakeCommand) Signature() string {
	return "make:provider"
}

// Description The console command description.
func (r *ProviderMakeCommand) Description() string {
	return "Create a new service provider class"
}

// Extend The console command extend.
func (r *ProviderMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the provider even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ProviderMakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "provider", ctx.Argument(0), support.Config.Paths.Providers)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	stub := r.getStub()

	if err := file.PutContent(make.GetFilePath(), r.populateStub(stub, make.GetPackageName(), make.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Provider created successfully")

	if env.IsBootstrapSetup() {
		if err := modify.AddProvider(make.GetPackageImportPath(), fmt.Sprintf("&%s.%s{}", make.GetPackageName(), make.GetStructName())); err != nil {
			ctx.Error(errors.ProviderRegisterFailed.Args(make.GetStructName(), err).Error())
			return nil
		}

		ctx.Success("Provider registered successfully")
	}

	return nil
}

func (r *ProviderMakeCommand) getStub() string {
	return Stubs{}.ServiceProvider()
}

// populateStub Populate the place-holders in the command stub.
func (r *ProviderMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyServiceProvider", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
