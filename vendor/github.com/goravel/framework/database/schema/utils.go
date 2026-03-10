package schema

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/database/driver"
)

type Expression string

func ColumnDefaultValue(def any) string {
	switch value := def.(type) {
	case bool:
		return "'" + cast.ToString(cast.ToInt(value)) + "'"
	case Expression:
		return string(value)
	default:
		return "'" + cast.ToString(def) + "'"
	}
}

func ColumnType(grammar driver.Grammar, column driver.ColumnDefinition) string {
	t := []rune(column.GetType())
	if len(t) == 0 {
		return ""
	}

	t[0] = unicode.ToUpper(t[0])
	methodName := fmt.Sprintf("Type%s", string(t))
	methodValue := reflect.ValueOf(grammar).MethodByName(methodName)
	if methodValue.IsValid() {
		args := []reflect.Value{reflect.ValueOf(column)}
		callResult := methodValue.Call(args)

		return callResult[0].String()
	}

	return column.GetType()
}
