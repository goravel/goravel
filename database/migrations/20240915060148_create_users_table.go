package migrations

import (
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/facades"
)

type M20240915060148CreateUsersTable struct {
}

// Signature The unique signature for the migration.
func (r *M20240915060148CreateUsersTable) Signature() string {
	return "20240915060148_create_users_table"
}

// Connection The database connection that should be used by the migration.
func (r *M20240915060148CreateUsersTable) Connection() string {
	return ""
}

// Up Run the migrations.
func (r *M20240915060148CreateUsersTable) Up() {
	facades.Schema().Create("users", func(table migration.Blueprint) {
		table.ID("id")
		table.String("migration")
		table.Integer("batch")
	})
}

// Down Reverse the migrations.
func (r *M20240915060148CreateUsersTable) Down() {
	facades.Schema().DropIfExists("users")
}
