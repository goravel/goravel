package gorm

import (
	"context"
	"maps"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/support/str"
)

type Event struct {
	columnNames map[string]string
	dest        any
	destOfMap   map[string]any
	model       any
	modelOfMap  map[string]any
	query       *Query
}

func NewEvent(query *Query, model, dest any) *Event {
	return &Event{
		dest:  dest,
		model: model,
		query: query,
	}
}

func (e *Event) Context() context.Context {
	return e.query.ctx
}

func (e *Event) GetAttribute(key string) any {
	destOfMap := e.getDestOfMap()
	value, exist := destOfMap[e.toDBColumnName(key)]
	if exist && e.validColumn(key) && e.validValue(key, value) {
		return value
	}

	return e.GetOriginal(key)
}

func (e *Event) GetOriginal(key string, def ...any) any {
	modelOfMap := e.getModelOfMap()
	value, exist := modelOfMap[e.toDBColumnName(key)]
	if exist {
		return value
	}

	if len(def) > 0 {
		return def[0]
	}

	return nil
}

func (e *Event) IsClean(fields ...string) bool {
	return !e.IsDirty(fields...)
}

func (e *Event) IsDirty(columns ...string) bool {
	destOfMap := e.getDestOfMap()

	if len(columns) == 0 {
		for destColumn, destValue := range destOfMap {
			if !e.validColumn(destColumn) || !e.validValue(destColumn, destValue) {
				continue
			}
			if e.dirty(destColumn, destValue) {
				return true
			}
		}
	} else {
		for _, column := range columns {
			if !e.validColumn(column) {
				continue
			}
			for destColumn, destValue := range destOfMap {
				if !e.validColumn(destColumn) || !e.validValue(destColumn, destValue) {
					continue
				}
				if e.equalColumnName(column, destColumn) && e.dirty(destColumn, destValue) {
					return true
				}
			}
		}
	}

	return false
}

func (e *Event) Query() orm.Query {
	return NewQuery(e.query.ctx, e.query.config, e.query.dbConfig, e.query.instance.Session(&gorm.Session{NewDB: true}), e.query.grammar, e.query.log, e.query.modelToObserver, nil)
}

func (e *Event) SetAttribute(key string, value any) {
	if e.dest == nil {
		return
	}

	destOfMap := e.getDestOfMap()
	if destOfMap == nil {
		return
	}

	destOfMap[e.toDBColumnName(key)] = value
	e.destOfMap = destOfMap

	if m, ok := e.dest.(map[string]any); ok {
		m[key] = value
	} else {
		destType := reflect.TypeOf(e.dest)
		destValue := reflect.ValueOf(e.dest)
		if destType.Kind() == reflect.Pointer {
			destType = destType.Elem()
			destValue = destValue.Elem()
		}

		if !destValue.CanAddr() {
			destValueCanAddr := reflect.New(destValue.Type())
			destValueCanAddr.Elem().Set(destValue)
			e.dest = destValueCanAddr.Interface()
			e.query.instance.Statement.Dest = e.dest
			destValue = destValueCanAddr.Elem()
		}

		for i := 0; i < destType.NumField(); i++ {
			if !destType.Field(i).IsExported() {
				continue
			}
			if e.equalColumnName(destType.Field(i).Name, key) {
				if value == nil {
					destValue.Field(i).Set(reflect.Zero(destValue.Field(i).Type()))
				} else {
					valueValue := reflect.ValueOf(value)
					destValue.Field(i).Set(valueValue)
				}
			}
		}
	}
}

func (e *Event) dirty(destColumn string, destValue any) bool {
	modelOfMap := e.getModelOfMap()
	dbDestColumn := e.toDBColumnName(destColumn)

	if modelValue, exist := modelOfMap[dbDestColumn]; exist {
		return !reflect.DeepEqual(modelValue, destValue)
	}

	return true
}

func (e *Event) equalColumnName(origin, source string) bool {
	originDbColumnName := e.toDBColumnName(origin)
	sourceDbColumnName := e.toDBColumnName(source)

	if originDbColumnName == "" || sourceDbColumnName == "" {
		return false
	}

	return originDbColumnName == sourceDbColumnName
}

func (e *Event) getColumnNames() map[string]string {
	if e.columnNames == nil {
		if e.model != nil {
			e.columnNames = fetchColumnNames(e.model)
		} else {
			e.columnNames = fetchColumnNames(e.dest)
		}
	}

	return e.columnNames
}

func (e *Event) getDestOfMap() map[string]any {
	if e.dest == nil {
		return nil
	}
	if e.destOfMap != nil {
		return e.destOfMap
	}

	destOfMap := make(map[string]any)
	if destMap, ok := e.dest.(map[string]any); ok {
		for key, value := range destMap {
			destOfMap[key] = value
			destOfMap[str.Of(key).Snake().String()] = value
		}
	} else {
		destType := reflect.TypeOf(e.dest)
		if destType.Kind() == reflect.Pointer {
			destType = destType.Elem()
		}
		if destType.Kind() == reflect.Struct {
			destOfMap = structToMap(e.dest)
		}
	}

	e.destOfMap = destOfMap

	return e.destOfMap
}

func (e *Event) getModelOfMap() map[string]any {
	if e.modelOfMap != nil {
		return e.modelOfMap
	}

	if e.model == nil {
		return map[string]any{}
	}

	e.modelOfMap = structToMap(e.model)

	return e.modelOfMap
}

func (e *Event) toDBColumnName(name string) string {
	dbColumnName, exist := e.getColumnNames()[name]
	if exist {
		return dbColumnName
	}

	return ""
}

func (e *Event) validColumn(name string) bool {
	dbColumn := e.toDBColumnName(name)
	if dbColumn == "" {
		return false
	}

	selectColumns := e.query.instance.Statement.Selects
	omitColumns := e.query.instance.Statement.Omits
	if len(selectColumns) > 0 {
		for _, selectColumn := range selectColumns {
			dbSelectColumn := e.toDBColumnName(selectColumn)
			if dbSelectColumn == "" {
				continue
			}

			if dbSelectColumn == dbColumn {
				return true
			}
		}

		return false
	}
	if len(omitColumns) > 0 {
		for _, omitColumn := range omitColumns {
			dbOmitColumn := e.toDBColumnName(omitColumn)
			if dbOmitColumn == "" {
				continue
			}

			if dbOmitColumn == dbColumn {
				return false
			}
		}

		return true
	}

	return true
}

func (e *Event) validValue(name string, value any) bool {
	dbColumn := e.toDBColumnName(name)
	if dbColumn == "" {
		return false
	}

	selectColumns := e.query.instance.Statement.Selects
	if len(selectColumns) > 0 {
		return e.validColumn(name)
	}

	if value == nil {
		return false
	}

	valueValue := reflect.ValueOf(value)

	return !valueValue.IsZero()
}

func fetchColumnNames(model any) map[string]string {
	res := make(map[string]string)
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return res
	}

	for i := 0; i < modelType.NumField(); i++ {
		if !modelType.Field(i).IsExported() {
			continue
		}
		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)
		if fieldValue.Kind() == reflect.Struct && fieldType.Anonymous {
			subStructMap := fetchColumnNames(fieldValue.Interface())
			maps.Copy(res, subStructMap)
			continue
		}

		dbColumn := structNameToDbColumnName(modelType.Field(i).Name, modelType.Field(i).Tag.Get("gorm"))
		res[modelType.Field(i).Name] = dbColumn
		res[dbColumn] = dbColumn
	}

	return res
}

func structToMap(data any) map[string]any {
	res := make(map[string]any)
	modelType := reflect.TypeOf(data)
	modelValue := reflect.ValueOf(data)

	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return res
	}

	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		dbColumn := structNameToDbColumnName(fieldType.Name, fieldType.Tag.Get("gorm"))
		if fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				res[dbColumn] = fieldValue.Interface()
				continue
			}
		}

		if (fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Pointer) && fieldType.Anonymous {
			subStructMap := structToMap(fieldValue.Interface())
			maps.Copy(res, subStructMap)
		} else {
			res[dbColumn] = fieldValue.Interface()
		}
	}

	return res
}

func structNameToDbColumnName(structName, tag string) string {
	if strings.Contains(tag, "column:") {
		tags := strings.Split(tag, ";")
		for _, item := range tags {
			if strings.Contains(item, "column:") {
				return strings.Trim(strings.ReplaceAll(item, "column:", ""), " ")
			}
		}
	}

	return str.Of(structName).Snake().String()
}
