package migration

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm/schema"

	contractsschema "github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/str"
)

const tablePrefix = "table."

// dataTypeMapping maps GORM DataType strings to blueprint method names.
// Used for custom TYPE tags that aren't standard GORM DataTypes.
var dataTypeMapping = map[string]string{
	"jsonb":     contractsschema.MethodJsonb,
	"json":      contractsschema.MethodJson,
	"text":      contractsschema.MethodText,
	"binary":    contractsschema.MethodBinary,
	"varbinary": contractsschema.MethodBinary,
	"blob":      contractsschema.MethodBinary,
	"decimal":   contractsschema.MethodDecimal,
	"numeric":   contractsschema.MethodDecimal,
	"uuid":      contractsschema.MethodUuid,
	"ulid":      contractsschema.MethodUlid,
	"date":      contractsschema.MethodDate,
	"time":      contractsschema.MethodTime,
}

// stringTypePrefixes contains SQL string type prefixes for type detection.
var stringTypePrefixes = []string{"char", "varchar", "nvarchar", "varchar2", "nchar"}

var (
	schemaCache           = &sync.Map{}
	defaultNamingStrategy = schema.NamingStrategy{}
)

// Generate creates migration schema lines from a model struct.
// It parses the model using GORM's schema parser and generates Blueprint method calls.
// Returns the table name, a slice of schema field lines, and an error if the model is invalid.
func Generate(model any) (string, []string, error) {
	sch, err := schema.Parse(model, schemaCache, defaultNamingStrategy)
	if err != nil {
		return "", nil, err
	}

	lines := renderSchema(sch)
	if len(lines) == 0 {
		return "", nil, errors.SchemaInvalidModel
	}

	return sch.Table, lines, nil
}

// renderSchema converts a GORM schema into Blueprint method call strings.
// It processes all fields and indexes from the schema, generating the corresponding
// table column definitions and index/constraint definitions for use in migrations.
func renderSchema(sch *schema.Schema) []string {
	var lines []string
	seen := make(map[string]bool)
	var fields []*schema.Field

	// 1. Render Columns
	for _, field := range sch.Fields {
		if shouldSkipField(field) || seen[field.DBName] {
			continue
		}
		seen[field.DBName] = true
		fields = append(fields, field)

		if line := renderField(field); line != "" {
			lines = append(lines, line)
		}
	}

	// 2. Render Indexes
	if idxLines := renderIndexes(sch, fields); len(idxLines) > 0 {
		lines = append(lines, "", strings.Join(idxLines, "\n"))
	}

	return lines
}

func shouldSkipField(field *schema.Field) bool {
	// Skip ignored or embedded fields
	if field.IgnoreMigration || field.EmbeddedSchema != nil || field.DBName == "" {
		return true
	}

	// Skip relation fields (e.g., User in "User `gorm:foreignKey:UserID`")
	// But DO NOT skip foreign key columns (e.g., UserID) they need to exist in the table
	relationships := &field.Schema.Relationships
	relationships.Mux.RLock()
	_, isRel := relationships.Relations[field.Name]
	relationships.Mux.RUnlock()

	return isRel
}

// renderField generates a single column definition with modifiers for a GORM field.
// It builds a Blueprint method chain including the column type, nullability, default value,
// unsigned modifier, and comments based on the field's properties and tags.
func renderField(f *schema.Field) string {
	method, args := fieldToMethod(f)
	if method == "" {
		return ""
	}

	b := &atom{Builder: strings.Builder{}}
	b.Grow(64)
	b.WriteString(tablePrefix)

	// Chain: table.Method("name", args...).Nullable()...
	b.WriteMethod(method, append([]any{f.DBName}, args...)...)

	// Type-specific modifiers first
	if method == contractsschema.MethodDecimal {
		if f.Scale > 0 {
			b.WriteMethod(contractsschema.MethodPlaces, f.Scale)
		}
		if f.Precision > 0 {
			b.WriteMethod(contractsschema.MethodTotal, f.Precision)
		}
	}

	// Check for unsigned modifier
	rawType := strings.ToLower(string(f.DataType))
	if !strings.Contains(method, "Unsigned") {
		if strings.Contains(rawType, "unsigned") || f.TagSettings["UNSIGNED"] != "" {
			b.WriteMethod(contractsschema.MethodUnsigned)
		}
	}

	// Nullable if:
	// - Pointer type (can be nil in Go), OR
	// - sql.Null* types (designed for nullable columns)
	// And not explicitly NOT NULL or primary key
	isNullable := f.FieldType.Kind() == reflect.Ptr || isSQLNullType(f.FieldType)
	if isNullable && !f.NotNull && !f.PrimaryKey {
		b.WriteMethod(contractsschema.MethodNullable)
	}
	if f.HasDefaultValue && f.DefaultValueInterface != nil {
		b.WriteMethod(contractsschema.MethodDefault, f.DefaultValueInterface)
	}
	if f.Comment != "" {
		b.WriteMethod(contractsschema.MethodComment, trimQuotes(f.Comment))
	}

	return b.String()
}

// fieldToMethod maps a GORM field to the appropriate Blueprint method name and arguments.
// It analyzes the field's data type, size, and attributes to determine the correct
// schema method (e.g., String, Integer, Boolean, TimestampTz) and any required parameters.
//
// Priority order:
//  1. migration tag (e.g., `gorm:"migration:Json"`) - highest priority, direct method name
//  2. Primary key with auto increment
//  3. DataType and type inference
func fieldToMethod(f *schema.Field) (string, []any) {
	// Highest priority: explicit migration tag (e.g., `gorm:"migration:Json"` or `gorm:"migration:long_text"`)
	// Allows direct Blueprint method specification. Normalizes to StudlyCase.
	if method, ok := f.TagSettings["MIGRATION"]; ok && method != "" {
		return str.Of(method).Studly().String(), nil
	}

	if f.PrimaryKey && f.AutoIncrement {
		if f.Size <= 32 {
			return contractsschema.MethodIncrements, nil
		}
		return contractsschema.MethodBigIncrements, nil
	}

	switch f.DataType {
	case schema.Bool:
		return contractsschema.MethodBoolean, nil
	case schema.Int:
		return intMethod(f.Size, false), nil
	case schema.Uint:
		return intMethod(f.Size, true), nil
	case schema.Float:
		if f.Size <= 32 {
			return contractsschema.MethodFloat, nil
		}
		return contractsschema.MethodDouble, nil
	case schema.String:
		if f.Size > 0 {
			return contractsschema.MethodString, []any{f.Size}
		}
		return contractsschema.MethodString, nil
	case schema.Time:
		if f.Precision > 0 {
			return contractsschema.MethodTimestampTz, []any{f.Precision}
		}
		return contractsschema.MethodTimestampTz, nil
	case schema.Bytes:
		return contractsschema.MethodBinary, nil
	}

	// String-based Type Inference (Enums, Custom types)
	sType := strings.ToLower(string(f.DataType))

	if strings.HasPrefix(sType, "enum") {
		return contractsschema.MethodEnum, []any{parseEnum(string(f.DataType))}
	}

	// Helper to check prefixes fast
	for _, p := range stringTypePrefixes {
		if strings.HasPrefix(sType, p) {
			if size := parseTypeSize(sType); size > 0 {
				return contractsschema.MethodString, []any{size}
			}
			return contractsschema.MethodString, nil
		}
	}

	if strings.HasPrefix(sType, "timestamp") || strings.HasPrefix(sType, "datetime") {
		if strings.Contains(sType, "tz") {
			return contractsschema.MethodTimestampTz, nil
		}
		return contractsschema.MethodTimestamp, nil
	}

	// Map lookup for fixed types (json, uuid, etc)
	for k, v := range dataTypeMapping {
		if strings.Contains(sType, k) {
			return v, nil
		}
	}

	// Fallback to Go type name
	goType := strings.ToLower(f.FieldType.String())
	if strings.Contains(goType, "json") {
		return contractsschema.MethodJson, nil
	}
	if strings.Contains(goType, "uuid") {
		return contractsschema.MethodUuid, nil
	}
	if strings.Contains(goType, "ulid") {
		return contractsschema.MethodUlid, nil
	}

	return contractsschema.MethodText, nil
}

// renderIndexes generates index and constraint definitions for a schema.
// It processes composite primary keys, unique indexes, regular indexes, and fulltext indexes,
// returning Blueprint method call strings for each index definition.
func renderIndexes(sch *schema.Schema, fields []*schema.Field) []string {
	var lines []string
	seen := make(map[string]bool)

	add := func(key, line string) {
		if !seen[key] {
			seen[key] = true
			lines = append(lines, line)
		}
	}

	// Composite Primary Keys (if > 1 PK field)
	if len(sch.PrimaryFields) > 1 {
		cols := getColNames(sch.PrimaryFields)
		add("PK:"+strings.Join(cols, ","), formatIndex(contractsschema.IndexMethodPrimary, cols, nil))
	}

	indexes := sch.ParseIndexes()
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].Name < indexes[j].Name })

	for _, idx := range indexes {
		if len(idx.Fields) == 0 {
			continue
		}

		cols := make([]string, 0, len(idx.Fields))
		for _, f := range idx.Fields {
			if f.DBName != "" {
				cols = append(cols, f.DBName)
			}
		}
		if len(cols) == 0 {
			continue
		}

		method := contractsschema.IndexMethodIndex
		switch idx.Class {
		case contractsschema.IndexClassUnique:
			method = contractsschema.IndexMethodUnique
		case contractsschema.IndexClassFullText:
			method = contractsschema.IndexMethodFullText
		case contractsschema.IndexClassPrimary:
			method = contractsschema.IndexMethodPrimary
		}

		add(idx.Class+":"+strings.Join(cols, ","), formatIndex(method, cols, idx))
	}

	// Field Unique Constraints
	for _, f := range fields {
		if f.Unique && f.DBName != "" {
			add("UQ:"+f.DBName, formatIndex(contractsschema.IndexMethodUnique, []string{f.DBName}, nil))
		}
	}

	return lines
}

func intMethod(size int, unsigned bool) string {
	switch {
	case size <= 8:
		if unsigned {
			return contractsschema.MethodUnsignedTinyInteger
		}
		return contractsschema.MethodTinyInteger
	case size <= 16:
		if unsigned {
			return contractsschema.MethodUnsignedSmallInteger
		}
		return contractsschema.MethodSmallInteger
	case size <= 32:
		if unsigned {
			return contractsschema.MethodUnsignedInteger
		}
		return contractsschema.MethodInteger
	default:
		if unsigned {
			return contractsschema.MethodUnsignedBigInteger
		}
		return contractsschema.MethodBigInteger
	}
}

func formatIndex(method string, cols []string, idx *schema.Index) string {
	b := &atom{Builder: strings.Builder{}}
	b.WriteString(tablePrefix)

	args := make([]any, len(cols))
	for i, v := range cols {
		args[i] = v
	}
	b.WriteMethod(method, args...)

	if idx != nil {
		if idx.Type != "" {
			b.WriteMethod(contractsschema.IndexMethodAlgorithm, idx.Type)
		}
		if idx.Name != "" {
			b.WriteMethod(contractsschema.IndexMethodName, idx.Name)
		}
	}
	return b.String()
}

func parseEnum(def string) []any {
	start, end := strings.IndexByte(def, '('), strings.LastIndexByte(def, ')')
	if start == -1 || end <= start+1 {
		return nil
	}

	var values []any
	var buf strings.Builder
	inQuote := false

	for _, r := range def[start+1 : end] {
		if r == '\'' {
			inQuote = !inQuote
			continue
		}
		if r == ',' && !inQuote {
			values = append(values, parseVal(buf.String()))
			buf.Reset()
			continue
		}
		buf.WriteRune(r)
	}
	if buf.Len() > 0 {
		values = append(values, parseVal(buf.String()))
	}
	return values
}

func parseVal(s string) any {
	s = strings.TrimSpace(s)
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if strings.Contains(s, ".") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return s
}

func trimQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '\'' || s[0] == '"') && s[0] == s[len(s)-1] {
		return s[1 : len(s)-1]
	}
	return s
}

func parseTypeSize(s string) int {
	start, end := strings.IndexByte(s, '('), strings.IndexByte(s, ')')
	if start > -1 && end > start {
		// Just grab the first number "varchar(255)" -> 255
		parts := strings.SplitN(s[start+1:end], ",", 2)
		if n, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			return n
		}
	}
	return 0
}

func getColNames(fields []*schema.Field) []string {
	names := make([]string, 0, len(fields))
	for _, f := range fields {
		if f.DBName != "" {
			names = append(names, f.DBName)
		}
	}
	return names
}

// isSQLNullType checks if the type is a known nullable database type.
// Supported types:
//   - database/sql.Null* (e.g., sql.NullString, sql.NullInt64, sql.NullTime)
//   - gorm.DeletedAt
//
// For other types, use pointer types or GORM tags to define nullability.
func isSQLNullType(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	// Check for database/sql.Null* types
	if t.PkgPath() == "database/sql" && strings.HasPrefix(t.Name(), "Null") {
		return true
	}

	// Check for gorm.DeletedAt
	if strings.HasSuffix(t.PkgPath(), "gorm.io/gorm") && t.Name() == "DeletedAt" {
		return true
	}

	return false
}

// atom is a string builder wrapper for generating Blueprint method chains.
type atom struct {
	strings.Builder
}

// WriteMethod appends a method call with arguments to the builder.
// Automatically adds a dot prefix if needed for method chaining.
func (r *atom) WriteMethod(name string, args ...any) {
	if r.Len() > 0 && r.String()[r.Len()-1] != '.' {
		r.WriteByte('.')
	}
	r.WriteString(name)
	r.WriteByte('(')
	for i, arg := range args {
		if i > 0 {
			r.WriteString(", ")
		}
		r.WriteValue(arg)
	}
	r.WriteByte(')')
}

// WriteValue formats and writes a single value to the builder.
// Handles strings, integers, floats, booleans, nil, and slices.
func (r *atom) WriteValue(v any) {
	switch val := v.(type) {
	case nil:
		r.WriteString("nil")
	case string:
		r.WriteString(strconv.Quote(val))
	case int:
		r.WriteString(strconv.Itoa(val))
	case int64:
		r.WriteString(strconv.FormatInt(val, 10))
	case float64:
		r.WriteString(strconv.FormatFloat(val, 'f', -1, 64))
	case bool:
		r.WriteString(strconv.FormatBool(val))
	case []any:
		r.WriteString("[]any{")
		for i, item := range val {
			if i > 0 {
				r.WriteString(", ")
			}
			r.WriteValue(item)
		}
		r.WriteByte('}')
	default:
		_, _ = fmt.Fprintf(r, "%v", val)
	}
}
