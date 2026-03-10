package console

import (
	"context"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/goravel/framework/contracts/console/command"
)

type Command interface {
	// Signature set the unique signature for the command.
	Signature() string
	// Description the console command description.
	Description() string
	// Extend the console command extend.
	Extend() command.Extend
	// Handle execute the console command.
	Handle(ctx Context) error
}

type Context interface {
	// Ask prompts the user for input.
	Ask(question string, option ...AskOption) (string, error)
	// CreateProgressBar creates a new progress bar instance.
	CreateProgressBar(total int) Progress
	// Choice prompts the user to select from a list of options.
	Choice(question string, options []Choice, option ...ChoiceOption) (string, error)
	// Comment writes a comment message to the console.
	Comment(message string)
	// Confirm prompts the user for a confirmation.
	Confirm(question string, option ...ConfirmOption) bool
	// Argument get the value of a command argument.
	Argument(index int) string
	// Arguments get all the arguments passed to command.
	Arguments() []string
	// Info writes an information message to the console.
	Info(message string)
	// Instance gets the underlying cli.Command instance.
	Instance() *cli.Command
	// Error writes an error message to the console.
	Error(message string)
	// Line writes a string to the console.
	Line(message string)
	// MultiSelect prompts the user to select multiple options from a list of options.
	MultiSelect(question string, options []Choice, option ...MultiSelectOption) ([]string, error)
	// NewLine writes a newline character to the console.
	NewLine(times ...int)
	// Option gets the value of a command option.
	Option(key string) string
	// OptionSlice looks up the value of a local StringSliceFlag, returns nil if not found
	OptionSlice(key string) []string
	// OptionBool looks up the value of a local BoolFlag, returns false if not found
	OptionBool(key string) bool
	// OptionFloat64 looks up the value of a local Float64Flag, returns zero if not found
	OptionFloat64(key string) float64
	// OptionFloat64Slice looks up the value of a local Float64SliceFlag, returns nil if not found
	OptionFloat64Slice(key string) []float64
	// OptionInt looks up the value of a local IntFlag, returns zero if not found
	OptionInt(key string) int
	// OptionIntSlice looks up the value of a local IntSliceFlag, returns nil if not found
	OptionIntSlice(key string) []int
	// OptionInt64 looks up the value of a local Int64Flag, returns zero if not found
	OptionInt64(key string) int64
	// OptionInt64Slice looks up the value of a local Int64SliceFlag, returns nil if not found
	OptionInt64Slice(key string) []int64
	// ArgumentString looks up the value of a local ArgumentString, returns empty string if not found
	ArgumentString(key string) string
	// ArgumentStringSlice looks up the value of a local ArgumentStringSlice, returns nil if not found
	ArgumentStringSlice(key string) []string
	// ArgumentFloat32 looks up the value of a local ArgumentFloat32, returns zero if not found
	ArgumentFloat32(key string) float32
	// ArgumentFloat32Slice looks up the value of a local ArgumentFloat32Slice, returns nil if not found
	ArgumentFloat32Slice(key string) []float32
	// ArgumentFloat64 looks up the value of a local ArgumentFloat64, returns zero if not found
	ArgumentFloat64(key string) float64
	// ArgumentFloat64Slice looks up the value of a local ArgumentFloat64Slice, returns nil if not found
	ArgumentFloat64Slice(key string) []float64
	// ArgumentInt looks up the value of a local ArgumentInt, returns zero if not found
	ArgumentInt(key string) int
	// ArgumentIntSlice looks up the value of a local ArgumentIntSlice, returns nil if not found
	ArgumentIntSlice(key string) []int
	// ArgumentInt8 looks up the value of a local ArgumentInt8, returns zero if not found
	ArgumentInt8(key string) int8
	// ArgumentInt8Slice looks up the value of a local ArgumentInt8Slice, returns nil if not found
	ArgumentInt8Slice(key string) []int8
	// ArgumentInt16 looks up the value of a local ArgumentInt16, returns zero if not found
	ArgumentInt16(key string) int16
	// ArgumentInt16Slice looks up the value of a local ArgumentInt16Slice, returns nil if not found
	ArgumentInt16Slice(key string) []int16
	// ArgumentInt32 looks up the value of a local ArgumentInt32, returns zero if not found
	ArgumentInt32(key string) int32
	// ArgumentInt32Slice looks up the value of a local ArgumentInt32Slice, returns nil if not found
	ArgumentInt32Slice(key string) []int32
	// ArgumentInt64 looks up the value of a local ArgumentInt64, returns zero if not found
	ArgumentInt64(key string) int64
	// ArgumentInt64Slice looks up the value of a local ArgumentInt64Slice, returns nil if not found
	ArgumentInt64Slice(key string) []int64
	// ArgumentUint looks up the value of a local ArgumentUint, returns zero if not found
	ArgumentUint(key string) uint
	// ArgumentUint8 looks up the value of a local ArgumentUint8, returns zero if not found
	ArgumentUint8(key string) uint8
	// ArgumentUint16 looks up the value of a local ArgumentUint16, returns zero if not found
	ArgumentUint16(key string) uint16
	// ArgumentUint32 looks up the value of a local ArgumentUint32, returns zero if not found
	ArgumentUint32(key string) uint32
	// ArgumentUint64 looks up the value of a local ArgumentUint64, returns zero if not found
	ArgumentUint64(key string) uint64
	// ArgumentUintSlice looks up the value of a local ArgumentUintSlice, returns nil if not found
	ArgumentUintSlice(key string) []uint
	// ArgumentUint8Slice looks up the value of a local ArgumentUint8Slice, returns nil if not found
	ArgumentUint8Slice(key string) []uint8
	// ArgumentUint16Slice looks up the value of a local ArgumentUint16Slice, returns nil if not found
	ArgumentUint16Slice(key string) []uint16
	// ArgumentUint32Slice looks up the value of a local ArgumentUint32Slice, returns nil if not found
	ArgumentUint32Slice(key string) []uint32
	// ArgumentUint64Slice looks up the value of a local ArgumentUint64Slice, returns nil if not found
	ArgumentUint64Slice(key string) []uint64
	// ArgumentTimestamp looks up the value of a local ArgumentTimestamp, returns zero if not found
	ArgumentTimestamp(key string) time.Time
	// ArgumentTimestampSlice looks up the value of a local ArgumentTimestampSlice, returns nil if not found
	ArgumentTimestampSlice(key string) []time.Time
	// Secret prompts the user for a password.
	Secret(question string, option ...SecretOption) (string, error)
	// Spinner creates a new spinner instance.
	Spinner(message string, option SpinnerOption) error
	// Success writes a success message to the console.
	Success(message string)
	// Warning writes a warning message to the console.
	Warning(message string)
	// WithProgressBar executes a callback with a progress bar.
	WithProgressBar(items []any, callback func(any) error) ([]any, error)
	// TwoColumnDetail writes a two column detail to the console.
	TwoColumnDetail(first, second string, filler ...rune)
	// Divider shows a terminal-width divider filled with given sting
	Divider(filler ...string)
	// Green writes green text to console
	Green(message string)
	// Greenln writes green line to console
	Greenln(message string)
	// Red writes red text to console
	Red(message string)
	// Redln writes red line to console
	Redln(message string)
	// Yellow writes yellow text to console
	Yellow(message string)
	// Yellowln writes yellow line to console
	Yellowln(message string)
	// Black writes black text to console
	Black(message string)
	// Blackln writes black line to console
	Blackln(message string)
}

type Progress interface {
	// Advance advances the progress bar by a given step.
	Advance(step ...int)
	// Finish completes the progress bar.
	Finish() error
	// ShowElapsedTime sets if the elapsed time should be displayed in the progress bar.
	ShowElapsedTime(b ...bool) Progress
	// ShowTitle sets the title of the progress bar.
	ShowTitle(b ...bool) Progress
	// SetTitle sets the message of the progress bar.
	SetTitle(message string)
	// Start starts the progress bar.
	Start() error
}

type Choice struct {
	// Key the choice key.
	Key string
	// Value the choice value.
	Value string
	// Selected determines if the choice is selected.
	Selected bool
}

type AskOption struct {
	// Validate the input validation function.
	Validate func(string) error
	// Default the default value for the input.
	Default string
	// Description the input description.
	Description string
	// Placeholder the input placeholder.
	Placeholder string
	// Prompt the prompt message.(use for single line input)
	Prompt string
	// Lines the number of lines for the input.(use for multiple lines text)
	Lines int
	// Limit the character limit for the input.
	Limit int
	// Multiple determines if input is single line or multiple lines text
	Multiple bool
}

type ChoiceOption struct {
	// Validate the input validation function.
	Validate func(string) error
	// Default the default value for the input.
	Default string
	// Description the input description.
	Description string
}

type ConfirmOption struct {
	// Affirmative label for the affirmative button.
	Affirmative string
	// Description the input description.
	Description string
	// Negative label for the negative button.
	Negative string
	// Default the default value for the input.
	Default bool
}

type SecretOption struct {
	// Validate the input validation function.
	Validate func(string) error
	// Default the default value for the input.
	Default string
	// Description the input description.
	Description string
	// Placeholder the input placeholder.
	Placeholder string
	// Limit the character limit for the input.
	Limit int
}

type MultiSelectOption struct {
	// Validate the input validation function.
	Validate func([]string) error
	// Description the input description.
	Description string
	// Default the default value for the input.
	Default []string
	// Limit the number of choices that can be selected.
	Limit int
	// Filterable determines if the choices can be filtered.
	Filterable bool
}

type SpinnerOption struct {
	// Ctx the context for the spinner.
	Ctx context.Context
	// Action the action to execute.
	Action func() error
}
