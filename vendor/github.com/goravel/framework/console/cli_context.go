package console

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v3"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/color"
	supportconsole "github.com/goravel/framework/support/console"
)

type CliContext struct {
	arguments []command.Argument
	instance  *cli.Command
}

func NewCliContext(instance *cli.Command, arguments []command.Argument) *CliContext {
	return &CliContext{arguments, instance}
}

func (r *CliContext) Ask(question string, option ...console.AskOption) (string, error) {
	var answer string
	multiple := false

	if len(option) > 0 {
		multiple = option[0].Multiple
		answer = option[0].Default
	}

	if multiple {
		input := huh.NewText().Title(question)
		if len(option) > 0 {
			input.CharLimit(option[0].Limit).Description(option[0].Description).Placeholder(option[0].Placeholder).Lines(option[0].Lines)
			if option[0].Validate != nil {
				input.Validate(option[0].Validate)
			}
		}

		err := input.Value(&answer).Run()
		if err != nil {
			return "", err
		}
	} else {
		input := huh.NewInput().Title(question)

		if len(option) > 0 {
			input.CharLimit(option[0].Limit).Description(option[0].Description).Placeholder(option[0].Placeholder).Prompt(option[0].Prompt)
			if option[0].Validate != nil {
				input.Validate(option[0].Validate)
			}
		}

		err := input.Value(&answer).Run()
		if err != nil {
			return "", err
		}
	}

	return answer, nil
}

func (r *CliContext) Argument(index int) string {
	return r.instance.Args().Get(index)
}

func (r *CliContext) Arguments() []string {
	return r.instance.Args().Slice()
}

func (r *CliContext) CreateProgressBar(total int) console.Progress {
	return NewProgressBar(total)
}

func (r *CliContext) Choice(question string, choices []console.Choice, option ...console.ChoiceOption) (string, error) {
	var answer string

	if len(option) > 0 {
		answer = option[0].Default
	}

	options := make([]huh.Option[string], len(choices))
	for i, choice := range choices {
		options[i] = huh.NewOption[string](choice.Key, choice.Value).Selected(choice.Selected)
	}

	input := huh.NewSelect[string]().Title(question).Options(options...)
	if len(option) > 0 {
		input.Description(option[0].Description)
		if option[0].Validate != nil {
			input.Validate(option[0].Validate)
		}
	}

	err := huh.NewForm(huh.NewGroup(input.Value(&answer))).Run()
	if err != nil {
		return "", err
	}
	return answer, err
}

func (r *CliContext) Comment(message string) {
	color.Debugln(message)
}

func (r *CliContext) Confirm(question string, option ...console.ConfirmOption) bool {
	input := huh.NewConfirm().Title(question)
	answer := false
	if len(option) > 0 {
		if len(option[0].Description) > 0 {
			input.Description(option[0].Description)
		}
		if len(option[0].Affirmative) > 0 {
			input.Affirmative(option[0].Affirmative)
		}
		if len(option[0].Negative) > 0 {
			input.Negative(option[0].Negative)
		}
		answer = option[0].Default
	}

	if err := input.Value(&answer).Run(); err != nil {
		r.Error(err.Error())

		return false
	}

	return answer
}

func (r *CliContext) Error(message string) {
	color.Errorln(message)
}

func (r *CliContext) Info(message string) {
	color.Infoln(message)
}

func (r *CliContext) Instance() *cli.Command {
	return r.instance
}

func (r *CliContext) Line(message string) {
	color.Default().Println(message)
}

func (r *CliContext) MultiSelect(question string, choices []console.Choice, option ...console.MultiSelectOption) ([]string, error) {
	var answer []string

	if len(option) > 0 {
		answer = option[0].Default
	}

	options := make([]huh.Option[string], len(choices))
	for i, choice := range choices {
		options[i] = huh.NewOption(choice.Key, choice.Value).Selected(choice.Selected)
	}

	input := huh.NewMultiSelect[string]().Title(question).Options(options...)
	if len(option) > 0 {
		input.Description(option[0].Description).Limit(option[0].Limit).Filterable(option[0].Filterable)
		if option[0].Validate != nil {
			input.Validate(option[0].Validate)
		}
	}

	err := huh.NewForm(huh.NewGroup(input.Value(&answer))).Run()
	if err != nil {
		return nil, err
	}

	return answer, err
}

func (r *CliContext) NewLine(times ...int) {
	numLines := 1
	if len(times) > 0 && times[0] > 0 {
		numLines = times[0]
	}
	for i := 0; i < numLines; i++ {
		color.Default().Println()
	}
}

func (r *CliContext) Option(key string) string {
	return r.instance.String(key)
}

func (r *CliContext) OptionSlice(key string) []string {
	return r.instance.StringSlice(key)
}

func (r *CliContext) OptionBool(key string) bool {
	return r.instance.Bool(key)
}

func (r *CliContext) OptionFloat64(key string) float64 {
	return r.instance.Float(key)
}

func (r *CliContext) OptionFloat64Slice(key string) []float64 {
	return r.instance.FloatSlice(key)
}

func (r *CliContext) OptionInt(key string) int {
	return r.instance.Int(key)
}

func (r *CliContext) OptionIntSlice(key string) []int {
	return r.instance.IntSlice(key)
}

func (r *CliContext) OptionInt64(key string) int64 {
	return r.instance.Int64(key)
}

func (r *CliContext) OptionInt64Slice(key string) []int64 {
	return r.instance.Int64Slice(key)
}

func (r *CliContext) ArgumentString(key string) string {
	value := r.instance.StringArgs(key)
	if len(value) > 0 {
		return value[0]
	}

	return cast.ToString(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentStringSlice(key string) []string {
	value := r.instance.StringArgs(key)
	if len(value) > 0 {
		return value
	}

	return cast.ToStringSlice(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentFloat32(key string) float32 {
	ret := r.instance.Float32Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToFloat32(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentFloat32Slice(key string) []float32 {
	value := r.instance.Float32Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]float32); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentFloat64(key string) float64 {
	ret := r.instance.Float64Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToFloat64(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentFloat64Slice(key string) []float64 {
	value := r.instance.Float64Args(key)
	if len(value) > 0 {
		return value
	}

	return cast.ToFloat64Slice(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt(key string) int {
	ret := r.instance.IntArgs(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToInt(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentIntSlice(key string) []int {
	value := r.instance.IntArgs(key)
	if len(value) > 0 {
		return value
	}

	return cast.ToIntSlice(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt8(key string) int8 {
	ret := r.instance.Int8Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToInt8(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt8Slice(key string) []int8 {
	value := r.instance.Int8Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]int8); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentInt16(key string) int16 {
	ret := r.instance.Int16Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToInt16(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt16Slice(key string) []int16 {
	value := r.instance.Int16Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]int16); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentInt32(key string) int32 {
	ret := r.instance.Int32Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToInt32(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt32Slice(key string) []int32 {
	value := r.instance.Int32Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]int32); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentInt64(key string) int64 {
	ret := r.instance.Int64Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToInt64(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentInt64Slice(key string) []int64 {
	value := r.instance.Int64Args(key)
	if len(value) > 0 {
		return value
	}

	return cast.ToInt64Slice(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint(key string) uint {
	ret := r.instance.UintArgs(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToUint(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUintSlice(key string) []uint {
	value := r.instance.UintArgs(key)
	if len(value) > 0 {
		return value
	}

	return cast.ToUintSlice(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint8(key string) uint8 {
	ret := r.instance.Uint8Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToUint8(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint8Slice(key string) []uint8 {
	value := r.instance.Uint8Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]uint8); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentUint16(key string) uint16 {
	ret := r.instance.Uint16Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToUint16(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint16Slice(key string) []uint16 {
	value := r.instance.Uint16Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]uint16); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentUint32(key string) uint32 {
	ret := r.instance.Uint32Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToUint32(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint32Slice(key string) []uint32 {
	value := r.instance.Uint32Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]uint32); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentUint64(key string) uint64 {
	ret := r.instance.Uint64Args(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToUint64(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentUint64Slice(key string) []uint64 {
	value := r.instance.Uint64Args(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]uint64); ok {
		return v
	}
	return nil
}

func (r *CliContext) ArgumentTimestamp(key string) time.Time {
	ret := r.instance.TimestampArgs(key)
	if len(ret) > 0 {
		return ret[0]
	}

	return cast.ToTime(r.getDefaultArgumentValue(key))
}

func (r *CliContext) ArgumentTimestampSlice(key string) []time.Time {
	value := r.instance.TimestampArgs(key)
	if len(value) > 0 {
		return value
	}

	defaultValue := r.getDefaultArgumentValue(key)
	if v, ok := defaultValue.([]time.Time); ok {
		return v
	}
	return nil
}

func (r *CliContext) Secret(question string, option ...console.SecretOption) (string, error) {
	var answer string
	if len(option) > 0 {
		answer = option[0].Default
	}

	input := huh.NewInput().Title(question)

	if len(option) > 0 {
		input.CharLimit(option[0].Limit).Description(option[0].Description).Placeholder(option[0].Placeholder).EchoMode(huh.EchoModePassword)
		if option[0].Validate != nil {
			input.Validate(option[0].Validate)
		}
	}

	err := input.Value(&answer).Run()
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (r *CliContext) Spinner(message string, option console.SpinnerOption) error {
	style := lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{TrueColor: "#3D8C8D", ANSI256: "30", ANSI: "6"})
	spin := spinner.New().Title(message).Style(style).TitleStyle(style)

	var err error
	if err := spin.Context(option.Ctx).Action(func() {
		err = option.Action()
	}).Run(); err != nil {
		return err
	}

	return err
}

func (r *CliContext) Success(message string) {
	color.Successln(message)
}

func (r *CliContext) Warning(message string) {
	color.Warningln(message)
}

func (r *CliContext) WithProgressBar(items []any, callback func(any) error) ([]any, error) {
	bar := r.CreateProgressBar(len(items))
	err := bar.Start()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		err := callback(item)
		if err != nil {
			return nil, err
		}
		bar.Advance()
	}

	err = bar.Finish()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *CliContext) TwoColumnDetail(first, second string, filler ...rune) {
	r.Line(supportconsole.TwoColumnDetail(first, second, filler...))
}

func (r *CliContext) Divider(filler ...string) {
	var str string
	if len(filler) == 0 || len(filler[0]) == 0 {
		str = "-"
	} else {
		str = filler[0]
	}

	width := pterm.GetTerminalWidth()
	strLen := utf8.RuneCountInString(str)

	repeat, remainder := width/strLen, width%strLen

	message := strings.Repeat(str, repeat)

	if remainder > 0 {
		message += string([]rune(str)[:remainder])
	}

	r.Line(message)
}

func (r *CliContext) Green(message string) {
	color.Green().Print(message)
}

func (r *CliContext) Greenln(message string) {
	color.Green().Println(message)
}

func (r *CliContext) Red(message string) {
	color.Red().Print(message)
}

func (r *CliContext) Redln(message string) {
	color.Red().Println(message)
}

func (r *CliContext) Yellow(message string) {
	color.Yellow().Print(message)
}

func (r *CliContext) Yellowln(message string) {
	color.Yellow().Println(message)
}

func (r *CliContext) Black(message string) {
	color.Black().Print(message)
}

func (r *CliContext) Blackln(message string) {
	color.Black().Println(message)
}

func (r *CliContext) getDefaultArgumentValue(key string) any {
	for _, arg := range r.arguments {
		if arg.GetName() == key && arg.GetValue() != nil {
			return arg.GetValue()
		}
	}

	return nil
}
