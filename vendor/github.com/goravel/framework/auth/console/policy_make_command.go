package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type PolicyMakeCommand struct {
}

func NewPolicyMakeCommand() *PolicyMakeCommand {
	return &PolicyMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *PolicyMakeCommand) Signature() string {
	return "make:policy"
}

// Description The console command description.
func (r *PolicyMakeCommand) Description() string {
	return "Create a new policy class"
}

// Extend The console command extend.
func (r *PolicyMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the policy even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *PolicyMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "policy", ctx.Argument(0), support.Config.Paths.Policies)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		return err
	}

	ctx.Success("Policy created successfully")

	return nil
}

func (r *PolicyMakeCommand) getStub() string {
	return PolicyStubs{}.Policy()
}

// populateStub Populate the place-holders in the command stub.
func (r *PolicyMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyPolicy", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
