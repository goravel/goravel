# Reflects

`reflects` Provide extends reflect util functions.

- some features

## Install

```bash
go get github.com/gookit/goutil/reflects
```

## Go docs

- [Go docs](https://pkg.go.dev/github.com/gookit/goutil/reflects)

## Usage

```go
import "github.com/gookit/goutil/reflects"

// get struct field value
reflects.GetFieldValue(obj, "Name")
```

## Functions API

> **Note**: doc by run `go doc ./reflects`

```go
func BaseTypeVal(v reflect.Value) (value any, err error)
func ConvSlice(oldSlRv reflect.Value, newElemTyp reflect.Type) (rv reflect.Value, err error)
func EachMap(mp reflect.Value, fn func(key, val reflect.Value))
func EachStrAnyMap(mp reflect.Value, fn func(key string, val any))
func Elem(v reflect.Value) reflect.Value
func FlatMap(rv reflect.Value, fn FlatFunc)
func HasChild(v reflect.Value) bool
func Indirect(v reflect.Value) reflect.Value
func IsAnyInt(k reflect.Kind) bool
func IsArrayOrSlice(k reflect.Kind) bool
func IsEmpty(v reflect.Value) bool
func IsEmptyValue(v reflect.Value) bool
func IsEqual(src, dst any) bool
func IsFunc(val any) bool
func IsIntx(k reflect.Kind) bool
func IsNil(v reflect.Value) bool
func IsSimpleKind(k reflect.Kind) bool
func IsUintX(k reflect.Kind) bool
func Len(v reflect.Value) int
func SetRValue(rv, val reflect.Value)
func SetUnexportedValue(rv reflect.Value, value any)
func SetValue(rv reflect.Value, val any) error
func SliceElemKind(typ reflect.Type) reflect.Kind
func SliceSubKind(typ reflect.Type) reflect.Kind
func String(rv reflect.Value) string
func ToString(rv reflect.Value) (str string, err error)
func UnexportedValue(rv reflect.Value) any
func ValToString(rv reflect.Value, defaultAsErr bool) (str string, err error)
func ValueByKind(val any, kind reflect.Kind) (rv reflect.Value, err error)
func ValueByType(val any, typ reflect.Type) (rv reflect.Value, err error)
type BKind uint
    func ToBKind(kind reflect.Kind) BKind
    func ToBaseKind(kind reflect.Kind) BKind
type FlatFunc func(path string, val reflect.Value)
type Type interface{ ... }
    func TypeOf(v any) Type
type Value struct{ ... }
    func ValueOf(v any) Value
    func Wrap(rv reflect.Value) Value
```

## Testings

```shell
go test -v ./reflects/...
```

Test limit by regexp:

```shell
go test -v -run ^TestSetByKeys ./reflects/...
```
