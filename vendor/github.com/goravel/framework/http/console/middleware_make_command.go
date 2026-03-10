package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type MiddlewareMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *MiddlewareMakeCommand) Signature() string {
	return "make:middleware"
}

// Description The console command description.
func (r *MiddlewareMakeCommand) Description() string {
	return "Create a new middleware class"
}

// Extend The console command extend.
func (r *MiddlewareMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the middleware even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *MiddlewareMakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "middleware", ctx.Argument(0), support.Config.Paths.Middleware)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(make.GetFilePath(), r.populateStub(r.getStub(), make.GetPackageName(), make.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Middleware created successfully")

	return nil
}

func (r *MiddlewareMakeCommand) getStub() string {
	return Stubs{}.Middleware()
}

// populateStub Populate the place-holders in the command stub.
func (r *MiddlewareMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyMiddleware", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
