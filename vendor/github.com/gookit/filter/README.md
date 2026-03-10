# Filter

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/filter)](https://github.com/gookit/filter)
[![Actions Status](https://github.com/gookit/filter/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/filter/actions)
[![Coverage Status](https://coveralls.io/repos/github/gookit/filter/badge.svg?branch=master)](https://coveralls.io/github/gookit/filter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/filter)](https://goreportcard.com/report/github.com/gookit/filter)
[![Go Reference](https://pkg.go.dev/badge/github.com/gookit/filter.svg)](https://pkg.go.dev/github.com/gookit/filter)

`filter` - provide filtering, sanitizing, and conversion of Golang data.

> 中文说明请查看 **[README.zh-CN](README.zh-CN.md)**

## GoDoc

- [godoc](https://pkg.go.dev/github.com/gookit/filter)

> **NOTE**: To filter and validate Map, Struct data. Please use [gookit/validate](https://github.com/gookit/validate)

## Install

```shell
go get github.com/gookit/filter
```

## Func Usage

Quick usage:

```go
str := filter.MustString(23) // "23"

intVal, err := filter.Int("20") // int(20)
strings := filter.Str2Slice("a,b, c", ",") // []string{"a", "b", "c"}
```

## Filtration

Filtering data:

```go
data := map[string]any{
    "name":     " inhere ",
    "age":      "50",
    "money":    "50.34",
    // 
    "remember": "yes",
    //
    "sub1": []string{"1", "2"},
    "tags": "go;lib",
    "str1": " word ",
    "ids":  []int{1, 2, 2, 1},
}
f := filter.New(data)
f.AddRule("money", "float")
f.AddRule("remember", "bool")
f.AddRule("sub1", "strings2ints")
f.AddRule("tags", "str2arr:;")
f.AddRule("ids", "unique")
f.AddRule("str1", "ltrim|rtrim")
f.AddRule("not-exist", "unique")
// add multi
f.AddRules(map[string]string{
    "age": "trim|int",
    "name": "trim|ucFirst",
})

// apply all added rules for data.
f.Filtering() 

// get filtered data
newData := f.CleanData()
fmt.Printf("%#v\n", newData)
// f.BindStruct(&user)
```

**Output**:

```go
map[string]interface {}{
    "remember":true, 
    "sub1":[]int{1, 2}, 
    "tags":[]string{"go", "lib"}, 
    "ids":[]int{2, 1}, 
    "str1":"word", 
    "name":"INHERE", 
    "age":50, 
    "money":50.34
}
```

## Filters & Converters

- `ToBool/Bool(s string) (bool, error)`
- `ToFloat/Float(v interface{}) (float64, error)`
- `ToInt/Int(v interface{}) (int, error)`
- `ToUint/Uint(v interface{}) (uint64, error)`
- `ToInt64/Int64(v interface{}) (int64, error)`
- `ToString/String(v interface{}) (string, error)`
- `MustBool(s string) bool`
- `MustFloat(s string) float64`
- `MustInt(s string) int`
- `MustInt64(s string) int64`
- `MustUint(s string) uint64`
- `MustString(v interface{}) string`
- `Trim(s string, cutSet ...string) string`
- `TrimLeft(s string, cutSet ...string) string`
- `TrimRight(s string, cutSet ...string) string`
- `TrimStrings(ss []string, cutSet ...string) (ns []string)`
- `Substr(s string, pos, length int) string`
- `Lower/Lowercase(s string) string`
- `Upper/Uppercase(s string) string`
- `LowerFirst(s string) string`
- `UpperFirst(s string) string`
- `UpperWord(s string) string`
- `Camel/CamelCase(s string, sep ...string) string`
- `Snake/SnakeCase(s string, sep ...string) string`
- `Email(s string) string`
- `URLDecode(s string) string`
- `URLEncode(s string) string`
- `EscapeJS(s string) string`
- `EscapeHTML(s string) string`
- `Unique(val interface{}) interface{}` Will remove duplicate values, use for `[]int` `[]int64` `[]string`
- `StrToSlice(s string, sep ...string) []string`
- `StrToInts(s string, sep ...string) (ints []int, err error)`
- `StrToTime(s string, layouts ...string) (t time.Time, err error)`
- `StringsToInts(ss []string) (ints []int, err error)`

## License

**[MIT](LICENSE)**
