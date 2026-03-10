package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type ListenerMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *ListenerMakeCommand) Signature() string {
	return "make:listener"
}

// Description The console command description.
func (r *ListenerMakeCommand) Description() string {
	return "Create a new listener class"
}

// Extend The console command extend.
func (r *ListenerMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the listener even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ListenerMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "listener", ctx.Argument(0), support.Config.Paths.Listeners)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Listener created successfully")

	return nil
}

func (r *ListenerMakeCommand) getStub() string {
	return Stubs{}.Listener()
}

// populateStub Populate the place-holders in the command stub.
func (r *ListenerMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyListener", structName)
	stub = strings.ReplaceAll(stub, "DummyName", str.Of(structName).Snake().String())
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
