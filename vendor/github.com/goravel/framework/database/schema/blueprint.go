package schema

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/support/convert"
)

const (
	CommandAdd          = "add"
	CommandComment      = "comment"
	CommandCreate       = "create"
	CommandDefault      = "default"
	CommandDrop         = "drop"
	CommandDropColumn   = "dropColumn"
	CommandDropForeign  = "dropForeign"
	CommandDropFullText = "dropFullText"
	CommandDropIfExists = "dropIfExists"
	CommandDropIndex    = "dropIndex"
	CommandDropPrimary  = "dropPrimary"
	CommandDropUnique   = "dropUnique"
	CommandForeign      = "foreign"
	CommandFullText     = "fullText"
	CommandIndex        = "index"
	CommandPrimary      = "primary"
	CommandRename       = "rename"
	CommandRenameColumn = "renameColumn"
	CommandRenameIndex  = "renameIndex"
	CommandTableComment = "tableComment"
	CommandUnique       = "unique"
	DefaultStringLength = 255
	DefaultUlidLength   = 26
)

type Blueprint struct {
	schema   schema.Schema
	prefix   string
	table    string
	columns  []*ColumnDefinition
	commands []*driver.Command
}

func NewBlueprint(schema schema.Schema, prefix, table string) *Blueprint {
	return &Blueprint{
		prefix: prefix,
		schema: schema,
		table:  table,
	}
}

func (r *Blueprint) BigIncrements(column string) driver.ColumnDefinition {
	return r.UnsignedBigInteger(column).AutoIncrement()
}

func (r *Blueprint) BigInteger(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("bigInteger", column)
}

func (r *Blueprint) Boolean(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("boolean", column)
}

func (r *Blueprint) Build(query orm.Query, grammar driver.Grammar) error {
	statements, err := r.ToSql(grammar)
	if err != nil {
		return err
	}

	for _, sql := range statements {
		if _, err = query.Exec(sql); err != nil {
			return err
		}
	}

	return nil
}

func (r *Blueprint) Char(column string, length ...int) driver.ColumnDefinition {
	defaultLength := DefaultStringLength
	if len(length) > 0 {
		defaultLength = length[0]
	}

	columnImpl := r.createAndAddColumn("char", column)
	columnImpl.length = &defaultLength

	return columnImpl
}

func (r *Blueprint) Column(column, ttype string) driver.ColumnDefinition {
	return r.createAndAddColumn(ttype, column)
}

func (r *Blueprint) Comment(comment string) {
	r.addCommand(&driver.Command{
		Name:  CommandTableComment,
		Value: comment,
	})
}

func (r *Blueprint) Create() {
	r.addCommand(&driver.Command{
		Name: CommandCreate,
	})
}

func (r *Blueprint) Decimal(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("decimal", column)
}

func (r *Blueprint) Date(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("date", column)
}

func (r *Blueprint) DateTime(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("dateTime", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) DateTimes(precision ...int) {
	_ = r.DateTime("created_at", precision...).Nullable()
	_ = r.DateTime("updated_at", precision...).Nullable()
}
func (r *Blueprint) DateTimeTz(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("dateTimeTz", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) Double(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("double", column)
}

func (r *Blueprint) Drop() {
	r.addCommand(&driver.Command{
		Name: CommandDrop,
	})
}

func (r *Blueprint) DropColumn(column ...string) {
	r.addCommand(&driver.Command{
		Name:    CommandDropColumn,
		Columns: column,
	})
}

func (r *Blueprint) DropForeign(column ...string) {
	r.indexCommand(CommandDropForeign, column, schema.IndexConfig{
		Name: r.createIndexName(CommandForeign, column),
	})
}

func (r *Blueprint) DropForeignByName(name string) {
	r.indexCommand(CommandDropForeign, nil, schema.IndexConfig{
		Name: name,
	})
}

func (r *Blueprint) DropFullText(column ...string) {
	r.indexCommand(CommandDropFullText, column, schema.IndexConfig{
		Name: r.createIndexName(CommandFullText, column),
	})
}

func (r *Blueprint) DropFullTextByName(name string) {
	r.indexCommand(CommandDropFullText, nil, schema.IndexConfig{
		Name: name,
	})
}

func (r *Blueprint) DropIfExists() {
	r.addCommand(&driver.Command{
		Name: CommandDropIfExists,
	})
}

func (r *Blueprint) DropIndex(column ...string) {
	r.indexCommand(CommandDropIndex, column, schema.IndexConfig{
		Name: r.createIndexName(CommandIndex, column),
	})
}

func (r *Blueprint) DropIndexByName(name string) {
	r.indexCommand(CommandDropIndex, nil, schema.IndexConfig{
		Name: name,
	})
}

func (r *Blueprint) DropPrimary(column ...string) {
	r.indexCommand(CommandDropPrimary, column, schema.IndexConfig{
		Name: r.createIndexName(CommandPrimary, column),
	})
}

func (r *Blueprint) DropSoftDeletes(column ...string) {
	if len(column) > 0 {
		r.DropColumn(column[0])
	} else {
		r.DropColumn("deleted_at")
	}
}

func (r *Blueprint) DropSoftDeletesTz(column ...string) {
	r.DropSoftDeletes(column...)
}

func (r *Blueprint) DropTimestamps() {
	r.DropColumn("created_at", "updated_at")
}

func (r *Blueprint) DropTimestampsTz() {
	r.DropTimestamps()
}

func (r *Blueprint) DropUnique(column ...string) {
	r.indexCommand(CommandDropUnique, column, schema.IndexConfig{
		Name: r.createIndexName(CommandUnique, column),
	})
}

func (r *Blueprint) DropUniqueByName(name string) {
	r.indexCommand(CommandDropUnique, nil, schema.IndexConfig{
		Name: name,
	})
}

func (r *Blueprint) Enum(column string, allowed []any) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("enum", column)
	columnImpl.allowed = allowed

	return columnImpl
}

func (r *Blueprint) Float(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("float", column)
	columnImpl.precision = convert.Pointer(53)

	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) Foreign(column ...string) schema.ForeignKeyDefinition {
	command := r.indexCommand(CommandForeign, column)

	return NewForeignKeyDefinition(command)
}

func (r *Blueprint) ForeignID(column string) schema.ForeignIDColumnDefinition {
	return &ForeignIDColumnDefinition{
		ColumnDefinition: r.UnsignedBigInteger(column).(*ColumnDefinition),
		blueprint:        r,
	}
}

func (r *Blueprint) ForeignUlid(column string, length ...int) schema.ForeignIDColumnDefinition {
	return &ForeignIDColumnDefinition{
		ColumnDefinition: r.Ulid(column, length...).(*ColumnDefinition),
		blueprint:        r,
	}
}

func (r *Blueprint) ForeignUuid(column string) schema.ForeignIDColumnDefinition {
	return &ForeignIDColumnDefinition{
		ColumnDefinition: r.Uuid(column).(*ColumnDefinition),
		blueprint:        r,
	}
}

func (r *Blueprint) FullText(column ...string) schema.IndexDefinition {
	command := r.indexCommand(CommandFullText, column)

	return NewIndexDefinition(command)
}

func (r *Blueprint) GetAddedColumns() []driver.ColumnDefinition {
	var columns []driver.ColumnDefinition
	for _, column := range r.columns {
		columns = append(columns, column)
	}

	return columns
}

func (r *Blueprint) GetCommands() []*driver.Command {
	return r.commands
}

func (r *Blueprint) GetTableName() string {
	return r.table
}

func (r *Blueprint) HasCommand(command string) bool {
	for _, c := range r.commands {
		if c.Name == command {
			return true
		}
	}

	return false
}

func (r *Blueprint) ID(column ...string) driver.ColumnDefinition {
	if len(column) > 0 {
		return r.BigIncrements(column[0])
	}

	return r.BigIncrements("id")
}

func (r *Blueprint) Increments(column string) driver.ColumnDefinition {
	return r.IntegerIncrements(column)
}

func (r *Blueprint) Index(column ...string) schema.IndexDefinition {
	command := r.indexCommand(CommandIndex, column)

	return NewIndexDefinition(command)
}

func (r *Blueprint) Integer(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("integer", column)
}

func (r *Blueprint) IntegerIncrements(column string) driver.ColumnDefinition {
	return r.UnsignedInteger(column).AutoIncrement()
}

func (r *Blueprint) Json(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("json", column)
}

func (r *Blueprint) Jsonb(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("jsonb", column)
}

func (r *Blueprint) LongText(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("longText", column)
}

func (r *Blueprint) MediumIncrements(column string) driver.ColumnDefinition {
	return r.UnsignedMediumInteger(column).AutoIncrement()
}

func (r *Blueprint) MediumInteger(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("mediumInteger", column)
}

func (r *Blueprint) MediumText(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("mediumText", column)
}

func (r *Blueprint) Morphs(name string, indexName ...string) {
	switch GetDefaultMorphKeyType() {
	case MorphKeyTypeUuid:
		r.UuidMorphs(name, indexName...)
	case MorphKeyTypeUlid:
		r.UlidMorphs(name, indexName...)
	default:
		r.NumericMorphs(name, indexName...)
	}
}

func (r *Blueprint) NullableMorphs(name string, indexName ...string) {
	r.String(name + "_type").Nullable()

	switch GetDefaultMorphKeyType() {
	case MorphKeyTypeUuid:
		r.Uuid(name + "_id").Nullable()
	case MorphKeyTypeUlid:
		r.Ulid(name + "_id").Nullable()
	default:
		r.UnsignedBigInteger(name + "_id").Nullable()
	}

	r.createMorphIndex(name, indexName...)
}

func (r *Blueprint) NumericMorphs(name string, indexName ...string) {
	r.String(name + "_type")
	r.UnsignedBigInteger(name + "_id")
	r.createMorphIndex(name, indexName...)
}

func (r *Blueprint) Primary(column ...string) {
	r.indexCommand(CommandPrimary, column)
}

func (r *Blueprint) Rename(to string) {
	command := &driver.Command{
		Name: CommandRename,
		To:   to,
	}

	r.addCommand(command)
}

func (r *Blueprint) RenameColumn(from, to string) {
	command := &driver.Command{
		Name: CommandRenameColumn,
		From: from,
		To:   to,
	}

	r.addCommand(command)
}

func (r *Blueprint) RenameIndex(from, to string) {
	command := &driver.Command{
		Name: CommandRenameIndex,
		From: from,
		To:   to,
	}

	r.addCommand(command)
}

func (r *Blueprint) SetTable(name string) {
	r.table = name
}

func (r *Blueprint) SmallIncrements(column string) driver.ColumnDefinition {
	return r.UnsignedSmallInteger(column).AutoIncrement()
}

func (r *Blueprint) SmallInteger(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("smallInteger", column)
}

func (r *Blueprint) SoftDeletes(column ...string) driver.ColumnDefinition {
	newColumn := "deleted_at"
	if len(column) > 0 {
		newColumn = column[0]
	}

	return r.Timestamp(newColumn).Nullable()
}

func (r *Blueprint) SoftDeletesTz(column ...string) driver.ColumnDefinition {
	newColumn := "deleted_at"
	if len(column) > 0 {
		newColumn = column[0]
	}

	return r.TimestampTz(newColumn).Nullable()
}

func (r *Blueprint) String(column string, length ...int) driver.ColumnDefinition {
	defaultLength := DefaultStringLength
	if len(length) > 0 {
		defaultLength = length[0]
	}

	columnImpl := r.createAndAddColumn("string", column)
	columnImpl.length = &defaultLength

	return columnImpl
}

func (r *Blueprint) Text(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("text", column)
}

func (r *Blueprint) Time(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("time", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) TimeTz(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("timeTz", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) Timestamp(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("timestamp", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) Timestamps(precision ...int) {
	r.Timestamp("created_at", precision...).Nullable()
	r.Timestamp("updated_at", precision...).Nullable()
}

func (r *Blueprint) TimestampsTz(precision ...int) {
	r.TimestampTz("created_at", precision...).Nullable()
	r.TimestampTz("updated_at", precision...).Nullable()
}

func (r *Blueprint) TimestampTz(column string, precision ...int) driver.ColumnDefinition {
	columnImpl := r.createAndAddColumn("timestampTz", column)
	if len(precision) > 0 {
		columnImpl.precision = &precision[0]
	}

	return columnImpl
}

func (r *Blueprint) TinyIncrements(column string) driver.ColumnDefinition {
	return r.UnsignedTinyInteger(column).AutoIncrement()
}

func (r *Blueprint) TinyInteger(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("tinyInteger", column)
}

func (r *Blueprint) TinyText(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("tinyText", column)
}

func (r *Blueprint) ToSql(grammar driver.Grammar) ([]string, error) {
	r.addImpliedCommands(grammar)

	var statements []string
	for _, command := range r.commands {
		if command.ShouldBeSkipped {
			continue
		}

		switch command.Name {
		case CommandAdd:
			if command.Column.IsChange() {
				if statement := grammar.CompileChange(r, command); len(statement) > 0 {
					statements = append(statements, statement...)
				}
				continue
			}
			statements = append(statements, grammar.CompileAdd(r, command))
		case CommandComment:
			if statement := grammar.CompileComment(r, command); statement != "" {
				statements = append(statements, statement)
			}
		case CommandCreate:
			statements = append(statements, grammar.CompileCreate(r))
		case CommandDefault:
			if statement := grammar.CompileDefault(r, command); statement != "" {
				statements = append(statements, statement)
			}
		case CommandDrop:
			statements = append(statements, grammar.CompileDrop(r))
		case CommandDropColumn:
			statements = append(statements, grammar.CompileDropColumn(r, command)...)
		case CommandDropForeign:
			statements = append(statements, grammar.CompileDropForeign(r, command))
		case CommandDropFullText:
			statements = append(statements, grammar.CompileDropFullText(r, command))
		case CommandDropIfExists:
			statements = append(statements, grammar.CompileDropIfExists(r))
		case CommandDropIndex:
			statements = append(statements, grammar.CompileDropIndex(r, command))
		case CommandDropPrimary:
			statements = append(statements, grammar.CompileDropPrimary(r, command))
		case CommandDropUnique:
			statements = append(statements, grammar.CompileDropUnique(r, command))
		case CommandForeign:
			statements = append(statements, grammar.CompileForeign(r, command))
		case CommandFullText:
			statements = append(statements, grammar.CompileFullText(r, command))
		case CommandIndex:
			statements = append(statements, grammar.CompileIndex(r, command))
		case CommandPrimary:
			statements = append(statements, grammar.CompilePrimary(r, command))
		case CommandRename:
			statements = append(statements, grammar.CompileRename(r, command))
		case CommandRenameColumn:
			columns, err := r.schema.GetColumns(r.GetTableName())
			if err != nil {
				return statements, err
			}
			statement, err := grammar.CompileRenameColumn(r, command, columns)
			if err != nil {
				return statements, err
			}
			statements = append(statements, statement)
		case CommandRenameIndex:
			indexes, err := r.schema.GetIndexes(r.GetTableName())
			if err != nil {
				return statements, err
			}
			statements = append(statements, grammar.CompileRenameIndex(r, command, indexes)...)
		case CommandTableComment:
			if statement := grammar.CompileTableComment(r, command); statement != "" {
				statements = append(statements, statement)
			}
		case CommandUnique:
			statements = append(statements, grammar.CompileUnique(r, command))
		}
	}

	return statements, nil
}

func (r *Blueprint) Unique(column ...string) schema.IndexDefinition {
	command := r.indexCommand(CommandUnique, column)

	return NewIndexDefinition(command)
}

func (r *Blueprint) UnsignedBigInteger(column string) driver.ColumnDefinition {
	return r.BigInteger(column).Unsigned()
}

func (r *Blueprint) UnsignedInteger(column string) driver.ColumnDefinition {
	return r.Integer(column).Unsigned()
}

func (r *Blueprint) UnsignedMediumInteger(column string) driver.ColumnDefinition {
	return r.MediumInteger(column).Unsigned()
}

func (r *Blueprint) UnsignedSmallInteger(column string) driver.ColumnDefinition {
	return r.SmallInteger(column).Unsigned()
}

func (r *Blueprint) UnsignedTinyInteger(column string) driver.ColumnDefinition {
	return r.TinyInteger(column).Unsigned()
}

func (r *Blueprint) Ulid(column string, length ...int) driver.ColumnDefinition {
	defaultLength := DefaultUlidLength
	if len(length) > 0 {
		defaultLength = length[0]
	}

	return r.Char(column, defaultLength)
}

func (r *Blueprint) UlidMorphs(name string, indexName ...string) {
	r.String(name + "_type")
	r.Ulid(name + "_id")
	r.createMorphIndex(name, indexName...)
}

func (r *Blueprint) Uuid(column string) driver.ColumnDefinition {
	return r.createAndAddColumn("uuid", column)
}

func (r *Blueprint) UuidMorphs(name string, indexName ...string) {
	r.String(name + "_type")
	r.Uuid(name + "_id")
	r.createMorphIndex(name, indexName...)
}

func (r *Blueprint) addAttributeCommands(grammar driver.Grammar) {
	attributeCommands := grammar.GetAttributeCommands()
	for _, column := range r.columns {
		for _, command := range attributeCommands {
			if command == CommandComment && (column.comment != nil || column.change) {
				r.addCommand(&driver.Command{
					Column: column,
					Name:   CommandComment,
				})
			}
			if command == CommandDefault && column.def != nil {
				r.addCommand(&driver.Command{
					Column: column,
					Name:   CommandDefault,
				})
			}
		}
	}
}

func (r *Blueprint) addCommand(command *driver.Command) {
	r.commands = append(r.commands, command)
}

func (r *Blueprint) addImpliedCommands(grammar driver.Grammar) {
	r.addAttributeCommands(grammar)
}

func (r *Blueprint) createAndAddColumn(ttype, name string) *ColumnDefinition {
	columnImpl := &ColumnDefinition{
		name:  &name,
		ttype: convert.Pointer(ttype),
	}

	r.columns = append(r.columns, columnImpl)

	if !r.isCreate() {
		r.addCommand(&driver.Command{
			Name:   CommandAdd,
			Column: columnImpl,
		})
	}

	return columnImpl
}

func (r *Blueprint) createIndexName(ttype string, columns []string) string {
	var table string
	if strings.Contains(r.table, ".") {
		lastDotIndex := strings.LastIndex(r.table, ".")
		table = r.table[:lastDotIndex+1] + r.prefix + r.table[lastDotIndex+1:]
	} else {
		table = r.prefix + r.table
	}

	index := strings.ToLower(fmt.Sprintf("%s_%s_%s", table, strings.Join(columns, "_"), ttype))

	index = strings.ReplaceAll(index, "-", "_")
	index = strings.ReplaceAll(index, ".", "_")

	return index
}

// createMorphIndex creates an index for morph columns with optional custom name
func (r *Blueprint) createMorphIndex(name string, indexName ...string) {
	if len(indexName) > 0 && indexName[0] != "" {
		r.Index(name+"_type", name+"_id").Name(indexName[0])
	} else {
		r.Index(name+"_type", name+"_id")
	}
}

func (r *Blueprint) indexCommand(name string, columns []string, config ...schema.IndexConfig) *driver.Command {
	command := &driver.Command{
		Columns: columns,
		Name:    name,
	}

	if len(config) > 0 {
		command.Algorithm = config[0].Algorithm
		command.Index = config[0].Name
		command.Language = config[0].Language
	} else {
		command.Index = r.createIndexName(name, columns)
	}

	r.addCommand(command)

	return command
}

func (r *Blueprint) isCreate() bool {
	for _, command := range r.commands {
		if command.Name == CommandCreate {
			return true
		}
	}

	return false
}
