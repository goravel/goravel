package postgres

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/goravel/framework/contracts/database/driver"
	databasedb "github.com/goravel/framework/database/db"
	"github.com/goravel/framework/database/schema"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/collect"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

var _ driver.Grammar = &Grammar{}

type Grammar struct {
	attributeCommands []string
	modifiers         []func(driver.Blueprint, driver.ColumnDefinition) string
	prefix            string
	serials           []string
	wrap              *schema.Wrap
}

func NewGrammar(prefix string) *Grammar {
	grammar := &Grammar{
		attributeCommands: []string{schema.CommandComment},
		prefix:            prefix,
		serials:           []string{"bigInteger", "integer", "mediumInteger", "smallInteger", "tinyInteger"},
		wrap:              schema.NewWrap(prefix),
	}
	grammar.modifiers = []func(driver.Blueprint, driver.ColumnDefinition) string{
		grammar.ModifyDefault,
		grammar.ModifyIncrement,
		grammar.ModifyNullable,
		grammar.ModifyGeneratedAsForChange,
		grammar.ModifyGeneratedAs,
	}

	return grammar
}

func (r *Grammar) CompileAdd(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("alter table %s add column %s", r.wrap.Table(blueprint.GetTableName()), r.getColumn(blueprint, command.Column))
}

func (r *Grammar) CompileChange(blueprint driver.Blueprint, command *driver.Command) []string {
	changes := []string{fmt.Sprintf("alter column %s type %s", r.wrap.Column(command.Column.GetName()), schema.ColumnType(r, command.Column))}
	for _, modifier := range r.modifiers {
		if change := modifier(blueprint, command.Column); change != "" {
			changes = append(changes, fmt.Sprintf("alter column %s%s", r.wrap.Column(command.Column.GetName()), change))
		}
	}

	return []string{
		fmt.Sprintf("alter table %s %s", r.wrap.Table(blueprint.GetTableName()), strings.Join(changes, ", ")),
	}
}

func (r *Grammar) CompileColumns(schema, table string) (string, error) {
	schema, table, err := parseSchemaAndTable(table, schema)
	if err != nil {
		return "", err
	}

	table = r.prefix + table

	return fmt.Sprintf(
		"select a.attname as name, t.typname as type_name, format_type(a.atttypid, a.atttypmod) as type, "+
			"(select tc.collcollate from pg_catalog.pg_collation tc where tc.oid = a.attcollation) as collation, "+
			"not a.attnotnull as nullable, "+
			"(select pg_get_expr(adbin, adrelid) from pg_attrdef where c.oid = pg_attrdef.adrelid and pg_attrdef.adnum = a.attnum) as default, "+
			"col_description(c.oid, a.attnum) as comment "+
			"from pg_attribute a, pg_class c, pg_type t, pg_namespace n "+
			"where c.relname = %s and n.nspname = %s and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid and n.oid = c.relnamespace "+
			"order by a.attnum", r.wrap.Quote(table), r.wrap.Quote(schema)), nil
}

func (r *Grammar) CompileComment(blueprint driver.Blueprint, command *driver.Command) string {
	comment := "NULL"
	if command.Column.IsSetComment() {
		comment = r.wrap.Quote(strings.ReplaceAll(command.Column.GetComment(), "'", "''"))
	}

	return fmt.Sprintf("comment on column %s.%s is %s",
		r.wrap.Table(blueprint.GetTableName()),
		r.wrap.Column(command.Column.GetName()),
		comment)
}

func (r *Grammar) CompileCreate(blueprint driver.Blueprint) string {
	return fmt.Sprintf("create table %s (%s)", r.wrap.Table(blueprint.GetTableName()), strings.Join(r.getColumns(blueprint), ", "))
}

func (r *Grammar) CompileDefault(_ driver.Blueprint, _ *driver.Command) string {
	return ""
}

func (r *Grammar) CompileDrop(blueprint driver.Blueprint) string {
	return fmt.Sprintf("drop table %s", r.wrap.Table(blueprint.GetTableName()))
}

func (r *Grammar) CompileDropAllDomains(domains []string) string {
	return fmt.Sprintf("drop domain %s cascade", strings.Join(r.EscapeNames(domains), ", "))
}

func (r *Grammar) CompileDropAllTables(schema string, tables []driver.Table) []string {
	excludedTables := r.EscapeNames([]string{"spatial_ref_sys"})
	escapedSchema := r.EscapeNames([]string{schema})[0]

	var dropTables []string
	for _, table := range tables {
		qualifiedName := fmt.Sprintf("%s.%s", table.Schema, table.Name)

		isExcludedTable := slices.Contains(excludedTables, qualifiedName) || slices.Contains(excludedTables, table.Name)
		isInCurrentSchema := escapedSchema == r.EscapeNames([]string{table.Schema})[0]

		if !isExcludedTable && isInCurrentSchema {
			dropTables = append(dropTables, qualifiedName)
		}
	}

	if len(dropTables) == 0 {
		return nil
	}

	return []string{fmt.Sprintf("drop table %s cascade", strings.Join(r.EscapeNames(dropTables), ", "))}
}

func (r *Grammar) CompileDropAllTypes(schema string, types []driver.Type) []string {
	var dropTypes, dropDomains []string

	for _, t := range types {
		if !t.Implicit && schema == t.Schema {
			if t.Type == "domain" {
				dropDomains = append(dropDomains, fmt.Sprintf("%s.%s", t.Schema, t.Name))
			} else {
				dropTypes = append(dropTypes, fmt.Sprintf("%s.%s", t.Schema, t.Name))
			}
		}
	}

	var sql []string
	if len(dropTypes) > 0 {
		sql = append(sql, fmt.Sprintf("drop type %s cascade", strings.Join(r.EscapeNames(dropTypes), ", ")))
	}
	if len(dropDomains) > 0 {
		sql = append(sql, fmt.Sprintf("drop domain %s cascade", strings.Join(r.EscapeNames(dropDomains), ", ")))
	}

	return sql
}

func (r *Grammar) CompileDropAllViews(schema string, views []driver.View) []string {
	var dropViews []string
	for _, view := range views {
		if schema == view.Schema {
			dropViews = append(dropViews, fmt.Sprintf("%s.%s", view.Schema, view.Name))
		}
	}
	if len(dropViews) == 0 {
		return nil
	}

	return []string{fmt.Sprintf("drop view %s cascade", strings.Join(r.EscapeNames(dropViews), ", "))}
}

func (r *Grammar) CompileDropColumn(blueprint driver.Blueprint, command *driver.Command) []string {
	columns := r.wrap.PrefixArray("drop column", r.wrap.Columns(command.Columns))

	return []string{
		fmt.Sprintf("alter table %s %s", r.wrap.Table(blueprint.GetTableName()), strings.Join(columns, ", ")),
	}
}

func (r *Grammar) CompileDropForeign(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("alter table %s drop constraint %s", r.wrap.Table(blueprint.GetTableName()), r.wrap.Column(command.Index))
}

func (r *Grammar) CompileDropFullText(blueprint driver.Blueprint, command *driver.Command) string {
	return r.CompileDropIndex(blueprint, command)
}

func (r *Grammar) CompileDropIfExists(blueprint driver.Blueprint) string {
	return fmt.Sprintf("drop table if exists %s", r.wrap.Table(blueprint.GetTableName()))
}

func (r *Grammar) CompileDropIndex(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("drop index %s", r.wrap.Column(command.Index))
}

func (r *Grammar) CompileDropPrimary(blueprint driver.Blueprint, command *driver.Command) string {
	tableName := blueprint.GetTableName()
	index := r.wrap.Column(fmt.Sprintf("%s%s_pkey", r.wrap.GetPrefix(), tableName))

	return fmt.Sprintf("alter table %s drop constraint %s", r.wrap.Table(tableName), index)
}

func (r *Grammar) CompileDropUnique(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("alter table %s drop constraint %s", r.wrap.Table(blueprint.GetTableName()), r.wrap.Column(command.Index))
}

func (r *Grammar) CompileForeign(blueprint driver.Blueprint, command *driver.Command) string {
	sql := fmt.Sprintf("alter table %s add constraint %s foreign key (%s) references %s (%s)",
		r.wrap.Table(blueprint.GetTableName()),
		r.wrap.Column(command.Index),
		r.wrap.Columnize(command.Columns),
		r.wrap.Table(command.On),
		r.wrap.Columnize(command.References))
	if command.OnDelete != "" {
		sql += " on delete " + command.OnDelete
	}
	if command.OnUpdate != "" {
		sql += " on update " + command.OnUpdate
	}

	return sql
}

func (r *Grammar) CompileForeignKeys(schema, table string) string {
	return fmt.Sprintf(
		`SELECT 
			c.conname AS name, 
			string_agg(la.attname, ',' ORDER BY conseq.ord) AS columns, 
			fn.nspname AS foreign_schema, 
			fc.relname AS foreign_table, 
			string_agg(fa.attname, ',' ORDER BY conseq.ord) AS foreign_columns, 
			c.confupdtype AS on_update, 
			c.confdeltype AS on_delete 
		FROM pg_constraint c 
		JOIN pg_class tc ON c.conrelid = tc.oid 
		JOIN pg_namespace tn ON tn.oid = tc.relnamespace 
		JOIN pg_class fc ON c.confrelid = fc.oid 
		JOIN pg_namespace fn ON fn.oid = fc.relnamespace 
		JOIN LATERAL unnest(c.conkey) WITH ORDINALITY AS conseq(num, ord) ON TRUE 
		JOIN pg_attribute la ON la.attrelid = c.conrelid AND la.attnum = conseq.num 
		JOIN pg_attribute fa ON fa.attrelid = c.confrelid AND fa.attnum = c.confkey[conseq.ord] 
		WHERE c.contype = 'f' AND tc.relname = %s AND tn.nspname = %s 
		GROUP BY c.conname, fn.nspname, fc.relname, c.confupdtype, c.confdeltype`,
		r.wrap.Quote(table),
		r.wrap.Quote(schema),
	)
}

func (r *Grammar) CompileFullText(blueprint driver.Blueprint, command *driver.Command) string {
	language := "english"
	if command.Language != "" {
		language = command.Language
	}

	columns := collect.Map(command.Columns, func(column string, _ int) string {
		return fmt.Sprintf("to_tsvector(%s, %s)", r.wrap.Quote(language), r.wrap.Column(column))
	})

	return fmt.Sprintf("create index %s on %s using gin(%s)", r.wrap.Column(command.Index), r.wrap.Table(blueprint.GetTableName()), strings.Join(columns, " || "))
}

func (r *Grammar) CompileIndex(blueprint driver.Blueprint, command *driver.Command) string {
	var algorithm string
	if command.Algorithm != "" {
		algorithm = " using " + command.Algorithm
	}

	return fmt.Sprintf("create index %s on %s%s (%s)",
		r.wrap.Column(command.Index),
		r.wrap.Table(blueprint.GetTableName()),
		algorithm,
		r.wrap.Columnize(command.Columns),
	)
}

func (r *Grammar) CompileIndexes(schema, table string) (string, error) {
	schema, table, err := parseSchemaAndTable(table, schema)
	if err != nil {
		return "", err
	}

	table = r.prefix + table

	return fmt.Sprintf(
		"select ic.relname as name, string_agg(a.attname, ',' order by indseq.ord) as columns, "+
			"am.amname as \"type\", i.indisunique as \"unique\", i.indisprimary as \"primary\" "+
			"from pg_index i "+
			"join pg_class tc on tc.oid = i.indrelid "+
			"join pg_namespace tn on tn.oid = tc.relnamespace "+
			"join pg_class ic on ic.oid = i.indexrelid "+
			"join pg_am am on am.oid = ic.relam "+
			"join lateral unnest(i.indkey) with ordinality as indseq(num, ord) on true "+
			"left join pg_attribute a on a.attrelid = i.indrelid and a.attnum = indseq.num "+
			"where tc.relname = %s and tn.nspname = %s "+
			"group by ic.relname, am.amname, i.indisunique, i.indisprimary",
		r.wrap.Quote(table),
		r.wrap.Quote(schema),
	), nil
}

func (r *Grammar) CompileJsonColumnsUpdate(values map[string]any) (map[string]any, error) {
	var (
		compiled = make(map[string]any)
		json     = App.GetJson()
	)

	for key, value := range values {
		if strings.Contains(key, "->") {
			segments := strings.Split(key, "->")
			column := segments[0]
			path := "{" + strings.Join(r.wrap.JsonPathAttributes(segments[1:], `"`), ",") + "}"

			binding, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			expr, ok := compiled[column]
			if !ok {
				expr = databasedb.Raw(r.wrap.Column(column) + "::jsonb")
			}

			compiled[column] = databasedb.Raw("jsonb_set(?,?,?)", expr, path, string(binding))

			continue
		}

		compiled[key] = value
	}

	return compiled, nil
}

func (r *Grammar) CompileJsonContains(column string, value any, isNot bool) (string, []any, error) {
	column = strings.ReplaceAll(r.CompileJsonSelector(column), "->>", "->")
	binding, err := App.GetJson().Marshal(value)
	if err != nil {
		return column, nil, err
	}

	return r.wrap.Not(fmt.Sprintf("(%s)::jsonb @> ?", column), isNot), []any{string(binding)}, nil
}

func (r *Grammar) CompileJsonContainsKey(column string, isNot bool) string {
	segments := strings.Split(column, "->")
	lastSegment := segments[len(segments)-1]
	segments = segments[:len(segments)-1]

	var jsonArrayIndex string
	if _, err := strconv.Atoi(lastSegment); err == nil {
		jsonArrayIndex = lastSegment
	} else if matches := regexp.MustCompile(`\[(-?[0-9]+)]$`).FindStringSubmatch(lastSegment); len(matches) == 2 {
		segments = append(segments, strings.TrimSuffix(lastSegment, matches[0]))
		jsonArrayIndex = matches[1]
	}

	column = strings.ReplaceAll(r.CompileJsonSelector(strings.Join(segments, "->")), "->>", "->")
	if len(jsonArrayIndex) > 0 {
		index := cast.ToInt(jsonArrayIndex)
		if index < 0 {
			index = -index
		} else {
			index = index + 1
		}
		return r.wrap.Not(fmt.Sprintf("case when %s then %s else false end",
			fmt.Sprintf("jsonb_typeof((%s)::jsonb) = 'array'", column),
			fmt.Sprintf("jsonb_array_length((%s)::jsonb) >= %d", column, index),
		), isNot)
	}

	return r.wrap.Not(fmt.Sprintf("coalesce((%s)::jsonb ? %s, false)", column, r.wrap.Quote(strings.ReplaceAll(lastSegment, "'", "''"))), isNot)
}

func (r *Grammar) CompileJsonLength(column string) string {
	column = strings.ReplaceAll(r.CompileJsonSelector(column), "->>", "->")

	return fmt.Sprintf("jsonb_array_length((%s)::jsonb)", column)
}

func (r *Grammar) CompileJsonSelector(column string) string {
	path := strings.Split(column, "->")
	field := r.wrap.Column(path[0])
	if len(path) == 1 {
		return field
	}

	wrappedPath := r.wrap.JsonPathAttributes(path[1:])
	if len(wrappedPath) > 1 {
		return field + "->" + strings.Join(wrappedPath[:len(wrappedPath)-1], "->") + "->>" + wrappedPath[len(wrappedPath)-1]
	}

	return field + "->>" + wrappedPath[0]
}

func (r *Grammar) CompileJsonValues(args ...any) []any {
	for i, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool:
			args[i] = fmt.Sprint(val.Interface())

		case reflect.Slice, reflect.Array:
			if length := val.Len(); length > 0 {
				values := make([]any, length)
				for j := 0; j < length; j++ {
					values[j] = val.Index(j).Interface()
				}
				args[i] = r.CompileJsonValues(values...)
			}
		default:

		}

	}
	return args
}

func (r *Grammar) CompileLockForUpdate(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	if conditions.LockForUpdate != nil && *conditions.LockForUpdate {
		builder = builder.Suffix("FOR UPDATE")
	}

	return builder
}

func (r *Grammar) CompileLockForUpdateForGorm() clause.Expression {
	return clause.Locking{Strength: "UPDATE"}
}

func (r *Grammar) CompilePlaceholderFormat() driver.PlaceholderFormat {
	return sq.Dollar
}

func (r *Grammar) CompilePrimary(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("alter table %s add primary key (%s)", r.wrap.Table(blueprint.GetTableName()), r.wrap.Columnize(command.Columns))
}

func (r *Grammar) CompilePrune(_ string) string {
	return "vacuum full"
}

func (r *Grammar) CompileInRandomOrder(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	if conditions.InRandomOrder != nil && *conditions.InRandomOrder {
		conditions.OrderBy = []string{"RANDOM()"}
	}

	return builder
}

func (r *Grammar) CompileRandomOrderForGorm() string {
	return "RANDOM()"
}

func (r *Grammar) CompileRename(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("alter table %s rename to %s", r.wrap.Table(blueprint.GetTableName()), r.wrap.Table(command.To))
}

func (r *Grammar) CompileRenameColumn(blueprint driver.Blueprint, command *driver.Command, _ []driver.Column) (string, error) {
	return fmt.Sprintf("alter table %s rename column %s to %s",
		r.wrap.Table(blueprint.GetTableName()),
		r.wrap.Column(command.From),
		r.wrap.Column(command.To),
	), nil
}

func (r *Grammar) CompileRenameIndex(blueprint driver.Blueprint, command *driver.Command, _ []driver.Index) []string {
	return []string{
		fmt.Sprintf("alter index %s rename to %s", r.wrap.Column(command.From), r.wrap.Column(command.To)),
	}
}

func (r *Grammar) CompileSharedLock(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	if conditions.SharedLock != nil && *conditions.SharedLock {
		builder = builder.Suffix("FOR SHARE")
	}

	return builder
}

func (r *Grammar) CompileSharedLockForGorm() clause.Expression {
	return clause.Locking{Strength: "SHARE"}
}

func (r *Grammar) CompileTables(_ string) string {
	return "select c.relname as name, n.nspname as schema, pg_total_relation_size(c.oid) as size, " +
		"obj_description(c.oid, 'pg_class') as comment from pg_class c, pg_namespace n " +
		"where c.relkind in ('r', 'p') and n.oid = c.relnamespace and n.nspname not in ('pg_catalog', 'information_schema') " +
		"order by c.relname"
}

func (r *Grammar) CompileTableComment(blueprint driver.Blueprint, command *driver.Command) string {
	return fmt.Sprintf("comment on table %s is '%s'",
		r.wrap.Table(blueprint.GetTableName()),
		strings.ReplaceAll(command.Value, "'", "''"),
	)
}

func (r *Grammar) CompileTypes() string {
	return `select t.typname as name, n.nspname as schema, t.typtype as type, t.typcategory as category, 
		((t.typinput = 'array_in'::regproc and t.typoutput = 'array_out'::regproc) or t.typtype = 'm') as implicit 
		from pg_type t 
		join pg_namespace n on n.oid = t.typnamespace 
		left join pg_class c on c.oid = t.typrelid 
		left join pg_type el on el.oid = t.typelem 
		left join pg_class ce on ce.oid = el.typrelid 
		where ((t.typrelid = 0 and (ce.relkind = 'c' or ce.relkind is null)) or c.relkind = 'c') 
		and not exists (select 1 from pg_depend d where d.objid in (t.oid, t.typelem) and d.deptype = 'e') 
		and n.nspname not in ('pg_catalog', 'information_schema')`
}

func (r *Grammar) CompileUnique(blueprint driver.Blueprint, command *driver.Command) string {
	sql := fmt.Sprintf("alter table %s add constraint %s unique (%s)",
		r.wrap.Table(blueprint.GetTableName()),
		r.wrap.Column(command.Index),
		r.wrap.Columnize(command.Columns))

	if command.Deferrable != nil {
		if *command.Deferrable {
			sql += " deferrable"
		} else {
			sql += " not deferrable"
		}
	}
	if command.Deferrable != nil && command.InitiallyImmediate != nil {
		if *command.InitiallyImmediate {
			sql += " initially immediate"
		} else {
			sql += " initially deferred"
		}
	}

	return sql
}

func (r *Grammar) CompileVersion() string {
	return "SELECT current_setting('server_version') AS value;"
}

func (r *Grammar) CompileViews(database string) string {
	return "select viewname as name, schemaname as schema, definition from pg_views where schemaname not in ('pg_catalog', 'information_schema') order by viewname"
}

func (r *Grammar) EscapeNames(names []string) []string {
	escapedNames := make([]string, 0, len(names))

	for _, name := range names {
		segments := strings.Split(name, ".")
		for i, segment := range segments {
			segments[i] = strings.Trim(segment, `'"`)
		}
		escapedName := `"` + strings.Join(segments, `"."`) + `"`
		escapedNames = append(escapedNames, escapedName)
	}

	return escapedNames
}

func (r *Grammar) GetAttributeCommands() []string {
	return r.attributeCommands
}

func (r *Grammar) ModifyDefault(_ driver.Blueprint, column driver.ColumnDefinition) string {
	if column.IsChange() {
		if column.GetAutoIncrement() || column.IsSetGeneratedAs() {
			return ""
		}
		if column.GetDefault() != nil {
			return fmt.Sprintf(" set default %s", schema.ColumnDefaultValue(column.GetDefault()))
		}
		return " drop default"
	}
	if column.GetDefault() != nil {
		return fmt.Sprintf(" default %s", schema.ColumnDefaultValue(column.GetDefault()))
	}

	return ""
}

func (r *Grammar) ModifyGeneratedAs(_ driver.Blueprint, column driver.ColumnDefinition) string {
	if !column.IsSetGeneratedAs() {
		return ""
	}

	option := "by default"
	if column.IsAlways() {
		option = "always"
	}

	identity := ""
	if generatedAs := column.GetGeneratedAs(); len(generatedAs) > 0 {
		identity = " (" + generatedAs + ")"
	}

	sql := fmt.Sprintf(" generated %s as identity%s", option, identity)
	if column.IsChange() {
		sql = " add" + sql
	}

	return sql

}

func (r *Grammar) ModifyGeneratedAsForChange(_ driver.Blueprint, column driver.ColumnDefinition) string {
	if column.IsChange() && column.IsSetGeneratedAs() && !column.GetAutoIncrement() {
		return " drop identity if exists"
	}

	return ""
}

func (r *Grammar) ModifyNullable(_ driver.Blueprint, column driver.ColumnDefinition) string {
	if column.IsChange() {
		if column.GetNullable() {
			return " drop not null"
		}
		return " set not null"
	}
	if column.GetNullable() {
		return " null"
	}
	return " not null"
}

func (r *Grammar) ModifyIncrement(blueprint driver.Blueprint, column driver.ColumnDefinition) string {
	if !column.IsChange() &&
		!blueprint.HasCommand("primary") &&
		(slices.Contains(r.serials, column.GetType()) || column.IsSetGeneratedAs()) &&
		column.GetAutoIncrement() {
		return " primary key"
	}

	return ""
}

func (r *Grammar) TypeBigInteger(column driver.ColumnDefinition) string {
	if column.GetAutoIncrement() && !column.IsChange() && !column.IsSetGeneratedAs() {
		return "bigserial"
	}

	return "bigint"
}

func (r *Grammar) TypeBoolean(column driver.ColumnDefinition) string {
	return "boolean"
}

func (r *Grammar) TypeChar(column driver.ColumnDefinition) string {
	length := column.GetLength()
	if length > 0 {
		return fmt.Sprintf("char(%d)", length)
	}

	return "char"
}

func (r *Grammar) TypeDate(column driver.ColumnDefinition) string {
	return "date"
}

func (r *Grammar) TypeDateTime(column driver.ColumnDefinition) string {
	return r.TypeTimestamp(column)
}

func (r *Grammar) TypeDateTimeTz(column driver.ColumnDefinition) string {
	return r.TypeTimestampTz(column)
}

func (r *Grammar) TypeDecimal(column driver.ColumnDefinition) string {
	return fmt.Sprintf("decimal(%d, %d)", column.GetTotal(), column.GetPlaces())
}

func (r *Grammar) TypeDouble(column driver.ColumnDefinition) string {
	return "double precision"
}

func (r *Grammar) TypeEnum(column driver.ColumnDefinition) string {
	return fmt.Sprintf(`varchar(255) check ("%s" in (%s))`, column.GetName(), strings.Join(r.wrap.Quotes(cast.ToStringSlice(column.GetAllowed())), ", "))
}

func (r *Grammar) TypeFloat(column driver.ColumnDefinition) string {
	precision := column.GetPrecision()
	if precision > 0 {
		return fmt.Sprintf("float(%d)", precision)
	}

	return "float"
}

func (r *Grammar) TypeInteger(column driver.ColumnDefinition) string {
	if column.GetAutoIncrement() && !column.IsChange() && !column.IsSetGeneratedAs() {
		return "serial"
	}

	return "integer"
}

func (r *Grammar) TypeJson(column driver.ColumnDefinition) string {
	return "json"
}

func (r *Grammar) TypeJsonb(column driver.ColumnDefinition) string {
	return "jsonb"
}

func (r *Grammar) TypeLongText(column driver.ColumnDefinition) string {
	return "text"
}

func (r *Grammar) TypeMediumInteger(column driver.ColumnDefinition) string {
	return r.TypeInteger(column)
}

func (r *Grammar) TypeMediumText(column driver.ColumnDefinition) string {
	return "text"
}

func (r *Grammar) TypeSmallInteger(column driver.ColumnDefinition) string {
	if column.GetAutoIncrement() && !column.IsChange() && !column.IsSetGeneratedAs() {
		return "smallserial"
	}

	return "smallint"
}

func (r *Grammar) TypeString(column driver.ColumnDefinition) string {
	length := column.GetLength()
	if length > 0 {
		return fmt.Sprintf("varchar(%d)", length)
	}

	return "varchar"
}

func (r *Grammar) TypeText(column driver.ColumnDefinition) string {
	return "text"
}

func (r *Grammar) TypeTime(column driver.ColumnDefinition) string {
	return fmt.Sprintf("time(%d) without time zone", column.GetPrecision())
}

func (r *Grammar) TypeTimeTz(column driver.ColumnDefinition) string {
	return fmt.Sprintf("time(%d) with time zone", column.GetPrecision())
}

func (r *Grammar) TypeTimestamp(column driver.ColumnDefinition) string {
	if column.GetUseCurrent() {
		column.Default(schema.Expression("CURRENT_TIMESTAMP"))
	}

	return fmt.Sprintf("timestamp(%d) without time zone", column.GetPrecision())
}

func (r *Grammar) TypeTimestampTz(column driver.ColumnDefinition) string {
	if column.GetUseCurrent() {
		column.Default(schema.Expression("CURRENT_TIMESTAMP"))
	}

	return fmt.Sprintf("timestamp(%d) with time zone", column.GetPrecision())
}

func (r *Grammar) TypeTinyInteger(column driver.ColumnDefinition) string {
	return r.TypeSmallInteger(column)
}

func (r *Grammar) TypeTinyText(column driver.ColumnDefinition) string {
	return "varchar(255)"
}

func (r *Grammar) TypeUuid(column driver.ColumnDefinition) string {
	return "uuid"
}

func (r *Grammar) getColumns(blueprint driver.Blueprint) []string {
	var columns []string
	for _, column := range blueprint.GetAddedColumns() {
		columns = append(columns, r.getColumn(blueprint, column))
	}

	return columns
}

func (r *Grammar) getColumn(blueprint driver.Blueprint, column driver.ColumnDefinition) string {
	sql := fmt.Sprintf("%s %s", r.wrap.Column(column.GetName()), schema.ColumnType(r, column))

	for _, modifier := range r.modifiers {
		sql += modifier(blueprint, column)
	}

	return sql
}

func parseSchemaAndTable(reference, defaultSchema string) (string, string, error) {
	if reference == "" {
		return "", "", errors.SchemaEmptyReferenceString
	}

	parts := strings.Split(reference, ".")
	if len(parts) > 2 {
		return "", "", errors.SchemaErrorReferenceFormat
	}

	schema := defaultSchema
	if len(parts) == 2 {
		schema = parts[0]
		parts = parts[1:]
	}

	table := parts[0]

	return schema, table, nil
}
