package db

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"gorm.io/gorm"

	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/str"
)

type Row struct {
	err error
	row map[string]any
}

func NewRow(row map[string]any, err error) *Row {
	return &Row{row: row, err: err}
}

func (r *Row) Err() error {
	return r.err
}

func (r *Row) Scan(value any) error {
	if r.err != nil {
		return r.err
	}

	msConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToStringHookFunc(), ToTimeHookFunc(), ToDeletedAtHookFunc(), ToScannerHookFunc(), ToSliceHookFunc(), ToMapHookFunc(),
		),
		Squash: true,
		Result: value,
		MatchName: func(mapKey, fieldName string) bool {
			return str.Of(mapKey).Studly().String() == fieldName || mapKey == str.Of(fieldName).Snake().String() || strings.EqualFold(mapKey, fieldName)
		},
	}

	decoder, err := mapstructure.NewDecoder(msConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(r.row)
}

// ToStringHookFunc is a hook function that converts []uint8 to string.
// Mysql returns []uint8 for String type when scanning the rows.
func ToStringHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf("") {
			return data, nil
		}

		dataSlice, ok := data.([]uint8)
		if ok {
			return string(dataSlice), nil
		}

		return data, nil
	}
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func ToDeletedAtHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(gorm.DeletedAt{}) {
			return data, nil
		}

		if f == reflect.TypeOf(time.Time{}) {
			return gorm.DeletedAt{Time: data.(time.Time), Valid: true}, nil
		}

		if f.Kind() == reflect.String {
			return gorm.DeletedAt{Time: carbon.Parse(data.(string)).StdTime(), Valid: true}, nil
		}

		return data, nil
	}
}

// ToScannerHookFunc is a hook function that handles types with custom Scan methods (sql.Scanner interface).
// This includes carbon types and other custom types implementing the Scan method.
func ToScannerHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		// Skip types that are handled by other specific hooks
		if t == reflect.TypeOf(time.Time{}) || t == reflect.TypeOf(gorm.DeletedAt{}) {
			return data, nil
		}

		// Skip if source and target are the same type
		if f == t {
			return data, nil
		}

		// Only process database types (string, []byte, []uint8, time.Time)
		if f.Kind() != reflect.String && f != reflect.TypeOf([]byte(nil)) && f != reflect.TypeOf([]uint8(nil)) && f != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// Check if the target type implements a Scan method
		scannerType := reflect.TypeOf((*interface{ Scan(any) error })(nil)).Elem()

		// Create a pointer to the target type to check for Scan method
		targetPtr := reflect.PointerTo(t)
		if !targetPtr.Implements(scannerType) {
			return data, nil
		}

		// Handle nil or empty data
		if data == nil {
			return reflect.Zero(t).Interface(), nil
		}
		if str, ok := data.(string); ok && str == "" {
			return reflect.Zero(t).Interface(), nil
		}

		// Create a new instance of the target type
		result := reflect.New(t)
		scanner := result.Interface().(interface{ Scan(any) error })

		// Convert string to []byte if needed (common for JSON fields from database)
		scanData := data
		if str, ok := data.(string); ok {
			scanData = []byte(str)
		}

		// Call the Scan method with the data
		if err := scanner.Scan(scanData); err != nil {
			return nil, err
		}

		return result.Elem().Interface(), nil
	}
}

// ToSliceHookFunc is a hook function that converts JSON string to slice.
func ToSliceHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t.Kind() != reflect.Slice || f.Kind() != reflect.String {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		// Return empty slice for empty string
		if str == "" {
			return reflect.MakeSlice(t, 0, 0).Interface(), nil
		}

		result := reflect.New(t).Interface()
		if err := json.Unmarshal([]byte(str), result); err != nil {
			return nil, err
		}

		return reflect.ValueOf(result).Elem().Interface(), nil
	}
}

// ToMapHookFunc is a hook function that converts JSON string to map.
func ToMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t.Kind() != reflect.Map || f.Kind() != reflect.String {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		// Return empty map for empty string
		if str == "" {
			return reflect.MakeMap(t).Interface(), nil
		}

		result := reflect.MakeMap(t).Interface()
		if err := json.Unmarshal([]byte(str), &result); err != nil {
			return nil, err
		}

		return result, nil
	}
}
