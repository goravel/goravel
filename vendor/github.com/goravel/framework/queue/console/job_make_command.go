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

type JobMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *JobMakeCommand) Signature() string {
	return "make:job"
}

// Description The console command description.
func (r *JobMakeCommand) Description() string {
	return "Create a new job class"
}

// Extend The console command extend.
func (r *JobMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the job even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *JobMakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "job", ctx.Argument(0), support.Config.Paths.Jobs)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(make.GetFilePath(), r.populateStub(r.getStub(), make.GetPackageName(), make.GetStructName(), make.GetSignature())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Job created successfully")

	if env.IsBootstrapSetup() {
		err = modify.AddJob(make.GetPackageImportPath(), fmt.Sprintf("&%s.%s{}", make.GetPackageName(), make.GetStructName()))
	} else {
		err = r.registerInKernel(make)
	}

	if err != nil {
		ctx.Error(errors.QueueJobRegisterFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Job registered successfully")

	return nil
}

func (r *JobMakeCommand) getStub() string {
	return JobStubs{}.Job()
}

// populateStub Populate the place-holders in the command stub.
func (r *JobMakeCommand) populateStub(stub string, packageName, structName, signature string) string {
	stub = strings.ReplaceAll(stub, "DummyJob", structName)
	stub = strings.ReplaceAll(stub, "DummySignature", str.Of(signature).Snake().String())
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}

func (r *JobMakeCommand) registerInKernel(make *supportconsole.Make) error {
	return modify.GoFile(filepath.Join("app", "providers", "queue_service_provider.go")).
		Find(match.Imports()).Modify(modify.AddImport(make.GetPackageImportPath())).
		Find(match.Jobs()).Modify(modify.Register(fmt.Sprintf("&%s.%s{}", make.GetPackageName(), make.GetStructName()))).
		Apply()
}
