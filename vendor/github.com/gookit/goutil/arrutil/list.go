package arrutil

import (
	"sort"
	"strings"

	"github.com/gookit/goutil/comdef"
)

// Ints type
type Ints[T comdef.Integer] []T

// String to string
func (is Ints[T]) String() string {
	return ToString(is)
}

// Has given element
func (is Ints[T]) Has(i T) bool {
	for _, iv := range is {
		if i == iv {
			return true
		}
	}
	return false
}

// First element value.
func (is Ints[T]) First(defVal ...T) T {
	if len(is) > 0 {
		return is[0]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty integer slice")
}

// Last element value.
func (is Ints[T]) Last(defVal ...T) T {
	if len(is) > 0 {
		return is[len(is)-1]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty integer slice")
}

// Sort the int slice
func (is Ints[T]) Sort() {
	sort.Sort(is)
}

// Len get length
func (is Ints[T]) Len() int {
	return len(is)
}

// Less compare two elements
func (is Ints[T]) Less(i, j int) bool {
	return is[i] < is[j]
}

// Swap elements by indexes
func (is Ints[T]) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

// Strings type
type Strings []string

// String to string
func (ss Strings) String() string {
	return strings.Join(ss, ",")
}

// Join to string
func (ss Strings) Join(sep string) string {
	return strings.Join(ss, sep)
}

// Has given element
func (ss Strings) Has(sub string) bool {
	return ss.Contains(sub)
}

// Contains given element
func (ss Strings) Contains(sub string) bool {
	for _, s := range ss {
		if s == sub {
			return true
		}
	}
	return false
}

// First element value.
func (ss Strings) First(defVal ...string) string {
	if len(ss) > 0 {
		return ss[0]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty string list")
}

// Last element value.
func (ss Strings) Last(defVal ...string) string {
	if len(ss) > 0 {
		return ss[len(ss)-1]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty string list")
}

// Sort the string slice
func (ss Strings) Sort() {
	sort.Strings(ss)
}

// SortedList definition for compared type
type SortedList[T comdef.Compared] []T

// Len get length
func (ls SortedList[T]) Len() int {
	return len(ls)
}

// Less compare two elements
func (ls SortedList[T]) Less(i, j int) bool {
	return ls[i] < ls[j]
}

// Swap elements by indexes
func (ls SortedList[T]) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

// IsEmpty check
func (ls SortedList[T]) IsEmpty() bool {
	return len(ls) == 0
}

// String to string
func (ls SortedList[T]) String() string {
	return ToString(ls)
}

// Has given element
func (ls SortedList[T]) Has(el T) bool {
	return ls.Contains(el)
}

// Contains given element
func (ls SortedList[T]) Contains(el T) bool {
	for _, v := range ls {
		if v == el {
			return true
		}
	}
	return false
}

// First element value.
func (ls SortedList[T]) First(defVal ...T) T {
	if len(ls) > 0 {
		return ls[0]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty list")
}

// Last element value.
func (ls SortedList[T]) Last(defVal ...T) T {
	if ln := len(ls); ln > 0 {
		return ls[ln-1]
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	panic("empty list")
}

// Remove given element
func (ls SortedList[T]) Remove(el T) SortedList[T] {
	return Filter(ls, func(v T) bool {
		return v != el
	})
}

// Filter the slice, default will filter zero value.
func (ls SortedList[T]) Filter(filter ...comdef.MatchFunc[T]) SortedList[T] {
	return Filter(ls, filter...)
}

// Map the slice to new slice. TODO syntax ERROR: Method cannot have type parameters
// func (ls SortedList[T]) Map[V any](mapFn MapFn[T, V]) SortedList[V] {
// 	return Map(ls, mapFn)
// }

// Sort the slice
func (ls SortedList[T]) Sort() {
	sort.Sort(ls)
}
