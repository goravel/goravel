package reflects

import (
	"fmt"
	"reflect"
)

// MakeSliceByElem create a new slice by the element type.
//
// - elType: the type of the element.
// - returns: the new slice.
//
// Usage:
//
//	sl := MakeSliceByElem(reflect.TypeOf(1), 10, 20)
//	sl.Index(0).SetInt(10)
//
//	// Or use reflect.AppendSlice() merge two slice
//	// Or use `for` with `reflect.Append()` add elements
func MakeSliceByElem(elTyp reflect.Type, len, cap int) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(elTyp), len, cap)
}

// FlatSlice flatten multi-level slice to given depth-level slice.
//
// Example:
//
//	FlatSlice([]any{ []any{3, 4}, []any{5, 6} }, 1) // Output: []any{3, 4, 5, 6}
//
// always return reflect.Value of []any. note: maybe flatSl.Cap != flatSl.Len
func FlatSlice(sl reflect.Value, depth int) reflect.Value {
	items := make([]reflect.Value, 0, sl.Cap())
	slCap := addSliceItem(sl, depth, func(item reflect.Value) {
		items = append(items, item)
	})

	flatSl := reflect.MakeSlice(reflect.SliceOf(anyType), 0, slCap)
	flatSl = reflect.Append(flatSl, items...)

	return flatSl
}

func addSliceItem(sl reflect.Value, depth int, collector func(item reflect.Value)) (c int) {
	for i := 0; i < sl.Len(); i++ {
		v := Elem(sl.Index(i))

		if depth > 0 {
			if v.Kind() != reflect.Slice {
				panic(fmt.Sprintf("depth: %d, the value of index %d is not slice", depth, i))
			}
			c += addSliceItem(v, depth-1, collector)
		} else {
			collector(v)
		}
	}

	if depth == 0 {
		c = sl.Cap()
	}
	return c
}
