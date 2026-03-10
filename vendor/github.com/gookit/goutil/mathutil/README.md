# Math Utils

- some features

## Install

```bash
go get github.com/gookit/goutil/mathutil
```

## Go docs

- [Go Docs](https://pkg.go.dev/github.com/gookit/goutil)

## Usage


## Functions

```go

func CompFloat[T comdef.Float](first, second T, op string) (ok bool)
func CompInt[T comdef.Xint](first, second T, op string) (ok bool)
func CompInt64(first, second int64, op string) bool
func CompValue[T comdef.XintOrFloat](first, second T, op string) (ok bool)
func Compare(first, second any, op string) (ok bool)
func DataSize(size uint64) string
func ElapsedTime(startTime time.Time) string
func Float(in any) (float64, error)
func FloatOr(in any, defVal float64) float64
func FloatOrDefault(in any, defVal float64) float64
func FloatOrErr(in any) (float64, error)
func FloatOrPanic(in any) float64
func GreaterOr[T comdef.XintOrFloat](val, min, defVal T) T
func GteOr[T comdef.XintOrFloat](val, min, defVal T) T
func HowLongAgo(sec int64) string
func InRange[T comdef.IntOrFloat](val, min, max T) bool
func InUintRange[T comdef.Uint](val, min, max T) bool
func Int(in any) (int, error)
func Int64(in any) (int64, error)
func Int64OrErr(in any) (int64, error)
func IntOr(in any, defVal int) int
func IntOrDefault(in any, defVal int) int
func IntOrErr(in any) (iVal int, err error)
func IntOrPanic(in any) int
func IsNumeric(c byte) bool
func LessOr[T comdef.XintOrFloat](val, max, devVal T) T
func LteOr[T comdef.XintOrFloat](val, max, devVal T) T
func Max[T comdef.XintOrFloat](x, y T) T
func MaxFloat(x, y float64) float64
func MaxI64(x, y int64) int64
func MaxInt(x, y int) int
func Min[T comdef.XintOrFloat](x, y T) T
func MustFloat(in any) float64
func MustInt(in any) int
func MustInt64(in any) int64
func MustString(val any) string
func MustUint(in any) uint64
func OrElse[T comdef.XintOrFloat](val, defVal T) T
func OutRange[T comdef.IntOrFloat](val, min, max T) bool
func Percent(val, total int) float64
func QuietFloat(in any) float64
func QuietInt(in any) int
func QuietInt64(in any) int64
func QuietString(val any) string
func QuietUint(in any) uint64
func RandInt(min, max int) int
func RandIntWithSeed(min, max int, seed int64) int
func RandomInt(min, max int) int
func RandomIntWithSeed(min, max int, seed int64) int
func SafeFloat(in any) float64
func SafeInt(in any) int
func SafeInt64(in any) int64
func SafeUint(in any) uint64
func StrInt(s string) int
func StrIntOr(s string, defVal int) int
func String(val any) string
func StringOrErr(val any) (string, error)
func StringOrPanic(val any) string
func SwapMax[T comdef.XintOrFloat](x, y T) (T, T)
func SwapMaxI64(x, y int64) (int64, int64)
func SwapMaxInt(x, y int) (int, int)
func SwapMin[T comdef.XintOrFloat](x, y T) (T, T)
func ToFloat(in any) (f64 float64, err error)
func ToFloatWithFunc(in any, usrFn func(any) (float64, error)) (f64 float64, err error)
func ToInt(in any) (iVal int, err error)
func ToInt64(in any) (i64 int64, err error)
func ToString(val any) (string, error)
func ToUint(in any) (u64 uint64, err error)
func ToUintWithFunc(in any, usrFn func(any) (uint64, error)) (u64 uint64, err error)
func TryToString(val any, defaultAsErr bool) (str string, err error)
func Uint(in any) (uint64, error)
func UintOrErr(in any) (uint64, error)
func ZeroOr[T comdef.XintOrFloat](val, defVal T) T

```

## Testings

```shell
go test -v ./mathutil/...
```

Test limit by regexp:

```shell
go test -v -run ^TestSetByKeys ./mathutil/...
```
