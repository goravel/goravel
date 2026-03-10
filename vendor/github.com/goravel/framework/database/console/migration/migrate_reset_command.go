package migration

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
)

type MigrateResetCommand struct {
	migrator migration.Migrator
}

func NewMigrateResetCommand(migrator migration.Migrator) *MigrateResetCommand {
	return &MigrateResetCommand{
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateResetCommand) Signature() string {
	return "migrate:reset"
}

// Description The console command description.
func (r *MigrateResetCommand) Description() string {
	return "Rollback all database migrations"
}

// Extend The console command extend.
func (r *MigrateResetCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
	}
}

// Handle Execute the console command.
func (r *MigrateResetCommand) Handle(ctx console.Context) error {
	if err := r.migrator.Reset(); err != nil {
		ctx.Error(err.Error())
	}

	ctx.Success("Migration reset success")

	return nil
}
