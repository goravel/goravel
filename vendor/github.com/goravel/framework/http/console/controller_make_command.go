package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type ControllerMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *ControllerMakeCommand) Signature() string {
	return "make:controller"
}

// Description The console command description.
func (r *ControllerMakeCommand) Description() string {
	return "Create a new controller class"
}

// Extend The console command extend.
func (r *ControllerMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:  "resource",
				Value: false,
				Usage: "resourceful controller",
			},
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the controller even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ControllerMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "controller", ctx.Argument(0), support.Config.Paths.Controllers)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	stub := r.getStub()
	if ctx.OptionBool("resource") {
		stub = r.getResourceStub()
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(stub, m.GetPackageName(), m.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Controller created successfully")

	return nil
}

func (r *ControllerMakeCommand) getStub() string {
	return Stubs{}.Controller()
}

func (r *ControllerMakeCommand) getResourceStub() string {
	return Stubs{}.ResourceController()
}

// populateStub Populate the place-holders in the command stub.
func (r *ControllerMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyController", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
