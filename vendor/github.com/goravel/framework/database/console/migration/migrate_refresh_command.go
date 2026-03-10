package migration

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/errors"
)

type MigrateRefreshCommand struct {
	artisan console.Artisan
}

func NewMigrateRefreshCommand(artisan console.Artisan) *MigrateRefreshCommand {
	return &MigrateRefreshCommand{
		artisan: artisan,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateRefreshCommand) Signature() string {
	return "migrate:refresh"
}

// Description The console command description.
func (r *MigrateRefreshCommand) Description() string {
	return "Reset and re-run all migrations"
}

// Extend The console command extend.
func (r *MigrateRefreshCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:  "step",
				Value: 0,
				Usage: "refresh steps",
			},
			&command.BoolFlag{
				Name:  "seed",
				Usage: "seed the database after running migrations",
			},
			&command.StringSliceFlag{
				Name:  "seeder",
				Usage: "specify the seeder(s) to use for seeding the database",
			},
		},
	}
}

// Handle Execute the console command.
func (r *MigrateRefreshCommand) Handle(ctx console.Context) error {
	if step := ctx.OptionInt("step"); step == 0 {
		if err := r.artisan.Call("migrate:reset"); err != nil {
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	} else {
		if err := r.artisan.Call(fmt.Sprintf("migrate:rollback --step %d", step)); err != nil {
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	}

	if err := r.artisan.Call("migrate"); err != nil {
		ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
		return nil
	}

	// Seed the database if the "seed" flag is provided
	if ctx.OptionBool("seed") {
		seeders := ctx.OptionSlice("seeder")
		seederFlag := ""
		if len(seeders) > 0 {
			seederFlag = " --seeder " + strings.Join(seeders, ",")
		}

		if err := r.artisan.Call("db:seed" + seederFlag); err != nil {
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	}
	ctx.Success("Migration refresh success")

	return nil
}
