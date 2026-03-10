package database

import (
	"reflect"
	"strings"
)

func GetID(dest any) any {
	if dest == nil {
		return nil
	}

	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)

	return GetIDByReflect(t, v)
}

func GetIDByReflect(t reflect.Type, v reflect.Value) any {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		if strings.Contains(t.Field(i).Tag.Get("gorm"), "primaryKey") {
			if v.Field(i).IsZero() {
				return nil
			}

			return v.Field(i).Interface()
		}
		if t.Field(i).Type.Kind() == reflect.Struct && t.Field(i).Anonymous {
			id := GetIDByReflect(t.Field(i).Type, v.Field(i))
			if id != nil {
				return id
			}
		}
	}

	return nil
}
