package command

import "time"

type Argument interface {
	GetMin() int
	GetMax() int
	GetName() string
	GetValue() any
}

type ArgumentBase[T any] struct {
	Name     string // the name of this argument
	Default  T      // the default value of this argument
	Usage    string // the usage text to show
	Required bool   // if this argument is required
}

func (a *ArgumentBase[T]) GetMin() int {
	if a.Required {
		return 1
	} else {
		return 0
	}
}

func (a *ArgumentBase[T]) GetMax() int {
	return 1
}

func (a *ArgumentBase[T]) GetName() string {
	return a.Name
}

func (a *ArgumentBase[T]) GetValue() any {
	return a.Default
}

type ArgumentsBase[T any] struct {
	Name    string // the name of this argument
	Default []T    // the default value of this argument
	Usage   string // the usage text to show
	Min     int    // the min num of occurrences of this argument
	Max     int    // the max num of occurrences of this argument, set to -1 for unlimited
}

func (a *ArgumentsBase[T]) GetMin() int {
	return a.Min
}

func (a *ArgumentsBase[T]) GetMax() int {
	return a.Max
}

func (a *ArgumentsBase[T]) GetName() string {
	return a.Name
}

func (a *ArgumentsBase[T]) GetValue() any {
	return a.Default
}

type (
	ArgumentFloat32 = ArgumentBase[float32]
	ArgumentFloat64 = ArgumentBase[float64]
	ArgumentInt     = ArgumentBase[int]
	ArgumentInt8    = ArgumentBase[int8]
	ArgumentInt16   = ArgumentBase[int16]
	ArgumentInt32   = ArgumentBase[int32]
	ArgumentInt64   = ArgumentBase[int64]
	ArgumentString  = ArgumentBase[string]
	ArgumentUint    = ArgumentBase[uint]
	ArgumentUint8   = ArgumentBase[uint8]
	ArgumentUint16  = ArgumentBase[uint16]
	ArgumentUint32  = ArgumentBase[uint32]
	ArgumentUint64  = ArgumentBase[uint64]

	ArgumentFloat32Slice = ArgumentsBase[float32]
	ArgumentFloat64Slice = ArgumentsBase[float64]
	ArgumentIntSlice     = ArgumentsBase[int]
	ArgumentInt8Slice    = ArgumentsBase[int8]
	ArgumentInt16Slice   = ArgumentsBase[int16]
	ArgumentInt32Slice   = ArgumentsBase[int32]
	ArgumentInt64Slice   = ArgumentsBase[int64]
	ArgumentStringSlice  = ArgumentsBase[string]
	ArgumentUintSlice    = ArgumentsBase[uint]
	ArgumentUint8Slice   = ArgumentsBase[uint8]
	ArgumentUint16Slice  = ArgumentsBase[uint16]
	ArgumentUint32Slice  = ArgumentsBase[uint32]
	ArgumentUint64Slice  = ArgumentsBase[uint64]
)

type ArgumentTimestamp struct {
	Name     string    // the name of this argument
	Value    time.Time // the default value of this argument
	Usage    string    // the usage text to show
	Required bool      // if this argument is required

	// Available layouts for timestamp value.
	// Note that value for formats with missing year/date will be interpreted as current year/date respectively.
	// Read more about time layouts: https://pkg.go.dev/time#pkg-constants
	Layouts []string
}

func (a *ArgumentTimestamp) GetMin() int {
	if a.Required {
		return 1
	} else {
		return 0
	}
}

func (a *ArgumentTimestamp) GetMax() int {
	return 1
}

func (a *ArgumentTimestamp) GetName() string {
	return a.Name
}

func (a *ArgumentTimestamp) GetValue() any {
	return a.Value
}

type ArgumentTimestampSlice struct {
	Name  string    // the name of this argument
	Value time.Time // the default value of this argument
	Usage string    // the usage text to show
	Min   int       // the min num of occurrences of this argument
	Max   int       // the max num of occurrences of this argument, set to -1 for unlimited

	// Available layouts for timestamp value.
	// Note that value for formats with missing year/date will be interpreted as current year/date respectively.
	// Read more about time layouts: https://pkg.go.dev/time#pkg-constants
	Layouts []string
}

func (a *ArgumentTimestampSlice) GetMin() int {
	return a.Min
}

func (a *ArgumentTimestampSlice) GetMax() int {
	return a.Max
}

func (a *ArgumentTimestampSlice) GetName() string {
	return a.Name
}

func (a *ArgumentTimestampSlice) GetValue() any {
	return a.Value
}
