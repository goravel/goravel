package schema

import (
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/support/convert"
)

type ColumnDefinition struct {
	def                any
	onUpdate           any
	autoIncrement      *bool
	comment            *string
	generatedAs        *string
	length             *int
	name               *string
	nullable           *bool
	places             *int
	precision          *int
	total              *int
	ttype              *string
	unsigned           *bool
	useCurrent         *bool
	useCurrentOnUpdate *bool
	after              string
	allowed            []any
	always             bool
	change             bool
	first              bool
}

func NewColumnDefinition(name string, ttype string) driver.ColumnDefinition {
	return &ColumnDefinition{
		name:  &name,
		ttype: convert.Pointer(ttype),
	}
}

func (r *ColumnDefinition) After(column string) driver.ColumnDefinition {
	r.after = column

	return r
}

func (r *ColumnDefinition) Always() driver.ColumnDefinition {
	r.always = true

	return r
}

func (r *ColumnDefinition) AutoIncrement() driver.ColumnDefinition {
	r.autoIncrement = convert.Pointer(true)

	return r
}

func (r *ColumnDefinition) Change() driver.ColumnDefinition {
	r.change = true

	return r
}

func (r *ColumnDefinition) Comment(comment string) driver.ColumnDefinition {
	r.comment = &comment

	return r
}

func (r *ColumnDefinition) Default(def any) driver.ColumnDefinition {
	r.def = def

	return r
}

func (r *ColumnDefinition) First() driver.ColumnDefinition {
	r.first = true

	return r
}

func (r *ColumnDefinition) GeneratedAs(expression ...string) driver.ColumnDefinition {
	expression = append(expression, "")
	r.generatedAs = &expression[0]

	return r
}

func (r *ColumnDefinition) GetAfter() string {
	return r.after
}

func (r *ColumnDefinition) GetAllowed() []any {
	return r.allowed
}

func (r *ColumnDefinition) GetAutoIncrement() bool {
	if r.autoIncrement != nil {
		return *r.autoIncrement
	}

	return false
}

func (r *ColumnDefinition) GetComment() string {
	if r.comment != nil {
		return *r.comment
	}

	return ""
}

func (r *ColumnDefinition) GetDefault() any {
	return r.def
}

func (r *ColumnDefinition) GetGeneratedAs() string {
	if r.generatedAs != nil {
		return *r.generatedAs
	}

	return ""
}

func (r *ColumnDefinition) GetName() string {
	if r.name != nil {
		return *r.name
	}

	return ""
}

func (r *ColumnDefinition) GetLength() int {
	if r.length != nil {
		return *r.length
	}

	return 0
}

func (r *ColumnDefinition) GetNullable() bool {
	if r.nullable != nil {
		return *r.nullable
	}

	return false
}

func (r *ColumnDefinition) GetOnUpdate() any {
	return r.onUpdate
}

func (r *ColumnDefinition) GetPlaces() int {
	if r.places != nil {
		return *r.places
	}

	return 2
}

func (r *ColumnDefinition) GetPrecision() int {
	if r.precision != nil {
		return *r.precision
	}

	return 0
}

func (r *ColumnDefinition) GetTotal() int {
	if r.total != nil {
		return *r.total
	}

	return 8
}

func (r *ColumnDefinition) GetType() string {
	if r.ttype != nil {
		return *r.ttype
	}

	return ""
}

func (r *ColumnDefinition) GetUnsigned() bool {
	if r.unsigned != nil {
		return *r.unsigned
	}

	return false
}

func (r *ColumnDefinition) GetUseCurrent() bool {
	if r.useCurrent != nil {
		return *r.useCurrent
	}

	return false
}

func (r *ColumnDefinition) GetUseCurrentOnUpdate() bool {
	if r.useCurrentOnUpdate != nil {
		return *r.useCurrentOnUpdate
	}

	return false
}

func (r *ColumnDefinition) IsAlways() bool {
	return r.always
}

func (r *ColumnDefinition) IsChange() bool {
	return r.change
}

func (r *ColumnDefinition) IsFirst() bool {
	return r.first
}

func (r *ColumnDefinition) IsSetComment() bool {
	return r != nil && r.comment != nil
}

func (r *ColumnDefinition) IsSetGeneratedAs() bool {
	return r != nil && r.generatedAs != nil
}

func (r *ColumnDefinition) Nullable() driver.ColumnDefinition {
	r.nullable = convert.Pointer(true)

	return r
}

func (r *ColumnDefinition) OnUpdate(value any) driver.ColumnDefinition {
	r.onUpdate = value

	return r
}

func (r *ColumnDefinition) Places(places int) driver.ColumnDefinition {
	r.places = convert.Pointer(places)

	return r
}

func (r *ColumnDefinition) Total(total int) driver.ColumnDefinition {
	r.total = convert.Pointer(total)

	return r
}

func (r *ColumnDefinition) Unsigned() driver.ColumnDefinition {
	r.unsigned = convert.Pointer(true)

	return r
}

func (r *ColumnDefinition) UseCurrent() driver.ColumnDefinition {
	r.useCurrent = convert.Pointer(true)

	return r
}

func (r *ColumnDefinition) UseCurrentOnUpdate() driver.ColumnDefinition {
	r.useCurrentOnUpdate = convert.Pointer(true)

	return r
}
