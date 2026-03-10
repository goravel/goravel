package reflects

import "reflect"

// BKind base data kind type, alias of reflect.Kind
//
// Diff with reflect.Kind:
//   - Int contains all intX types
//   - Uint contains all uintX types
//   - Float contains all floatX types
//   - Array for array and slice types
//   - Complex contains all complexX types
type BKind = reflect.Kind

// base kinds
const (
	// Int for all intX types
	Int = reflect.Int
	// Uint for all uintX types
	Uint = reflect.Uint
	// Float for all floatX types
	Float = reflect.Float32
	// Array for array,slice types
	Array = reflect.Array
	// Complex for all complexX types
	Complex = reflect.Complex64
)

// ToBaseKind convert reflect.Kind to base kind
func ToBaseKind(kind reflect.Kind) BKind {
	return ToBKind(kind)
}

// ToBKind convert reflect.Kind to base kind
func ToBKind(kind reflect.Kind) BKind {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Uint
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Complex64, reflect.Complex128:
		return Complex
	case reflect.Array, reflect.Slice:
		return Array
	default:
		// like: string, map, struct, ptr, func, interface ...
		return kind
	}
}

// Type struct
type Type interface {
	reflect.Type
	// BaseKind value
	BaseKind() BKind
	// RealType returns a ptr type's real type. otherwise, will return self.
	RealType() reflect.Type
	// SafeElem returns a type's element type. otherwise, will return self.
	SafeElem() reflect.Type
}

type xType struct {
	reflect.Type
	baseKind BKind
}

// TypeOf value
func TypeOf(v any) Type {
	rftTyp := reflect.TypeOf(v)

	return &xType{
		Type:     rftTyp,
		baseKind: ToBKind(rftTyp.Kind()),
	}
}

// BaseKind value
func (t *xType) BaseKind() BKind {
	return t.baseKind
}

// RealType returns a ptr type's real type. otherwise, will return self.
func (t *xType) RealType() reflect.Type {
	return TypeReal(t.Type)
}

// SafeElem returns the array, slice, chan, map type's element type. otherwise, will return self.
func (t *xType) SafeElem() reflect.Type {
	return TypeElem(t.Type)
}
