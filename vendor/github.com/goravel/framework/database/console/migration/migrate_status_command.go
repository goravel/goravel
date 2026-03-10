package migration

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/support/color"
)

type MigrateStatusCommand struct {
	migrator migration.Migrator
}

func NewMigrateStatusCommand(migrator migration.Migrator) *MigrateStatusCommand {
	return &MigrateStatusCommand{
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateStatusCommand) Signature() string {
	return "migrate:status"
}

// Description The console command description.
func (r *MigrateStatusCommand) Description() string {
	return "Show the status of each migration"
}

// Extend The console command extend.
func (r *MigrateStatusCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
	}
}

// Handle Execute the console command.
func (r *MigrateStatusCommand) Handle(ctx console.Context) error {
	migrationStatus, err := r.migrator.Status()
	if err != nil {
		ctx.Error(err.Error())
	}
	if len(migrationStatus) > 0 {
		ctx.NewLine()
		ctx.TwoColumnDetail(color.Gray().Sprint("Migration name"), color.Gray().Sprint("Batch / Status"))
		for i := range migrationStatus {
			var status string
			if migrationStatus[i].Ran {
				status = color.Default().Sprintf("[%d] <fg=green;op=bold>Ran</>", migrationStatus[i].Batch)
			} else {
				status = color.Yellow().Sprint("<op=bold>Pending</>")
			}
			ctx.TwoColumnDetail(migrationStatus[i].Name, status)
		}
		ctx.NewLine()
	}

	return nil
}
