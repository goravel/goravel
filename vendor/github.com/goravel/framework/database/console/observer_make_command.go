package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type ObserverMakeCommand struct {
}

func NewObserverMakeCommand() *ObserverMakeCommand {
	return &ObserverMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *ObserverMakeCommand) Signature() string {
	return "make:observer"
}

// Description The console command description.
func (r *ObserverMakeCommand) Description() string {
	return "Create a new observer class"
}

// Extend The console command extend.
func (r *ObserverMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the observer even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ObserverMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "observer", ctx.Argument(0), support.Config.Paths.Observers)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		return err
	}

	ctx.Success("Observer created successfully")

	return nil
}

func (r *ObserverMakeCommand) getStub() string {
	return Stubs{}.Observer()
}

// populateStub Populate the place-holders in the command stub.
func (r *ObserverMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyObserver", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
