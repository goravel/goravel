package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20210101000002CreateJobsTable struct{}

// Signature The unique signature for the migration.
func (r *M20210101000002CreateJobsTable) Signature() string {
	return "20210101000002_create_jobs_table"
}

// Up Run the migrations.
func (r *M20210101000002CreateJobsTable) Up() error {
	if !facades.Schema().HasTable("jobs") {
		if err := facades.Schema().Create("jobs", func(table schema.Blueprint) {
			table.ID()
			table.String("queue")
			table.LongText("payload")
			table.UnsignedTinyInteger("attempts").Default(0)
			table.DateTimeTz("reserved_at").Nullable()
			table.DateTimeTz("available_at")
			table.DateTimeTz("created_at").UseCurrent()
			table.Index("queue")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("failed_jobs") {
		if err := facades.Schema().Create("failed_jobs", func(table schema.Blueprint) {
			table.ID()
			table.String("uuid")
			table.Text("connection")
			table.Text("queue")
			table.LongText("payload")
			table.LongText("exception")
			table.DateTimeTz("failed_at").UseCurrent()
			table.Unique("uuid")
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20210101000002CreateJobsTable) Down() error {
	if err := facades.Schema().DropIfExists("jobs"); err != nil {
		return err
	}

	if err := facades.Schema().DropIfExists("failed_jobs"); err != nil {
		return err
	}

	return nil
}
