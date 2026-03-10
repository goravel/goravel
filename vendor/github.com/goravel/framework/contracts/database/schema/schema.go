package schema

import (
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/database/orm"
)

type Schema interface {
	// Connection Get the connection for the schema.
	Connection(name string) Schema
	// Create a new table on the schema.
	Create(table string, callback func(table Blueprint)) error
	// Drop a table from the schema.
	Drop(table string) error
	// DropAllTables Drop all tables from the schema.
	DropAllTables() error
	// DropAllTypes Drop all types from the schema.
	DropAllTypes() error
	// DropAllViews Drop all views from the schema.
	DropAllViews() error
	// DropColumns Drop columns from a table on the schema.
	DropColumns(table string, columns []string) error
	// DropIfExists Drop a table from the schema if exists.
	DropIfExists(table string) error
	// Extend the schema with given extend parameter.
	Extend(extend Extension) Schema
	// GetColumnListing Get the column listing for a given table.
	GetColumnListing(table string) []string
	// GetColumns Get the columns for a given table.
	GetColumns(table string) ([]driver.Column, error)
	// GetConnection Get the connection of the schema.
	GetConnection() string
	// GetForeignKeys Get the foreign keys for a given table.
	GetForeignKeys(table string) ([]driver.ForeignKey, error)
	// GetIndexListing Get the names of the indexes for a given table.
	GetIndexListing(table string) []string
	// GetIndexes Get the indexes for a given table.
	GetIndexes(table string) ([]driver.Index, error)
	// GetModel Get the model from the registered models by name.
	GetModel(name string) any
	// GetTableListing Get the table listing for the database.
	GetTableListing() []string
	// GetTables Get the tables that belong to the database.
	GetTables() ([]driver.Table, error)
	// GetTypes Get the types that belong to the database.
	GetTypes() ([]driver.Type, error)
	// GetViews Get the views that belong to the database.
	GetViews() ([]driver.View, error)
	// GoTypes returns the mapping of schema types to Go types.
	GoTypes() []GoType
	// HasColumn Determine if the given table has a given column.
	HasColumn(table, column string) bool
	// HasColumns Determine if the given table has given columns.
	HasColumns(table string, columns []string) bool
	// HasIndex Determine if the given table has a given index.
	HasIndex(table, index string) bool
	// HasTable Determine if the given table exists.
	HasTable(name string) bool
	// HasType Determine if the given type exists.
	HasType(name string) bool
	// HasView Determine if the given view exists.
	HasView(name string) bool
	// Migrations Get the migrations.
	Migrations() []Migration
	// Orm Get the orm instance.
	Orm() orm.Orm
	// Prune reclaims space or optimizes underlying storage.
	Prune() error
	// Register migrations.
	Register([]Migration)
	// Rename a table on the schema.
	Rename(from, to string) error
	// SetConnection Set the connection of the schema.
	SetConnection(name string)
	// Sql Execute a sql directly.
	Sql(sql string) error
	// Table Modify a table on the schema.
	Table(table string, callback func(table Blueprint)) error
}

type Migration interface {
	// Signature Get the migration signature.
	Signature() string
	// Up Run the migrations.
	Up() error
	// Down Reverse the migrations.
	Down() error
}

type Connection interface {
	// Connection Get the connection for the migration.
	Connection() string
}

// Extension represents an extension for the schema
type Extension struct {
	GoTypes []GoType
	// Models is a list of model instances to register for runtime discovery and migration generation.
	Models []any
}

// GoType represents a database column type to Go type mapping.
type GoType struct {
	Pattern    string
	Type       string
	Import     string
	NullType   string
	NullImport string
}
