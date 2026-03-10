package console

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	contractsseeder "github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/errors"
)

type SeedCommand struct {
	config config.Config
	seeder contractsseeder.Facade
}

func NewSeedCommand(config config.Config, seeder contractsseeder.Facade) *SeedCommand {
	return &SeedCommand{
		config: config,
		seeder: seeder,
	}
}

// Signature The name and signature of the console command.
func (r *SeedCommand) Signature() string {
	return "db:seed"
}

// Description The console command description.
func (r *SeedCommand) Description() string {
	return "Seed the database with records"
}

// Extend The console command extend.
func (r *SeedCommand) Extend() command.Extend {
	return command.Extend{
		Category: "db",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "force the operation to run when in production",
			},
			&command.StringSliceFlag{
				Name:    "seeder",
				Aliases: []string{"s"},
				Usage:   "specify the seeder(s) to run",
			},
		},
	}
}

// Handle executes the console command.
func (r *SeedCommand) Handle(ctx console.Context) error {
	force := ctx.OptionBool("force")
	if err := r.ConfirmToProceed(force); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	names := ctx.OptionSlice("seeder")
	seeders, err := r.GetSeeders(names)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}
	if len(seeders) == 0 {
		ctx.Success("no seeders found")
		return nil
	}

	if err := r.seeder.Call(seeders); err != nil {
		ctx.Error(errors.DatabaseFailToRunSeeder.Args(err).Error())
		return nil
	}
	ctx.Success("Database seeding completed successfully.")

	return nil
}

// ConfirmToProceed determines if the command should proceed based on user confirmation.
func (r *SeedCommand) ConfirmToProceed(force bool) error {
	if force || (r.config.GetString("app.env") != "production") {
		return nil
	}

	return errors.DatabaseForceIsRequiredInProduction
}

// GetSeeders returns a seeder instances
func (r *SeedCommand) GetSeeders(names []string) ([]contractsseeder.Seeder, error) {
	if len(names) == 0 {
		return r.seeder.GetSeeders(), nil
	}
	var seeders []contractsseeder.Seeder
	for _, name := range names {
		seeder := r.seeder.GetSeeder(name)
		if seeder == nil {
			return nil, errors.DatabaseSeederNotFound.Args(name)
		}
		seeders = append(seeders, seeder)
	}
	return seeders, nil
}
