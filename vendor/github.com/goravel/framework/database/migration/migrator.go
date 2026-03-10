package migration

import (
	"fmt"
	"slices"

	"github.com/goravel/framework/contracts/console"
	contractsmigration "github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/contracts/database/orm"
	contractsschema "github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/support/collect"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
	supportfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type Migrator struct {
	artisan    console.Artisan
	creator    *Creator
	repository contractsmigration.Repository
	schema     contractsschema.Schema
}

func NewMigrator(artisan console.Artisan, schema contractsschema.Schema, table string) *Migrator {
	return &Migrator{
		artisan:    artisan,
		creator:    NewCreator(),
		repository: NewRepository(schema, table),
		schema:     schema,
	}
}

func (r *Migrator) Create(name string, modelName string) (string, error) {
	table, create := TableGuesser{}.Guess(name)

	var schemaFields []string
	if modelName != "" {
		model := r.schema.GetModel(modelName)
		if model == nil {
			return "", errors.SchemaModelNotFound.Args(modelName)
		}

		var err error
		table, schemaFields, err = Generate(model)
		if err != nil {
			return "", err
		}
	}

	stub := r.creator.GetStub(table, create)

	// Prepend timestamp to the file name.
	fileName := r.creator.GetFileName(name)
	facadesImport := packages.Paths().Facades().Import()
	if !env.IsBootstrapSetup() {
		facadesImport = "github.com/goravel/framework/facades"
	}

	templateData := StubData{
		FacadesPackage: packages.Paths().Facades().Package(),
		FacadesImport:  facadesImport,
		Package:        packages.Paths().Migrations().Package(),
		SchemaFields:   schemaFields,
		Signature:      fileName,
		StructName:     str.Of(fileName).Prepend("m_").Studly().String(),
		Table:          table,
	}

	content, err := r.creator.PopulateStub(stub, templateData)
	if err != nil {
		return "", err
	}

	if err := supportfile.PutContent(r.creator.GetPath(fileName), content); err != nil {
		return "", err
	}

	return fileName, nil
}

func (r *Migrator) Fresh() error {
	if err := r.artisan.Call("db:wipe --force"); err != nil {
		return err
	}
	if err := r.artisan.Call("migrate"); err != nil {
		return err
	}

	return nil
}

func (r *Migrator) Reset() error {
	if !r.repository.RepositoryExists() {
		color.Warningln("Migration table not found")

		return nil
	}

	ran, err := r.repository.GetRan()
	if err != nil {
		return err
	}

	return r.Rollback(len(ran), 0)
}

func (r *Migrator) Rollback(step, batch int) error {
	if !r.repository.RepositoryExists() {
		color.Warningln("Migration table not found")

		return nil
	}

	files, err := r.getFilesForRollback(step, batch)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		color.Infoln("Nothing to rollback")

		return nil
	}

	color.Infoln("Rolling back migration")

	for _, file := range files {
		migration := r.getMigrationViaFile(file)
		if migration == nil {
			color.Warningf("Migration not found: %s\n", file.Migration)

			continue
		}

		if err := r.runDown(migration); err != nil {
			return err
		}

		color.Infoln("Rolled back:", migration.Signature())
	}

	return nil
}

func (r *Migrator) Run() error {
	if err := r.prepareDatabase(); err != nil {
		return err
	}

	ran, err := r.repository.GetRan()
	if err != nil {
		return err
	}

	pendingMigrations := r.pendingMigrations(ran)

	return r.runPending(pendingMigrations)
}

func (r *Migrator) Status() ([]contractsmigration.Status, error) {
	if !r.repository.RepositoryExists() {
		color.Warningln("Migration table not found")

		return nil, nil
	}

	batches, err := r.repository.GetMigrations()
	if err != nil {
		return nil, err
	}

	migrationStatus := r.getStatusForMigrations(batches)
	if len(migrationStatus) == 0 {
		color.Warningln("No migrations found")

		return nil, nil
	}

	return migrationStatus, nil
}

func (r *Migrator) getFilesForRollback(step, batch int) ([]contractsmigration.File, error) {
	if step > 0 {
		return r.repository.GetMigrationsByStep(step)
	}

	if batch > 0 {
		return r.repository.GetMigrationsByBatch(batch)
	}

	return r.repository.GetLast()
}

func (r *Migrator) getMigrationViaFile(file contractsmigration.File) contractsschema.Migration {
	for _, migration := range r.schema.Migrations() {
		if migration.Signature() == file.Migration {
			return migration
		}
	}

	return nil
}

func (r *Migrator) getStatusForMigrations(batches []contractsmigration.File) []contractsmigration.Status {
	var migrationStatus []contractsmigration.Status

	for _, migration := range r.schema.Migrations() {
		var file contractsmigration.File
		collect.Each(batches, func(item contractsmigration.File, index int) {
			if item.Migration == migration.Signature() {
				file = item
				return
			}
		})

		if file.ID > 0 {
			migrationStatus = append(migrationStatus, contractsmigration.Status{
				Name:  migration.Signature(),
				Batch: file.Batch,
				Ran:   true,
			})
		} else {
			migrationStatus = append(migrationStatus, contractsmigration.Status{
				Name: migration.Signature(),
				Ran:  false,
			})
		}
	}

	return migrationStatus
}

func (r *Migrator) pendingMigrations(ran []string) []contractsschema.Migration {
	var pendingMigrations []contractsschema.Migration
	for _, migration := range r.schema.Migrations() {
		if !slices.Contains(ran, migration.Signature()) {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	return pendingMigrations
}

func (r *Migrator) prepareDatabase() error {
	if r.repository.RepositoryExists() {
		return nil
	}

	return r.repository.CreateRepository()
}

func (r *Migrator) printTitle(maxNameLength int) {
	color.Default().Print(fmt.Sprintf("%-*s", maxNameLength, "Migration name"))
	color.Default().Println(" | Batch / Status")
	for i := 0; i < maxNameLength+17; i++ {
		color.Default().Print("-")
	}
	color.Default().Println()
}

func (r *Migrator) runPending(migrations []contractsschema.Migration) error {
	if len(migrations) == 0 {
		color.Infoln("Nothing to migrate")

		return nil
	}

	batch, err := r.repository.GetNextBatchNumber()
	if err != nil {
		return err
	}

	color.Infoln("Running migration")

	for _, migration := range migrations {
		color.Infoln("Running:", migration.Signature())

		if err := r.runUp(migration, batch); err != nil {
			return err
		}
	}

	return nil
}

func (r *Migrator) runDown(migration contractsschema.Migration) error {
	defaultConnection := r.schema.GetConnection()
	defaultQuery := r.schema.Orm().Query()
	if connectionMigration, ok := migration.(contractsschema.Connection); ok {
		r.schema.SetConnection(connectionMigration.Connection())
	}

	defer func() {
		// reset the connection and query to default, to avoid err and panic
		r.schema.Orm().SetQuery(defaultQuery)
		r.schema.SetConnection(defaultConnection)
	}()

	if err := r.schema.Orm().Transaction(func(tx orm.Query) error {
		r.schema.Orm().SetQuery(tx)

		return migration.Down()
	}); err != nil {
		return err
	}

	// repository.Log should be called in the default connection.
	// The code below can't be set in the transaction, because the connection will conflict.
	r.schema.Orm().SetQuery(defaultQuery)
	r.schema.SetConnection(defaultConnection)

	return r.repository.Delete(migration.Signature())
}

func (r *Migrator) runUp(migration contractsschema.Migration, batch int) error {
	defaultConnection := r.schema.GetConnection()
	defaultQuery := r.schema.Orm().Query()
	if connectionMigration, ok := migration.(contractsschema.Connection); ok {
		r.schema.SetConnection(connectionMigration.Connection())
	}

	defer func() {
		// reset the connection and query to default, to avoid err and panic
		r.schema.Orm().SetQuery(defaultQuery)
		r.schema.SetConnection(defaultConnection)
	}()

	if err := r.schema.Orm().Transaction(func(tx orm.Query) error {
		r.schema.Orm().SetQuery(tx)

		return migration.Up()
	}); err != nil {
		return err
	}

	// repository.Log should be called in the default connection.
	// The code below can't be set in the transaction, because the connection will conflict.
	r.schema.Orm().SetQuery(defaultQuery)
	r.schema.SetConnection(defaultConnection)

	return r.repository.Log(migration.Signature(), batch)
}
