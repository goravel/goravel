package filter

import (
	"errors"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

// Some alias methods.
var (
	Lower = strings.ToLower
	Upper = strings.ToUpper
	Title = strings.ToTitle

	// EscapeJS escape javascript string
	EscapeJS = template.JSEscapeString
	// EscapeHTML escape html string
	EscapeHTML = template.HTMLEscapeString
	// error for params
	errInvalidParam = errors.New("invalid input parameter")
)

/*************************************************************
 * value to int,bool,float,string
 *************************************************************/

// Int convert string to int
func Int(in any) (int, error) { return ToInt(in) }

// MustInt convert string to int, alias of the mathutil.SafeInt
func MustInt(in any) int {
	return mathutil.SafeInt(in)
}

// ToInt convert string to int
func ToInt(in any) (int, error) { return mathutil.ToInt(in) }

// Int64 convert value to int64
func Int64(in any) (int64, error) { return ToInt64(in) }

// ToInt64 convert value to int64
func ToInt64(val any) (int64, error) { return mathutil.ToInt64(val) }

// MustInt64 convert value to int64, alias of the mathutil.SafeInt64
func MustInt64(in any) int64 {
	return mathutil.SafeInt64(in)
}

// Uint convert string to uint
func Uint(in any) (uint, error) { return ToUint(in) }

// ToUint convert string to uint
func ToUint(in any) (uint, error) { return mathutil.ToUint(in) }

// MustUint convert string to uint, will ignore error
func MustUint(in any) uint {
	return mathutil.SafeUint(in)
}

// Uint64 convert string to uint64
func Uint64(in any) (uint64, error) { return ToUint64(in) }

// ToUint64 convert string to uint64
func ToUint64(in any) (uint64, error) { return mathutil.ToUint64(in) }

// MustUint64 convert string to uint64, alias of the mathutil.SafeUint64
func MustUint64(in any) uint64 {
	return mathutil.SafeUint64(in)
}

// Float convert string to float
func Float(s string) (float64, error) { return ToFloat(s) }

// ToFloat convert string to float
func ToFloat(s string) (float64, error) {
	return strconv.ParseFloat(Trim(s), 0)
}

// MustFloat convert string to float, will ignore error
func MustFloat(s string) float64 {
	val, _ := strconv.ParseFloat(Trim(s), 0)
	return val
}

// ToBool convert string to bool
func ToBool(s string) (bool, error) { return Bool(s) }

// Bool parse string to bool
func Bool(s string) (bool, error) { return strutil.ToBool(s) }

// MustBool convert, will ignore error.
func MustBool(s string) bool {
	val, _ := Bool(Trim(s))
	return val
}

// String convert val to string
func String(val any) (string, error) { return ToString(val) }

// MustString convert value to string, will ignore error
func MustString(in any) string {
	return strutil.SafeString(in)
}

// ToString convert value to string
func ToString(val any) (string, error) { return strutil.ToString(val) }

/*************************************************************
 * change string case
 *************************************************************/

// Lowercase alias of the strings.ToLower()
func Lowercase(s string) string { return strings.ToLower(s) }

// Uppercase alias of the strings.ToUpper()
func Uppercase(s string) string { return strings.ToUpper(s) }

// UpperWord Change the first character of each word to uppercase
func UpperWord(s string) string { return strutil.UpperWord(s) }

// LowerFirst lower first char
func LowerFirst(s string) string { return strutil.LowerFirst(s) }

// UpperFirst upper first char
func UpperFirst(s string) string { return strutil.UpperFirst(s) }

// Snake alias of the SnakeCase
func Snake(s string, sep ...string) string { return SnakeCase(s, sep...) }

// SnakeCase convert. eg "RangePrice" -> "range_price"
func SnakeCase(s string, sep ...string) string { return strutil.SnakeCase(s, sep...) }

// Camel alias of the CamelCase
func Camel(s string, sep ...string) string { return strutil.CamelCase(s, sep...) }

// CamelCase convert string to camel case.
//
// Support:
//
//	"range_price" -> "rangePrice"
//	"range price" -> "rangePrice"
//	"range-price" -> "rangePrice"
func CamelCase(s string, sep ...string) string { return strutil.CamelCase(s, sep...) }

/*************************************************************
 * string to slice, time
 *************************************************************/

// StrToInts split string to slice and convert item to int.
func StrToInts(s string, sep ...string) ([]int, error) { return strutil.ToIntSlice(s, sep...) }

// StrToArray alias of the StrToSlice()
func StrToArray(s string, sep ...string) []string { return StrToSlice(s, sep...) }

// StrToSlice split string to array.
func StrToSlice(s string, sep ...string) []string {
	if len(sep) > 0 {
		return strutil.Split(s, sep[0])
	}
	return strutil.Split(s, ",")
}

// StringsToInts string slice to int slice
func StringsToInts(ss []string) ([]int, error) { return arrutil.StringsToInts(ss) }

// StrToTime convert date string to time.Time
func StrToTime(s string, layouts ...string) (time.Time, error) {
	return strutil.ToTime(s, layouts...)
}
