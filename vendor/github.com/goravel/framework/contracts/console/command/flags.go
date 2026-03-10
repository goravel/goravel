package command

const (
	FlagTypeBool         = "bool"
	FlagTypeFloat64      = "float64"
	FlagTypeFloat64Slice = "float64_slice"
	FlagTypeInt          = "int"
	FlagTypeIntSlice     = "int_slice"
	FlagTypeInt64        = "int64"
	FlagTypeInt64Slice   = "int64_slice"
	FlagTypeString       = "string"
	FlagTypeStringSlice  = "string_slice"
)

type Extend struct {
	ArgsUsage string
	Category  string
	Flags     []Flag
	Arguments []Argument
}

type Flag interface {
	// Type gets a flag type.
	Type() string
}

type BoolFlag struct {
	Name               string
	Usage              string
	Aliases            []string
	DisableDefaultText bool
	Required           bool
	Value              bool
}

func (receiver *BoolFlag) Type() string {
	return FlagTypeBool
}

type Float64Flag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    float64
	Required bool
}

func (receiver *Float64Flag) Type() string {
	return FlagTypeFloat64
}

type Float64SliceFlag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    []float64
	Required bool
}

func (receiver *Float64SliceFlag) Type() string {
	return FlagTypeFloat64Slice
}

type IntFlag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    int
	Required bool
}

func (receiver *IntFlag) Type() string {
	return FlagTypeInt
}

type IntSliceFlag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    []int
	Required bool
}

func (receiver *IntSliceFlag) Type() string {
	return FlagTypeIntSlice
}

type Int64Flag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    int64
	Required bool
}

func (receiver *Int64Flag) Type() string {
	return FlagTypeInt64
}

type Int64SliceFlag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    []int64
	Required bool
}

func (receiver *Int64SliceFlag) Type() string {
	return FlagTypeInt64Slice
}

type StringFlag struct {
	Name     string
	Usage    string
	Value    string
	Aliases  []string
	Required bool
}

func (receiver *StringFlag) Type() string {
	return FlagTypeString
}

type StringSliceFlag struct {
	Name     string
	Usage    string
	Aliases  []string
	Value    []string
	Required bool
}

func (receiver *StringSliceFlag) Type() string {
	return FlagTypeStringSlice
}
