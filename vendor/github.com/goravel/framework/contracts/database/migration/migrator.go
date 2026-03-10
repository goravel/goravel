package migration

const (
	MigratorDefault = "default"
	MigratorSql     = "sql"
)

type Status struct {
	Name  string
	Batch int
	Ran   bool
}

type Migrator interface {
	// Create a new migration file.
	Create(name string, modelName string) (string, error)
	// Fresh drops all tables and re-runs all migrations from scratch.
	Fresh() error
	// Reset the migrations.
	Reset() error
	// Rollback the last migration operation.
	Rollback(step, batch int) error
	// Run the migrations according to paths.
	Run() error
	// Status gets the migration's status.
	Status() ([]Status, error)
}
