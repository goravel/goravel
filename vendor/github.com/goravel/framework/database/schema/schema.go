package schema

import (
	"reflect"
	"slices"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database/driver"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsschema "github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
)

var _ contractsschema.Schema = (*Schema)(nil)

type Schema struct {
	config           config.Config
	driver           driver.Driver
	grammar          driver.Grammar
	log              log.Log
	migrations       []contractsschema.Migration
	orm              contractsorm.Orm
	prefix           string
	processor        driver.Processor
	schema           string
	goTypes          []contractsschema.GoType
	models           []any
	modelsByFullName map[string]any
}

func NewSchema(config config.Config, log log.Log, orm contractsorm.Orm, driver driver.Driver, migrations []contractsschema.Migration) (*Schema, error) {
	writers := driver.Pool().Writers
	if len(writers) == 0 {
		return nil, errors.DatabaseConfigNotFound
	}

	prefix := writers[0].Prefix
	schema := writers[0].Schema
	grammar := driver.Grammar()
	processor := driver.Processor()

	return &Schema{
		config:           config,
		driver:           driver,
		grammar:          grammar,
		log:              log,
		migrations:       migrations,
		orm:              orm,
		prefix:           prefix,
		processor:        processor,
		schema:           schema,
		goTypes:          defaultGoTypes(),
		models:           make([]any, 0),
		modelsByFullName: make(map[string]any),
	}, nil
}

func (r *Schema) Connection(name string) contractsschema.Schema {
	schema, err := NewSchema(r.config, r.log, r.orm.Connection(name), r.driver, r.migrations)
	if err != nil {
		r.log.Panic(errors.SchemaConnectionNotFound.Args(name).SetModule(errors.ModuleSchedule).Error())
		return nil
	}

	return schema
}

func (r *Schema) Create(table string, callback func(table contractsschema.Blueprint)) error {
	blueprint := r.createBlueprint(table)
	blueprint.Create()
	callback(blueprint)

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToCreateTable.Args(table, err)
	}

	return nil
}

func (r *Schema) Drop(table string) error {
	blueprint := r.createBlueprint(table)
	blueprint.Drop()

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToDropTable.Args(table, err)
	}

	return nil
}

func (r *Schema) DropAllTables() error {
	tables, err := r.GetTables()
	if err != nil {
		return err
	}

	sqls := r.grammar.CompileDropAllTables(r.schema, tables)
	if sqls == nil {
		return nil
	}

	return r.orm.Transaction(func(tx contractsorm.Query) error {
		for _, sql := range sqls {
			if _, err := tx.Exec(sql); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Schema) DropAllTypes() error {
	types, err := r.GetTypes()
	if err != nil {
		return err
	}

	return r.orm.Transaction(func(tx contractsorm.Query) error {
		for _, sql := range r.grammar.CompileDropAllTypes(r.schema, types) {
			if _, err := tx.Exec(sql); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Schema) DropAllViews() error {
	views, err := r.GetViews()
	if err != nil {
		return err
	}

	sqls := r.grammar.CompileDropAllViews(r.schema, views)
	if sqls == nil {
		return nil
	}

	return r.orm.Transaction(func(tx contractsorm.Query) error {
		for _, sql := range sqls {
			if _, err := tx.Exec(sql); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Schema) DropColumns(table string, columns []string) error {
	blueprint := r.createBlueprint(table)
	blueprint.DropColumn(columns...)

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToDropColumns.Args(table, err)
	}

	return nil
}

func (r *Schema) DropIfExists(table string) error {
	blueprint := r.createBlueprint(table)
	blueprint.DropIfExists()

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToDropTable.Args(table, err)
	}

	return nil
}

func (r *Schema) Extend(extend contractsschema.Extension) contractsschema.Schema {
	r.extendGoTypes(extend.GoTypes)
	r.extendModels(extend.Models)
	return r
}

func (r *Schema) GetColumnListing(table string) []string {
	columns, err := r.GetColumns(table)
	if err != nil {
		r.log.Errorf("failed to get %s columns: %v", table, err)
		return nil
	}

	var names []string
	for _, column := range columns {
		names = append(names, column.Name)
	}

	return names
}

func (r *Schema) GetColumns(table string) ([]driver.Column, error) {
	var dbColumns []driver.DBColumn
	sql, err := r.grammar.CompileColumns(r.schema, table)
	if err != nil {
		return nil, err
	}

	if err := r.orm.Query().Raw(sql).Scan(&dbColumns); err != nil {
		return nil, err
	}

	return r.processor.ProcessColumns(dbColumns), nil
}

func (r *Schema) GetConnection() string {
	return r.orm.Name()
}

func (r *Schema) GetForeignKeys(table string) ([]driver.ForeignKey, error) {
	table = r.prefix + table

	var dbForeignKeys []driver.DBForeignKey
	if err := r.orm.Query().Raw(r.grammar.CompileForeignKeys(r.schema, table)).Scan(&dbForeignKeys); err != nil {
		return nil, err
	}

	return r.processor.ProcessForeignKeys(dbForeignKeys), nil
}

func (r *Schema) GetIndexListing(table string) []string {
	indexes, err := r.GetIndexes(table)
	if err != nil {
		r.log.Errorf("failed to get %s indexes: %v", table, err)
		return nil
	}

	var names []string
	for _, index := range indexes {
		names = append(names, index.Name)
	}

	return names
}

func (r *Schema) GetIndexes(table string) ([]driver.Index, error) {
	var dbIndexes []driver.DBIndex
	sql, err := r.grammar.CompileIndexes(r.schema, table)
	if err != nil {
		return nil, err
	}

	if err := r.orm.Query().Raw(sql).Scan(&dbIndexes); err != nil {
		return nil, err
	}

	return r.processor.ProcessIndexes(dbIndexes), nil
}

// GetModel retrieves a registered model by name. If the name doesn't contain a package
// (no dot), it automatically prepends "models." to the name. Returns nil if the model
// is not found in the registry.
func (r *Schema) GetModel(name string) any {
	// If no dot, assume "models" package
	if !strings.Contains(name, ".") {
		name = "models." + name
	}

	return r.modelsByFullName[name]
}

func (r *Schema) GetTableListing() []string {
	tables, err := r.GetTables()
	if err != nil {
		r.log.Errorf("failed to get tables: %v", err)
		return nil
	}

	var names []string
	for _, table := range tables {
		names = append(names, table.Name)
	}

	return names
}

func (r *Schema) GetTables() ([]driver.Table, error) {
	var tables []driver.Table
	if err := r.orm.Query().Raw(r.grammar.CompileTables(r.orm.DatabaseName())).Scan(&tables); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *Schema) GetTypes() ([]driver.Type, error) {
	var types []driver.Type
	if err := r.orm.Query().Raw(r.grammar.CompileTypes()).Scan(&types); err != nil {
		return nil, err
	}

	return r.processor.ProcessTypes(types), nil
}

func (r *Schema) GetViews() ([]driver.View, error) {
	var views []driver.View
	if err := r.orm.Query().Raw(r.grammar.CompileViews(r.orm.DatabaseName())).Scan(&views); err != nil {
		return nil, err
	}

	return views, nil
}

func (r *Schema) GoTypes() []contractsschema.GoType {
	return r.goTypes
}

func (r *Schema) HasColumn(table, column string) bool {
	return slices.Contains(r.GetColumnListing(table), column)
}

func (r *Schema) HasColumns(table string, columns []string) bool {
	columnListing := r.GetColumnListing(table)
	for _, column := range columns {
		if !slices.Contains(columnListing, column) {
			return false
		}
	}

	return true
}

func (r *Schema) HasIndex(table, index string) bool {
	indexListing := r.GetIndexListing(table)

	return slices.Contains(indexListing, index)
}

func (r *Schema) HasTable(name string) bool {
	var schema string
	if strings.Contains(name, ".") {
		lastDotIndex := strings.LastIndex(name, ".")
		schema = name[:lastDotIndex]
		name = name[lastDotIndex+1:]
	}

	tableName := r.prefix + name

	tables, err := r.GetTables()
	if err != nil {
		r.log.Errorf(errors.SchemaFailedToGetTables.Args(r.orm.Name(), err).Error())
		return false
	}

	for _, table := range tables {
		if table.Name == tableName {
			if schema == "" || schema == table.Schema {
				return true
			}
		}
	}

	return false
}

func (r *Schema) HasType(name string) bool {
	types, err := r.GetTypes()
	if err != nil {
		r.log.Errorf(errors.SchemaFailedToGetTables.Args(r.orm.Name(), err).Error())
		return false
	}

	for _, t := range types {
		if t.Name == name {
			return true
		}
	}

	return false
}

func (r *Schema) HasView(name string) bool {
	views, err := r.GetViews()
	if err != nil {
		r.log.Errorf(errors.SchemaFailedToGetTables.Args(r.orm.Name(), err).Error())
		return false
	}

	for _, view := range views {
		if view.Name == name {
			return true
		}
	}

	return false
}

func (r *Schema) Migrations() []contractsschema.Migration {
	return r.migrations
}

func (r *Schema) Orm() contractsorm.Orm {
	return r.orm
}

func (r *Schema) Prune() error {
	if sql := r.grammar.CompilePrune(r.orm.DatabaseName()); len(sql) > 0 {
		_, err := r.orm.Query().Exec(sql)

		return err
	}

	return nil
}

func (r *Schema) Register(migrations []contractsschema.Migration) {
	existingSignatures := make(map[string]bool)

	for _, migration := range migrations {
		signature := migration.Signature()

		if existingSignatures[signature] {
			color.Errorf("Duplicate migration signature: %s in %T\n", signature, migration)
		} else {
			existingSignatures[signature] = true
			r.migrations = append(r.migrations, migration)
		}
	}
}

func (r *Schema) Rename(from, to string) error {
	blueprint := r.createBlueprint(from)
	blueprint.Rename(to)

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToRenameTable.Args(from, err)
	}

	return nil
}

func (r *Schema) SetConnection(name string) {
	r.orm = r.orm.Connection(name)
}

func (r *Schema) Sql(sql string) error {
	_, err := r.orm.Query().Exec(sql)

	return err
}

func (r *Schema) Table(table string, callback func(table contractsschema.Blueprint)) error {
	blueprint := r.createBlueprint(table)
	callback(blueprint)

	if err := r.build(blueprint); err != nil {
		return errors.SchemaFailedToChangeTable.Args(table, err)
	}

	return nil
}

func (r *Schema) build(blueprint contractsschema.Blueprint) error {
	if r.orm.Query().InTransaction() {
		return blueprint.Build(r.orm.Query(), r.grammar)
	}

	return r.orm.Transaction(func(tx contractsorm.Query) error {
		return blueprint.Build(tx, r.grammar)
	})
}

func (r *Schema) createBlueprint(table string) contractsschema.Blueprint {
	return NewBlueprint(r, r.prefix, table)
}

// extendGoTypes merges user-provided GoType overrides and additions into the schema's default mappings.
// New patterns (not present in defaults) are prepended for highest priority. Existing patterns are updated with non-zero override fields.
func (r *Schema) extendGoTypes(overrides []contractsschema.GoType) {
	if len(overrides) == 0 {
		return
	}

	defaults := r.goTypes
	defaultPatterns := make(map[string]bool, len(defaults))
	for _, d := range defaults {
		defaultPatterns[d.Pattern] = true
	}

	overrideMap := make(map[string]contractsschema.GoType, len(overrides))
	var newPatterns []contractsschema.GoType
	for _, o := range overrides {
		overrideMap[o.Pattern] = o
		if !defaultPatterns[o.Pattern] {
			newPatterns = append(newPatterns, o)
		}
	}

	result := make([]contractsschema.GoType, 0, len(defaults)+len(newPatterns))
	result = append(result, newPatterns...)

	for _, d := range defaults {
		if o, exists := overrideMap[d.Pattern]; exists {
			if o.Type != "" {
				d.Type = o.Type
			}
			if o.NullType != "" {
				d.NullType = o.NullType
			}
			if o.Import != "" {
				d.Import = o.Import
			}
			if o.NullImport != "" {
				d.NullImport = o.NullImport
			}
		}
		result = append(result, d)
	}

	r.goTypes = result
}

func (r *Schema) extendModels(models []any) {
	for _, m := range models {
		fullName := getModelName(m)
		if fullName == "" {
			continue
		}

		// Use full name for duplicate detection
		if _, exists := r.modelsByFullName[fullName]; exists {
			continue
		}

		r.models = append(r.models, m)
		r.modelsByFullName[fullName] = m
	}
}

// modelType returns the reflect.Type for a model, dereferencing pointers.
func modelType(m any) reflect.Type {
	if m == nil {
		return nil
	}
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// getModelName returns package.TypeName (e.g., "models.User").
func getModelName(m any) string {
	t := modelType(m)
	if t == nil {
		return ""
	}
	if pkg := t.PkgPath(); pkg != "" {
		if i := strings.LastIndexByte(pkg, '/'); i >= 0 {
			pkg = pkg[i+1:]
		}
		return pkg + "." + t.Name()
	}
	return t.Name()
}

func defaultGoTypes() []contractsschema.GoType {
	return []contractsschema.GoType{
		// Special cases first - these need to be matched before general patterns
		{Pattern: "(?i)^tinyint\\(1\\)$", Type: "bool", NullType: "*bool"}, // MySQL boolean representation

		// Boolean types
		{Pattern: "(?i)^bool(ean)?$", Type: "bool", NullType: "*bool"},
		{Pattern: "(?i)^bit(\\(1\\))?$", Type: "bool", NullType: "*bool"}, // Single bit as boolean

		// Integer types - ordered from most specific to general
		// Unsigned variants (MySQL)
		{Pattern: "(?i)^bigint\\s+unsigned(\\(\\d+\\))?$", Type: "uint64", NullType: "*uint64"},
		{Pattern: "(?i)^int(eger)?\\s+unsigned(\\(\\d+\\))?$", Type: "uint32", NullType: "*uint32"},
		{Pattern: "(?i)^mediumint\\s+unsigned(\\(\\d+\\))?$", Type: "uint32", NullType: "*uint32"},
		{Pattern: "(?i)^smallint\\s+unsigned(\\(\\d+\\))?$", Type: "uint16", NullType: "*uint16"},
		{Pattern: "(?i)^tinyint\\s+unsigned(\\(\\d+\\))?$", Type: "uint8", NullType: "*uint8"},

		// PostgreSQL serials
		{Pattern: "(?i)^bigserial(\\(\\d+\\))?$", Type: "int64", NullType: "*int64"},
		{Pattern: "(?i)^serial(\\(\\d+\\))?$", Type: "int", NullType: "*int"},
		{Pattern: "(?i)^smallserial(\\(\\d+\\))?$", Type: "int16", NullType: "*int16"},

		// Standard integer types
		{Pattern: "(?i)^bigint(\\(\\d+\\))?$", Type: "int64", NullType: "*int64"},
		{Pattern: "(?i)^int8(\\(\\d+\\))?$", Type: "int64", NullType: "*int64"},      // PostgreSQL
		{Pattern: "(?i)^mediumint(\\(\\d+\\))?$", Type: "int32", NullType: "*int32"}, // MySQL
		{Pattern: "(?i)^int(eger)?(\\(\\d+\\))?$", Type: "int", NullType: "*int"},
		{Pattern: "(?i)^int4(\\(\\d+\\))?$", Type: "int32", NullType: "*int32"}, // PostgreSQL
		{Pattern: "(?i)^smallint(\\(\\d+\\))?$", Type: "int16", NullType: "*int16"},
		{Pattern: "(?i)^int2(\\(\\d+\\))?$", Type: "int16", NullType: "*int16"},  // PostgreSQL
		{Pattern: "(?i)^tinyint(\\(\\d+\\))?$", Type: "int8", NullType: "*int8"}, // MySQL (when not tinyint(1))
		{Pattern: "(?i)^year(\\(\\d+\\))?$", Type: "int16", NullType: "*int16"},  // MySQL YEAR type

		// Bit/binary integer types
		{Pattern: "(?i)^bit(\\(\\d+\\))?$", Type: "[]byte", NullType: "[]byte"}, // Multi-bit

		// Fixed-precision types
		{Pattern: "(?i)^money(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"}, // PostgreSQL, SQL Server
		{Pattern: "(?i)^decimal(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"},
		{Pattern: "(?i)^dec(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"},
		{Pattern: "(?i)^numeric(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"},
		{Pattern: "(?i)^fixed(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"}, // MySQL

		// Floating point types
		{Pattern: "(?i)^double precision(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"},
		{Pattern: "(?i)^double(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"},
		{Pattern: "(?i)^float8(\\(\\d+,?\\d*\\))?$", Type: "float64", NullType: "*float64"}, // PostgreSQL
		{Pattern: "(?i)^float4(\\(\\d+,?\\d*\\))?$", Type: "float32", NullType: "*float32"}, // PostgreSQL
		{Pattern: "(?i)^float(\\(\\d+,?\\d*\\))?$", Type: "float32", NullType: "*float32"},  // MySQL
		{Pattern: "(?i)^real(\\(\\d+,?\\d*\\))?$", Type: "float32", NullType: "*float32"},

		// String types - longer/specific types first
		{Pattern: "(?i)^character\\s+varying(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^varchar(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^character(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^nvarchar(\\(\\d+\\))?$", Type: "string", NullType: "*string"}, // SQL Server
		{Pattern: "(?i)^nchar(\\(\\d+\\))?$", Type: "string", NullType: "*string"},    // SQL Server
		{Pattern: "(?i)^national\\s+char(acter)?(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^national\\s+varchar(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^longtext$", Type: "string", NullType: "*string"},   // MySQL
		{Pattern: "(?i)^mediumtext$", Type: "string", NullType: "*string"}, // MySQL
		{Pattern: "(?i)^tinytext$", Type: "string", NullType: "*string"},   // MySQL
		{Pattern: "(?i)^ntext$", Type: "string", NullType: "*string"},      // SQL Server
		{Pattern: "(?i)^text(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^char(\\(\\d+\\))?$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^varchar2(\\(\\d+\\))?$", Type: "string", NullType: "*string"},  // Oracle
		{Pattern: "(?i)^nvarchar2(\\(\\d+\\))?$", Type: "string", NullType: "*string"}, // Oracle
		{Pattern: "(?i)^citext$", Type: "string", NullType: "*string"},                 // PostgreSQL

		// JSON types
		{Pattern: "(?i)^jsonb$", Type: "string", NullType: "*string"}, // PostgreSQL
		{Pattern: "(?i)^json$", Type: "string", NullType: "*string"},

		// Date and Time types
		{Pattern: "(?i)^timestamptz(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                             // PostgreSQL
		{Pattern: "(?i)^timestamp(\\(\\d+\\))?\\s+with(out)?\\s+time\\s+zone$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"}, // PostgreSQL
		{Pattern: "(?i)^timestamp(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},
		{Pattern: "(?i)^datetime(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                           // MySQL
		{Pattern: "(?i)^datetime2(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                          // SQL Server
		{Pattern: "(?i)^datetimeoffset(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                     // SQL Server
		{Pattern: "(?i)^smalldatetime$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                                   // SQL Server
		{Pattern: "(?i)^timetz(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},                             // PostgreSQL
		{Pattern: "(?i)^time(\\(\\d+\\))?\\s+with(out)?\\s+time\\s+zone$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"}, // PostgreSQL
		{Pattern: "(?i)^time(\\(\\d+\\))?$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},
		{Pattern: "(?i)^date$", Type: "carbon.DateTime", NullType: "*carbon.DateTime", Import: "github.com/goravel/framework/support/carbon"},
		{Pattern: "(?i)^interval$", Type: "string", NullType: "*string"}, // PostgreSQL

		// Range types (PostgreSQL)
		{Pattern: "(?i)^int4range$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^int8range$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^numrange$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^tsrange$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^tstzrange$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^daterange$", Type: "string", NullType: "*string"},

		// Enum types
		{Pattern: "(?i)^enum\\([^)]*\\)$", Type: "string", NullType: "*string"}, // MySQL
		{Pattern: "(?i)^set\\([^)]*\\)$", Type: "string", NullType: "*string"},  // MySQL

		// Binary types - larger types first
		{Pattern: "(?i)^longblob$", Type: "[]byte", NullType: "[]byte"},               // MySQL
		{Pattern: "(?i)^mediumblob$", Type: "[]byte", NullType: "[]byte"},             // MySQL
		{Pattern: "(?i)^tinyblob$", Type: "[]byte", NullType: "[]byte"},               // MySQL
		{Pattern: "(?i)^blob(\\(\\d+\\))?$", Type: "[]byte", NullType: "[]byte"},      // MySQL
		{Pattern: "(?i)^image$", Type: "[]byte", NullType: "[]byte"},                  // SQL Server
		{Pattern: "(?i)^varbinary(\\(\\d+\\))?$", Type: "[]byte", NullType: "[]byte"}, // MySQL/SQL Server
		{Pattern: "(?i)^binary(\\(\\d+\\))?$", Type: "[]byte", NullType: "[]byte"},    // MySQL/SQL Server
		{Pattern: "(?i)^bytea$", Type: "[]byte", NullType: "[]byte"},                  // PostgreSQL

		// Network types (PostgreSQL)
		{Pattern: "(?i)^macaddr8$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^macaddr$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^cidr$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^inet$", Type: "string", NullType: "*string"},

		// Geometric types (PostgreSQL)
		{Pattern: "(?i)^circle$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^polygon$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^path$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^box$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^lseg$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^line$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^point$", Type: "string", NullType: "*string"},

		// UUID/GUID types
		{Pattern: "(?i)^uuid$", Type: "string", NullType: "*string"},             // PostgreSQL
		{Pattern: "(?i)^uniqueidentifier$", Type: "string", NullType: "*string"}, // SQL Server

		// XML and other text types
		{Pattern: "(?i)^xml$", Type: "string", NullType: "*string"},
		{Pattern: "(?i)^rowversion$", Type: "[]byte", NullType: "[]byte"}, // SQL Server

		// Spatial/Geometry types
		{Pattern: "(?i)^geometry$", Type: "string", NullType: "*string"},    // MySQL/PostgreSQL/SQL Server
		{Pattern: "(?i)^geography$", Type: "string", NullType: "*string"},   // SQL Server
		{Pattern: "(?i)^st_geometry$", Type: "string", NullType: "*string"}, // PostgreSQL PostGIS

		// Miscellaneous specialized types
		{Pattern: "(?i)^hstore$", Type: "map[string]string", NullType: "*map[string]string"}, // PostgreSQL
		{Pattern: "(?i)^hierarchyid$", Type: "string", NullType: "*string"},                  // SQL Server

		// SQLite specific
		{Pattern: "(?i)^rowid$", Type: "int64", NullType: "*int64"}, // SQLite

		// Array types (PostgreSQL)
		{Pattern: "(?i)^(.+)\\[\\]$", Type: "string", NullType: "*string"}, // Match any array type like "text[]"

		// Fallback for unknown types
		{Pattern: ".*", Type: "any", NullType: "any"},
	}
}
