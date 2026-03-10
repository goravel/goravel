// Package filter provide data filter, sanitize, convert process.
//
// Source code and other details for the project are available at GitHub:
//
//	https://github.com/gookit/filter
//
// More usage please see README and tests
package filter

import (
	"fmt"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

// Apply a filter by name. for filter value.
func Apply(name string, val any, args []string) (any, error) {
	var err error
	realName := Name(name)

	// don't limit value type
	if _, ok := dontLimitType[realName]; ok {
		switch realName {
		case "int":
			val, err = mathutil.ToInt(val)
		case "uint":
			val, err = mathutil.ToUint(val)
		case "int64":
			val, err = mathutil.ToInt64(val)
		case "float":
			val, err = mathutil.ToFloat(val)
		case "unique":
			val = Unique(val)
		case "trimStrings":
			if ss, ok := val.([]string); ok {
				val = arrutil.TrimStrings(ss)
			} else {
				err = errInvalidParam
			}
		case "stringsToInts":
			if ss, ok := val.([]string); ok {
				val, err = arrutil.StringsToInts(ss)
			} else {
				err = errInvalidParam
			}
		}
		return val, err
	}

	// check val is string
	var str string

	// up: support filter pointer string value
	if poStr, ok := val.(*string); ok {
		str = *poStr
	} else if str, ok = val.(string); !ok {
		return nil, fmt.Errorf("filter: '%s' only use for string type, input %T", name, val)
	}

	// val is must be string.
	switch realName {
	case "bool":
		val, err = strutil.ToBool(str)
	case "trim":
		val = strutil.Trim(str, args...)
	case "trimLeft":
		val = strutil.TrimLeft(str, args...)
	case "trimRight":
		val = strutil.TrimRight(str, args...)
	case "title":
		val = Title(str)
	case "email":
		val = strutil.FilterEmail(str)
	case "substr":
		val = strutil.Substr(str, MustInt(args[0]), MustInt(args[1]))
	case "lower":
		val = strutil.Lowercase(str)
	case "upper":
		val = strutil.Uppercase(str)
	case "lowerFirst":
		val = strutil.LowerFirst(str)
	case "upperFirst":
		val = strutil.UpperFirst(str)
	case "upperWord":
		val = strutil.UpperWord(str)
	case "snakeCase":
		val = strutil.SnakeCase(str, args...)
	case "camelCase":
		val = strutil.CamelCase(str, args...)
	case "URLEncode":
		val = strutil.URLEncode(str)
	case "URLDecode":
		val = strutil.URLDecode(str)
	case "escapeJS":
		val = strutil.EscapeJS(str)
	case "escapeHTML":
		val = strutil.EscapeHTML(str)
	case "strToInts":
		val, err = strutil.ToInts(str, args...)
	case "strToSlice":
		val = strutil.ToSlice(str, args...)
	case "strToTime":
		val, err = strutil.ToTime(str, args...)
	}

	return val, err
}

// GetByPath get value from a map[string]any. eg "top" "top.sub"
func GetByPath(key string, mp map[string]any) (any, bool) {
	return maputil.GetByPath(key, mp)
}

func parseArgString(argStr string) (ss []string) {
	if argStr == "" { // no arg
		return
	}

	if len(argStr) == 1 { // one char
		return []string{argStr}
	}
	return strutil.Split(argStr, ",")
}
