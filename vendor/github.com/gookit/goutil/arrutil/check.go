package arrutil

import (
	"reflect"
	"strings"

	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/mathutil"
)

// SliceHas check the slice contains the given value
func SliceHas[T comdef.ScalarType](slice []T, val T) bool {
	for _, ele := range slice {
		if ele == val {
			return true
		}
	}
	return false
}

// IntsHas check the []comdef.Integer contains the given value
func IntsHas[T comdef.Integer](ints []T, val T) bool {
	for _, ele := range ints {
		if ele == val {
			return true
		}
	}
	return false
}

// Int64sHas check the []int64 contains the given value
func Int64sHas(ints []int64, val int64) bool {
	for _, ele := range ints {
		if ele == val {
			return true
		}
	}
	return false
}

// StringsHas check the []string contains the given element
func StringsHas[T ~string](ss []T, val T) bool {
	for _, ele := range ss {
		if ele == val {
			return true
		}
	}
	return false
}

// InStrings check elem in the ss. alias of StringsHas()
func InStrings[T ~string](elem T, ss []T) bool {
	return StringsHas(ss, elem)
}

// NotIn check the given value whether not in the list
func NotIn[T comdef.ScalarType](value T, list []T) bool {
	return !In(value, list)
}

// In check the given value whether in the list
func In[T comdef.ScalarType](value T, list []T) bool {
	for _, elem := range list {
		if elem == value {
			return true
		}
	}
	return false
}

// ContainsAll check given values is sub-list of sample list.
func ContainsAll[T comdef.ScalarType](list, values []T) bool {
	return IsSubList(values, list)
}

// IsSubList check given values is sub-list of sample list.
func IsSubList[T comdef.ScalarType](values, list []T) bool {
	for _, value := range values {
		if !In(value, list) {
			return false
		}
	}
	return true
}

// IsParent check given values is parent-list of samples.
func IsParent[T comdef.ScalarType](values, list []T) bool {
	return IsSubList(list, values)
}

// HasValue check array(strings, intXs, uintXs) should be contained the given value(int(X),string).
func HasValue(arr, val any) bool { return Contains(arr, val) }

// Contains check slice/array(strings, intXs, uintXs) should be contained the given value(int(X),string).
//
// TIP: Difference the In(), Contains() will try to convert value type,
// and Contains() support array type.
func Contains(arr, val any) bool {
	if val == nil || arr == nil {
		return false
	}

	// if is string value
	if strVal, ok := val.(string); ok {
		if ss, ok := arr.([]string); ok {
			return StringsHas(ss, strVal)
		}

		rv := reflect.ValueOf(arr)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			for i := 0; i < rv.Len(); i++ {
				if v, ok := rv.Index(i).Interface().(string); ok && strings.EqualFold(v, strVal) {
					return true
				}
			}
		}

		return false
	}

	// as int value
	intVal, err := mathutil.Int64(val)
	if err != nil {
		return false
	}

	if int64s, err := ToInt64s(arr); err == nil {
		return Int64sHas(int64s, intVal)
	}
	return false
}

// NotContains check array(strings, ints, uints) should be not contains the given value.
func NotContains(arr, val any) bool {
	return !Contains(arr, val)
}
