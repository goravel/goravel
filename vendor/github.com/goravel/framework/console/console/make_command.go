package console

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type MakeCommand struct {
}

func NewMakeCommand() *MakeCommand {
	return &MakeCommand{}
}

// Signature The name and signature of the console command.
func (r *MakeCommand) Signature() string {
	return "make:command"
}

// Description The console command description.
func (r *MakeCommand) Description() string {
	return "Create a new Artisan command"
}

// Extend The console command extend.
func (r *MakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:               "force",
				Aliases:            []string{"f"},
				Value:              false,
				Usage:              "Create the command even if it already exists",
				DisableDefaultText: true,
			},
		},
	}
}

// Handle Execute the console command.
func (r *MakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "command", ctx.Argument(0), support.Config.Paths.Commands)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(make.GetFilePath(), r.populateStub(r.getStub(), make.GetPackageName(), make.GetStructName(), make.GetSignature())); err != nil {
		return err
	}

	ctx.Success("Console command created successfully")

	if env.IsBootstrapSetup() {
		err = modify.AddCommand(make.GetPackageImportPath(), fmt.Sprintf("&%s.%s{}", make.GetPackageName(), make.GetStructName()))
	} else {
		err = r.registerInKernel(make)
	}

	if err != nil {
		ctx.Error(errors.ConsoleCommandRegisterFailed.Args(make.GetSignature(), err).Error())
		return nil
	}

	ctx.Success("Console command registered successfully")

	return nil
}

func (r *MakeCommand) getStub() string {
	return Stubs{}.Command()
}

// populateStub Populate the place-holders in the command stub.
func (r *MakeCommand) populateStub(stub string, packageName, structName, signature string) string {
	stub = strings.ReplaceAll(stub, "DummyCommand", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)
	stub = strings.ReplaceAll(stub, "DummySignature", str.Of(signature).Kebab().Prepend("app:").String())

	return stub
}

func (r *MakeCommand) registerInKernel(make *supportconsole.Make) error {
	if err := modify.GoFile(filepath.Join("app", "console", "kernel.go")).
		Find(match.Imports()).Modify(modify.AddImport(make.GetPackageImportPath())).
		Find(match.Commands()).Modify(modify.Register(fmt.Sprintf("&%s.%s{}", make.GetPackageName(), make.GetStructName()))).
		Apply(); err != nil {

		return err
	}

	return nil
}
