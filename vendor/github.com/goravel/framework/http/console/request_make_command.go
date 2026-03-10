package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type RequestMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *RequestMakeCommand) Signature() string {
	return "make:request"
}

// Description The console command description.
func (r *RequestMakeCommand) Description() string {
	return "Create a new request class"
}

// Extend The console command extend.
func (r *RequestMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the request even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *RequestMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "request", ctx.Argument(0), support.Config.Paths.Requests)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err = file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Request created successfully")

	return nil
}

func (r *RequestMakeCommand) getStub() string {
	return Stubs{}.Request()
}

// populateStub Populate the place-holders in the command stub.
func (r *RequestMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyRequest", structName)
	stub = strings.ReplaceAll(stub, "DummyField", "Name string `form:\"name\" json:\"name\"`")
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
