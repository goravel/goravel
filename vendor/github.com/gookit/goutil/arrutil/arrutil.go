// Package arrutil provides some util functions for array, slice
package arrutil

import (
	"github.com/gookit/goutil/mathutil"
)

// GetRandomOne get random element from an array/slice
func GetRandomOne[T any](arr []T) T { return RandomOne(arr) }

// RandomOne get random element from an array/slice
func RandomOne[T any](arr []T) T {
	if ln := len(arr); ln > 0 {
		i := mathutil.RandomInt(0, len(arr))
		return arr[i]
	}
	panic("cannot get value from nil or empty slice")
}
