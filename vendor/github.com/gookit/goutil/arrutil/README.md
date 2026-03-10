# Array/Slice Utils

## Install

```shell
go get github.com/gookit/goutil/arrutil
```

## Go docs

- [Go docs](https://pkg.go.dev/github.com/gookit/goutil/arrutil)

## Functions API

> **Note**: doc by run `go doc ./arrutil`

```go
func AnyToString(arr any) string
func CloneSlice(data any) interface{}
func Contains(arr, val any) bool
func ExceptWhile(data any, fn Predicate) interface{}
func Excepts(first, second any, fn Comparer) interface{}
func Find(source any, fn Predicate) (interface{}, error)
func FindOrDefault(source any, fn Predicate, defaultValue any) interface{}
func FormatIndent(arr any, indent string) string
func GetRandomOne(arr any) interface{}
func HasValue(arr, val any) bool
func InStrings(elem string, ss []string) bool
func Int64sHas(ints []int64, val int64) bool
func Intersects(first any, second any, fn Comparer) interface{}
func IntsHas(ints []int, val int) bool
func JoinSlice(sep string, arr ...any) string
func JoinStrings(sep string, ss ...string) string
func MakeEmptySlice(itemType reflect.Type) interface{}
func Map[T any, V any](list []T, mapFn func(obj T) (val V, find bool)) []V
func Column[T any, V any](list []T, mapFn func(obj T) (val V, find bool)) []V
func MustToInt64s(arr any) []int64
func MustToStrings(arr any) []string
func NotContains(arr, val any) bool
func RandomOne(arr any) interface{}
func Reverse(ss []string)
func SliceToInt64s(arr []any) []int64
func SliceToString(arr ...any) string
func SliceToStrings(arr []any) []string
func StringsFilter(ss []string, filter ...func(s string) bool) []string
func StringsHas(ss []string, val string) bool
func StringsJoin(sep string, ss ...string) string
func StringsMap(ss []string, mapFn func(s string) string) []string
func StringsRemove(ss []string, s string) []string
func StringsToInts(ss []string) (ints []int, err error)
func StringsToSlice(ss []string) []interface{}
func TakeWhile(data any, fn Predicate) interface{}
func ToInt64s(arr any) (ret []int64, err error)
func ToString(arr []any) string
func ToStrings(arr any) (ret []string, err error)
func TrimStrings(ss []string, cutSet ...string) []string
func TwowaySearch(data any, item any, fn Comparer) (int, error)
func Union(first, second any, fn Comparer) interface{}
func Unique(arr any) interface{}
type ArrFormatter struct{ ... }
    func NewFormatter(arr any) *ArrFormatter
```

## Code Check & Testing

```bash
gofmt -w -l ./
golint ./...
```

**Testing**:

```shell
go test -v ./cliutil/...
```

**Test limit by regexp**:

```shell
go test -v -run ^TestSetByKeys ./cliutil/...
```

## Refers

- https://github.com/elliotchance/pie
