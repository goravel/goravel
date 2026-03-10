package migration

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/str"
)

type MigrateMakeCommand struct {
	app      foundation.Application
	migrator migration.Migrator
}

func NewMigrateMakeCommand(app foundation.Application, migrator migration.Migrator) *MigrateMakeCommand {
	return &MigrateMakeCommand{app: app, migrator: migrator}
}

// Signature The name and signature of the console command.
func (r *MigrateMakeCommand) Signature() string {
	return "make:migration"
}

// Description The console command description.
func (r *MigrateMakeCommand) Description() string {
	return "Create a new migration file"
}

// Extend The console command extend.
func (r *MigrateMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "The model name to use for generating the migration schema",
			},
		},
	}
}

// Handle Executes the console command.
func (r *MigrateMakeCommand) Handle(ctx console.Context) error {
	makeMigration, err := supportconsole.NewMake(ctx, "migration", ctx.Argument(0), support.Config.Paths.Migrations)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}
	modelName := ctx.Option("model")

	fileName, err := r.migrator.Create(makeMigration.GetName(), modelName)
	if err != nil {
		ctx.Error(errors.MigrationCreateFailed.Args(err).Error())
		return nil
	}

	ctx.Success(fmt.Sprintf("Created Migration: %s", makeMigration.GetName()))

	structName := str.Of(fileName).Prepend("m_").Studly().String()
	if env.IsBootstrapSetup() {
		err = modify.AddMigration(makeMigration.GetPackageImportPath(), fmt.Sprintf("&%s.%s{}", makeMigration.GetPackageName(), structName))
	} else {
		err = r.registerInKernel(makeMigration.GetPackageImportPath(), structName)
	}

	if err != nil {
		ctx.Error(errors.MigrationRegisterFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Migration registered successfully")

	return nil
}

// DEPRECATED: The kernel file will be removed in future versions.
func (r *MigrateMakeCommand) registerInKernel(pkg, structName string) error {
	return modify.GoFile(r.app.DatabasePath("kernel.go")).
		Find(match.Imports()).Modify(modify.AddImport(pkg)).
		Find(match.Migrations()).Modify(modify.Register(fmt.Sprintf("&migrations.%s{}", structName))).
		Apply()
}
