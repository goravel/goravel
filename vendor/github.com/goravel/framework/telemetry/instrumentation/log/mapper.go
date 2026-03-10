// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
//
// Adapted from: https://github.com/open-telemetry/opentelemetry-go-contrib
// Modified by Goravel to support framework-specific log contracts.

package log

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"

	contractslog "github.com/goravel/framework/contracts/log"
)

func toSeverity(level contractslog.Level) otellog.Severity {
	switch level {
	case contractslog.LevelPanic:
		return otellog.SeverityFatal4
	case contractslog.LevelFatal:
		return otellog.SeverityFatal
	case contractslog.LevelError:
		return otellog.SeverityError
	case contractslog.LevelWarning:
		return otellog.SeverityWarn
	case contractslog.LevelInfo:
		return otellog.SeverityInfo
	case contractslog.LevelDebug:
		return otellog.SeverityDebug
	default:
		return otellog.SeverityInfo
	}
}

func toValue(v any) otellog.Value {
	if v == nil {
		return otellog.Value{}
	}

	switch val := v.(type) {
	case otellog.Value:
		return val
	case string:
		return otellog.StringValue(val)
	case bool:
		return otellog.BoolValue(val)
	case int:
		return otellog.Int64Value(int64(val))
	case int64:
		return otellog.Int64Value(val)
	case float64:
		return otellog.Float64Value(val)
	case error:
		return otellog.StringValue(val.Error())
	case []byte:
		return otellog.BytesValue(val)
	case time.Time:
		return otellog.StringValue(val.Format(time.RFC3339Nano))
	case map[string]any:
		return toMapValue(val)
	case []string:
		return toStringSliceValue(val)
	case int32:
		return otellog.Int64Value(int64(val))
	case int16:
		return otellog.Int64Value(int64(val))
	case int8:
		return otellog.Int64Value(int64(val))
	case uint:
		return toUintValue(uint64(val))
	case uint64:
		return toUintValue(val)
	case uint32:
		return otellog.Int64Value(int64(val))
	case uint16:
		return otellog.Int64Value(int64(val))
	case uint8:
		return otellog.Int64Value(int64(val))
	case uintptr:
		return toUintValue(uint64(val))
	case float32:
		return otellog.Float64Value(float64(val))
	case time.Duration:
		return otellog.Int64Value(val.Nanoseconds())
	case complex64:
		r := otellog.Float64("r", float64(real(val)))
		i := otellog.Float64("i", float64(imag(val)))
		return otellog.MapValue(r, i)
	case complex128:
		r := otellog.Float64("r", real(val))
		i := otellog.Float64("i", imag(val))
		return otellog.MapValue(r, i)
	case fmt.Stringer:
		return otellog.StringValue(val.String())
	case attribute.Value:
		return otellog.ValueFromAttribute(val)
	}

	return toReflectedValue(v)
}

func toReflectedValue(v any) otellog.Value {
	t := reflect.TypeOf(v)
	if t == nil {
		return otellog.Value{}
	}
	val := reflect.ValueOf(v)

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		items := make([]otellog.Value, val.Len())
		for i := 0; i < val.Len(); i++ {
			items[i] = toValue(val.Index(i).Interface())
		}
		return otellog.SliceValue(items...)

	case reflect.Map:
		kvs := make([]otellog.KeyValue, 0, val.Len())
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			var keyStr string
			if k.Kind() == reflect.String {
				keyStr = k.String()
			} else {
				keyStr = fmt.Sprintf("%+v", k.Interface())
			}
			kvs = append(kvs, otellog.KeyValue{
				Key:   keyStr,
				Value: toValue(iter.Value().Interface()),
			})
		}
		return otellog.MapValue(kvs...)

	case reflect.Struct:
		return otellog.StringValue(fmt.Sprintf("%+v", v))

	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return otellog.Value{}
		}
		return toValue(val.Elem().Interface())

	case reflect.String:
		return otellog.StringValue(val.String())
	case reflect.Bool:
		return otellog.BoolValue(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return otellog.Int64Value(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return toUintValue(val.Uint())
	case reflect.Float32, reflect.Float64:
		return otellog.Float64Value(val.Float())

	default:
		return otellog.StringValue(fmt.Sprintf("unhandled: (%s) %+v", t, v))
	}
}

func toMapValue(m map[string]any) otellog.Value {
	kvs := make([]otellog.KeyValue, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, otellog.KeyValue{
			Key:   k,
			Value: toValue(v),
		})
	}
	return otellog.MapValue(kvs...)
}

func toStringSliceValue(s []string) otellog.Value {
	items := make([]otellog.Value, len(s))
	for i, v := range s {
		items[i] = otellog.StringValue(v)
	}
	return otellog.SliceValue(items...)
}

func toUintValue(v uint64) otellog.Value {
	if v > math.MaxInt64 {
		return otellog.StringValue(strconv.FormatUint(v, 10))
	}
	return otellog.Int64Value(int64(v))
}
