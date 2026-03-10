package console

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/errors"
	supportconsole "github.com/goravel/framework/support/console"
)

type WipeCommand struct {
	config config.Config
	schema schema.Schema
}

func NewWipeCommand(config config.Config, schema schema.Schema) *WipeCommand {
	return &WipeCommand{
		config: config,
		schema: schema,
	}
}

// Signature The name and signature of the console command.
func (r *WipeCommand) Signature() string {
	return "db:wipe"
}

// Description The console command description.
func (r *WipeCommand) Description() string {
	return "Drop all tables, views, and types"
}

// Extend The console command extend.
func (r *WipeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "db",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "database",
				Aliases: []string{"d"},
				Usage:   "The database connection to use",
			},
			&command.BoolFlag{
				Name:    "drop-views",
				Aliases: []string{"dv"},
				Usage:   "Drop all tables and views",
			},
			&command.BoolFlag{
				Name:    "drop-types",
				Aliases: []string{"dt"},
				Usage:   "Drop all tables and types (Postgres only)",
			},
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force the operation to run when in production",
			},
		},
	}
}

// Handle Execute the console command.
func (r *WipeCommand) Handle(ctx console.Context) error {
	if !supportconsole.ConfirmToProceed(ctx, r.config.GetString("app.env")) {
		ctx.Warning(errors.ConsoleRunInProduction.Error())
		return nil
	}

	database := ctx.Option("database")

	if ctx.OptionBool("drop-views") {
		if err := r.dropAllViews(database); err != nil {
			ctx.Error(errors.ConsoleDropAllViewsFailed.Args(err).Error())
			return nil
		}

		ctx.Success("Dropped all views successfully")
	}

	if err := r.dropAllTables(database); err != nil {
		ctx.Error(errors.ConsoleDropAllTablesFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Dropped all tables successfully")

	if ctx.OptionBool("drop-types") {
		if err := r.dropAllTypes(database); err != nil {
			ctx.Error(errors.ConsoleDropAllTypesFailed.Args(err).Error())
			return nil
		}

		ctx.Success("Dropped all types successfully")
	}

	if err := r.schema.Connection(database).Prune(); err != nil {
		ctx.Error(errors.ConsolePruneFailed.Args(err).Error())
		return nil
	}

	return nil
}

func (r *WipeCommand) dropAllTables(database string) error {
	return r.schema.Connection(database).DropAllTables()
}

func (r *WipeCommand) dropAllViews(database string) error {
	return r.schema.Connection(database).DropAllViews()
}

func (r *WipeCommand) dropAllTypes(database string) error {
	return r.schema.Connection(database).DropAllTypes()
}
