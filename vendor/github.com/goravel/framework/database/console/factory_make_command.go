package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type FactoryMakeCommand struct {
}

func NewFactoryMakeCommand() *FactoryMakeCommand {
	return &FactoryMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *FactoryMakeCommand) Signature() string {
	return "make:factory"
}

// Description The console command description.
func (r *FactoryMakeCommand) Description() string {
	return "Create a new factory class"
}

// Extend The console command extend.
func (r *FactoryMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the factory even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *FactoryMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "factory", ctx.Argument(0), support.Config.Paths.Factories)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Factory created successfully")

	return nil
}

func (r *FactoryMakeCommand) getStub() string {
	return Stubs{}.Factory()
}

// populateStub Populate the place-holders in the command stub.
func (r *FactoryMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyFactory", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
