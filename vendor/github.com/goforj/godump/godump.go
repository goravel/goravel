package godump

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/tabwriter"
	"unicode/utf8"
	"unsafe"
)

const (
	colorReset   = "\033[0m"
	colorGray    = "\033[90m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorRedBg   = "\033[48;2;34;16;16m"
	colorGreenBg = "\033[48;2;16;34;22m"
	colorLime    = "\033[1;38;5;113m"
	colorCyan    = "\033[38;5;38m"
	colorNote    = "\033[38;5;38m"
	colorRef     = "\033[38;5;247m"
	colorMeta    = "\033[38;5;170m"
	colorDefault = "\033[38;5;208m"
	indentWidth  = 2
)

// Default configuration values for the Dumper.
const (
	defaultDisableStringer = false
	defaultMaxDepth        = 15
	defaultMaxItems        = 100
	defaultMaxStringLen    = 100000
	defaultMaxStackDepth   = 10
	initialCallerSkip      = 2
)

const (
	// FieldMatchExact matches field names exactly (case-insensitive).
	FieldMatchExact FieldMatchMode = iota
	// FieldMatchContains matches if the field name contains a substring (case-insensitive).
	FieldMatchContains
	// FieldMatchPrefix matches if the field name starts with a substring (case-insensitive).
	FieldMatchPrefix
	// FieldMatchSuffix matches if the field name ends with a substring (case-insensitive).
	FieldMatchSuffix
)

// FieldMatchMode controls how field names are matched.
type FieldMatchMode int

var defaultRedactedFields = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"access_key",
	"accesskey",
	"private_key",
	"privatekey",
	"client_secret",
	"clientsecret",
	"refresh_token",
	"session",
	"cookie",
	"jwt",
	"bearer",
	"authorization",
	"signature",
	"signing_key",
}

// defaultDumper is the default Dumper instance used by Dump and DumpStr functions.
var defaultDumper = NewDumper()

// exitFunc is a function that can be overridden for testing purposes.
var exitFunc = os.Exit

// Colorizer is a function type that takes a color code and a string, returning the colorized string.
type Colorizer func(code, str string) string

// colorizeUnstyled returns the string without any colorization.
//
// It satisfies the [Colorizer] interface.
func colorizeUnstyled(code, str string) string {
	return str // No colorization
}

// colorizeANSI colorizes the string using ANSI escape codes.
//
// It satisfies the [Colorizer] interface.
func colorizeANSI(code, str string) string {
	return code + str + colorReset
}

// htmlColorMap maps color codes to HTML colors.
var htmlColorMap = map[string]string{
	colorGray:    "#999",
	colorYellow:  "#ffb400",
	colorRed:     "#ff5f5f",
	colorGreen:   "#55d655",
	colorLime:    "#80ff80",
	colorNote:    "#40c0ff",
	colorRef:     "#aaa",
	colorMeta:    "#d087d0",
	colorDefault: "#ff7f00",
}

// colorizeHTML colorizes the string using HTML span tags.
//
// It satisfies the [Colorizer] interface.
func colorizeHTML(code, str string) string {
	return fmt.Sprintf(`<span style="color:%s">%s</span>`, htmlColorMap[code], str)
}

// Dumper holds configuration for dumping structured data.
// It controls depth, item count, and string length limits.
type Dumper struct {
	maxDepth           int
	maxItems           int
	maxStringLen       int
	writer             io.Writer
	skippedStackFrames int
	disableStringer    bool
	disableColor       bool
	disableHeader      bool
	includeFields      []string
	excludeFields      []string
	redactFields       []string
	fieldMatchMode     FieldMatchMode
	redactMatchMode    FieldMatchMode

	// callerFn is used to get the caller information.
	// It defaults to [runtime.Caller], it is here to be overridden for testing purposes.
	callerFn func(skip int) (uintptr, string, int, bool)

	// colorizer is used to apply color formatting to the output.
	colorizer Colorizer
}

// Option defines a functional option for configuring a Dumper.
type Option func(*Dumper) *Dumper

// dumpState tracks reference ids for a single dump call.
type dumpState struct {
	nextRefID int
	refs      map[uintptr]int
}

// newDumpState initializes per-dump reference tracking.
func newDumpState() *dumpState {
	return &dumpState{
		nextRefID: 1,
		refs:      map[uintptr]int{},
	}
}

// WithMaxDepth limits how deep the structure will be dumped.
// Param n must be 0 or greater or this will be ignored, and default MaxDepth will be 15.
// @group Options
//
// Example: limit depth
//
//	// Default: 15
//	v := map[string]map[string]int{"a": {"b": 1}}
//	d := godump.NewDumper(godump.WithMaxDepth(1))
//	d.Dump(v)
//	// #map[string]map[string]int {
//	//   a => #map[string]int {
//	//     b => 1 #int
//	//   }
//	// }
func WithMaxDepth(n int) Option {
	return func(d *Dumper) *Dumper {
		if n >= 0 {
			d.maxDepth = n
		}
		return d
	}
}

// WithMaxItems limits how many items from an array, slice, or map can be printed.
// Param n must be 0 or greater or this will be ignored, and default MaxItems will be 100.
// @group Options
//
// Example: limit items
//
//	// Default: 100
//	v := []int{1, 2, 3}
//	d := godump.NewDumper(godump.WithMaxItems(2))
//	d.Dump(v)
//	// #[]int [
//	//   0 => 1 #int
//	//   1 => 2 #int
//	//   ... (truncated)
//	// ]
func WithMaxItems(n int) Option {
	return func(d *Dumper) *Dumper {
		if n >= 0 {
			d.maxItems = n
		}
		return d
	}
}

// WithMaxStringLen limits how long printed strings can be.
// Param n must be 0 or greater or this will be ignored, and default MaxStringLen will be 100000.
// @group Options
//
// Example: limit string length
//
//	// Default: 100000
//	v := "hello world"
//	d := godump.NewDumper(godump.WithMaxStringLen(5))
//	d.Dump(v)
//	// "hello…" #string
func WithMaxStringLen(n int) Option {
	return func(d *Dumper) *Dumper {
		if n >= 0 {
			d.maxStringLen = n
		}
		return d
	}
}

// WithWriter routes output to the provided writer.
// @group Options
//
// Example: write to buffer
//
//	// Default: stdout
//	var b strings.Builder
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper(godump.WithWriter(&b))
//	d.Dump(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func WithWriter(w io.Writer) Option {
	return func(d *Dumper) *Dumper {
		d.writer = w
		return d
	}
}

// WithSkipStackFrames skips additional stack frames for header reporting.
// This is useful when godump is wrapped and the actual call site is deeper.
// @group Options
//
// Example: skip wrapper frames
//
//	// Default: 0
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper(godump.WithSkipStackFrames(2))
//	d.Dump(v)
//	// <#dump // ../../../../usr/local/go/src/runtime/asm_arm64.s:1223
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func WithSkipStackFrames(n int) Option {
	return func(d *Dumper) *Dumper {
		if n >= 0 {
			d.skippedStackFrames = n
		}
		return d
	}
}

// WithDisableStringer disables using the fmt.Stringer output.
// When enabled, the underlying type is rendered instead of String().
// @group Options
//
// Example: show raw types
//
//	// Default: false
//	v := time.Duration(3)
//	d := godump.NewDumper(godump.WithDisableStringer(true))
//	d.Dump(v)
//	// 3 #time.Duration
func WithDisableStringer(b bool) Option {
	return func(d *Dumper) *Dumper {
		d.disableStringer = b
		return d
	}
}

// WithoutColor disables colorized output for the dumper.
// @group Options
//
// Example: disable colors
//
//	// Default: false
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper(godump.WithoutColor())
//	d.Dump(v)
//	// (prints without color)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func WithoutColor() Option {
	return func(d *Dumper) *Dumper {
		d.disableColor = true
		d.colorizer = colorizeUnstyled
		return d
	}
}

// WithoutHeader disables printing the source location header.
// @group Options
//
// Example: disable header
//
//	// Default: false
//	d := godump.NewDumper(godump.WithoutHeader())
//	d.Dump("hello")
//	// "hello" #string
func WithoutHeader() Option {
	return func(d *Dumper) *Dumper {
		d.disableHeader = true
		return d
	}
}

// WithOnlyFields limits struct output to fields that match the provided names.
// @group Options
//
// Example: include-only fields
//
//	// Default: none
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//	}
//	d := godump.NewDumper(
//		godump.WithOnlyFields("ID", "Email"),
//	)
//	d.Dump(User{ID: 1, Email: "user@example.com", Password: "secret"})
//	// #godump.User {
//	//   +ID    => 1 #int
//	//   +Email => "user@example.com" #string
//	// }
func WithOnlyFields(names ...string) Option {
	return func(d *Dumper) *Dumper {
		d.includeFields = append(d.includeFields, names...)
		return d
	}
}

// WithExcludeFields omits struct fields that match the provided names.
// @group Options
//
// Example: exclude fields
//
//	// Default: none
//	type User struct {
//		ID       int
//		Email    string
//		Password string
//	}
//	d := godump.NewDumper(
//		godump.WithExcludeFields("Password"),
//	)
//	d.Dump(User{ID: 1, Email: "user@example.com", Password: "secret"})
//	// #godump.User {
//	//   +ID    => 1 #int
//	//   +Email => "user@example.com" #string
//	// }
func WithExcludeFields(names ...string) Option {
	return func(d *Dumper) *Dumper {
		d.excludeFields = append(d.excludeFields, names...)
		return d
	}
}

// WithFieldMatchMode sets how field names are matched for WithExcludeFields.
// @group Options
//
// Example: use substring matching
//
//	// Default: FieldMatchExact
//	type User struct {
//		UserID int
//	}
//	d := godump.NewDumper(
//		godump.WithExcludeFields("id"),
//		godump.WithFieldMatchMode(godump.FieldMatchContains),
//	)
//	d.Dump(User{UserID: 10})
//	// #godump.User {
//	// }
func WithFieldMatchMode(mode FieldMatchMode) Option {
	return func(d *Dumper) *Dumper {
		d.fieldMatchMode = mode
		return d
	}
}

// WithRedactFields replaces matching struct fields with a redacted placeholder.
// @group Options
//
// Example: redact fields
//
//	// Default: none
//	type User struct {
//		ID       int
//		Password string
//	}
//	d := godump.NewDumper(
//		godump.WithRedactFields("Password"),
//	)
//	d.Dump(User{ID: 1, Password: "secret"})
//	// #godump.User {
//	//   +ID       => 1 #int
//	//   +Password => <redacted> #string
//	// }
func WithRedactFields(names ...string) Option {
	return func(d *Dumper) *Dumper {
		d.redactFields = append(d.redactFields, names...)
		return d
	}
}

// WithRedactSensitive enables default redaction for common sensitive fields.
// @group Options
//
// Example: redact common sensitive fields
//
//	// Default: disabled
//	type User struct {
//		Password string
//		Token    string
//	}
//	d := godump.NewDumper(
//		godump.WithRedactSensitive(),
//	)
//	d.Dump(User{Password: "secret", Token: "abc"})
//	// #godump.User {
//	//   +Password => <redacted> #string
//	//   +Token    => <redacted> #string
//	// }
func WithRedactSensitive() Option {
	return func(d *Dumper) *Dumper {
		d.redactFields = append(d.redactFields, defaultRedactedFields...)
		d.redactMatchMode = FieldMatchContains
		return d
	}
}

// WithRedactMatchMode sets how field names are matched for WithRedactFields.
// @group Options
//
// Example: use substring matching
//
//	// Default: FieldMatchExact
//	type User struct {
//		APIKey string
//	}
//	d := godump.NewDumper(
//		godump.WithRedactFields("key"),
//		godump.WithRedactMatchMode(godump.FieldMatchContains),
//	)
//	d.Dump(User{APIKey: "abc"})
//	// #godump.User {
//	//   +APIKey => <redacted> #string
//	// }
func WithRedactMatchMode(mode FieldMatchMode) Option {
	return func(d *Dumper) *Dumper {
		d.redactMatchMode = mode
		return d
	}
}

// NewDumper creates a new Dumper with the given options applied.
// Defaults are used for any setting not overridden.
// @group Builder
//
// Example: build a custom dumper
//
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper(
//		godump.WithMaxDepth(10),
//		godump.WithWriter(os.Stdout),
//	)
//	d.Dump(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func NewDumper(opts ...Option) *Dumper {
	d := &Dumper{
		maxDepth:        defaultMaxDepth,
		maxItems:        defaultMaxItems,
		maxStringLen:    defaultMaxStringLen,
		disableStringer: defaultDisableStringer,
		writer:          os.Stdout,
		colorizer:       nil, // ensure no detection is made if we don't need it
		callerFn:        runtime.Caller,
		fieldMatchMode:  FieldMatchExact,
		redactMatchMode: FieldMatchExact,
	}
	for _, opt := range opts {
		d = opt(d)
	}
	return d
}

// Dump prints the values to stdout with colorized output.
// @group Dump
//
// Example: print to stdout
//
//	v := map[string]int{"a": 1}
//	godump.Dump(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func Dump(vs ...any) {
	defaultDumper.Dump(vs...)
}

// Dump prints the values to stdout with colorized output.
// @group Dump
//
// Example: print with a custom dumper
//
//	d := godump.NewDumper()
//	v := map[string]int{"a": 1}
//	d.Dump(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func (d *Dumper) Dump(vs ...any) {
	fmt.Fprint(d.writer, d.DumpStr(vs...))
}

// Fdump writes the formatted dump of values to the given io.Writer.
// @group Dump
//
// Example: dump to writer
//
//	var b strings.Builder
//	v := map[string]int{"a": 1}
//	godump.Fdump(&b, v)
//	// outputs to strings builder
func Fdump(w io.Writer, vs ...any) {
	NewDumper(WithWriter(w)).Dump(vs...)
}

// DumpStr returns a string representation of the values with colorized output.
// @group Dump
//
// Example: get a string dump
//
//	v := map[string]int{"a": 1}
//	out := godump.DumpStr(v)
//	godump.Dump(out)
//	// "#map[string]int {\n  a => 1 #int\n}" #string
func DumpStr(vs ...any) string {
	return defaultDumper.DumpStr(vs...)
}

// DumpStr returns a string representation of the values with colorized output.
// @group Dump
//
// Example: get a string dump with a custom dumper
//
//	d := godump.NewDumper()
//	v := map[string]int{"a": 1}
//	out := d.DumpStr(v)
//	_ = out
//	// "#map[string]int {\n  a => 1 #int\n}" #string
func (d *Dumper) DumpStr(vs ...any) string {
	local := d.clone()
	state := newDumpState()
	var sb strings.Builder
	local.printDumpHeader(&sb)
	tw := tabwriter.NewWriter(&sb, 0, 0, 1, ' ', 0)
	local.writeDump(tw, state, vs...)
	tw.Flush()
	return sb.String()
}

// DumpJSONStr pretty-prints values as JSON and returns it as a string.
// @group JSON
//
// Example: dump JSON string
//
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper()
//	out := d.DumpJSONStr(v)
//	_ = out
//	// {"a":1}
func (d *Dumper) DumpJSONStr(vs ...any) string {
	if len(vs) == 0 {
		return `{"error": "DumpJSON called with no arguments"}`
	}

	var data any = vs
	if len(vs) == 1 {
		data = vs[0]
	}

	b, err := json.MarshalIndent(data, "", strings.Repeat(" ", indentWidth))
	if err != nil {
		//nolint:errchkjson // fallback handles this manually below
		errorJSON, _ := json.Marshal(map[string]string{"error": err.Error()})
		return string(errorJSON)
	}
	return string(b)
}

// DumpJSON prints a pretty-printed JSON string to the configured writer.
// @group JSON
//
// Example: print JSON
//
//	v := map[string]int{"a": 1}
//	d := godump.NewDumper()
//	d.DumpJSON(v)
//	// {
//	//   "a": 1
//	// }
func (d *Dumper) DumpJSON(vs ...any) {
	output := d.DumpJSONStr(vs...)
	fmt.Fprintln(d.writer, output)
}

// DumpHTML dumps the values as HTML with colorized output.
// @group HTML
//
// Example: dump HTML
//
//	v := map[string]int{"a": 1}
//	html := godump.DumpHTML(v)
//	_ = html
//	// (html output)
func DumpHTML(vs ...any) string {
	return defaultDumper.DumpHTML(vs...)
}

// DumpHTML dumps the values as HTML with colorized output.
// @group HTML
//
// Example: dump HTML with a custom dumper
//
//	d := godump.NewDumper()
//	v := map[string]int{"a": 1}
//	html := d.DumpHTML(v)
//	_ = html
//	fmt.Println(html)
//	// (html output)
func (d *Dumper) DumpHTML(vs ...any) string {
	var sb strings.Builder
	sb.WriteString(`<div style='background-color:black;'><pre style="background-color:black; color:white; padding:5px; border-radius: 5px">` + "\n")

	htmlDumper := d.clone()
	if !htmlDumper.disableColor {
		htmlDumper.colorizer = colorizeHTML // use HTML colorizer
	}

	sb.WriteString(htmlDumper.DumpStr(vs...))

	sb.WriteString("</pre></div>")
	return sb.String()
}

// DumpJSON dumps the values as a pretty-printed JSON string.
// If there is more than one value, they are dumped as a JSON array.
// It returns an error string if marshaling fails.
// @group JSON
//
// Example: print JSON
//
//	v := map[string]int{"a": 1}
//	godump.DumpJSON(v)
//	// {
//	//   "a": 1
//	// }
func DumpJSON(vs ...any) {
	defaultDumper.DumpJSON(vs...)
}

// DumpJSONStr dumps the values as a JSON string.
// @group JSON
//
// Example: JSON string
//
//	v := map[string]int{"a": 1}
//	out := godump.DumpJSONStr(v)
//	_ = out
//	// {"a":1}
func DumpJSONStr(vs ...any) string {
	return defaultDumper.DumpJSONStr(vs...)
}

// Dd is a debug function that prints the values and exits the program.
// @group Dump
//
// Example: dump and exit
//
//	v := map[string]int{"a": 1}
//	godump.Dd(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func Dd(vs ...any) {
	defaultDumper.Dd(vs...)
}

// Dd is a debug function that prints the values and exits the program.
// @group Debug
//
// Example: dump and exit with a custom dumper
//
//	d := godump.NewDumper()
//	v := map[string]int{"a": 1}
//	d.Dd(v)
//	// #map[string]int {
//	//   a => 1 #int
//	// }
func (d *Dumper) Dd(vs ...any) {
	d.Dump(vs...)
	exitFunc(1)
}

// clone creates a copy of the [Dumper] with the same configuration.
// This is useful for creating a new dumper with the same settings without modifying the original.
func (d *Dumper) clone() *Dumper {
	newDumper := *d
	return &newDumper
}

// colorize applies the configured [Colorizer] to the string with the given color code.
func (d *Dumper) colorize(code, str string) string {
	if d.colorizer == nil {
		// this avoids detecting color if not needed
		if d.disableColor {
			d.colorizer = colorizeUnstyled
			return d.colorizer(code, str)
		}
		d.colorizer = newColorizer()
	}
	return d.colorizer(code, str)
}

// ensureColorizer initializes the colorizer when none is configured.
func (d *Dumper) ensureColorizer() {
	if d.colorizer == nil {
		if d.disableColor {
			d.colorizer = colorizeUnstyled
			return
		}
		d.colorizer = newColorizer()
	}
}

// printDumpHeader prints the header for the dump output, including the file and line number.
func (d *Dumper) printDumpHeader(out io.Writer) {
	if d.disableHeader {
		return
	}
	file, line := d.findFirstNonInternalFrame(d.skippedStackFrames)
	if file == "" {
		return
	}

	relPath := file
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, file); err == nil {
			relPath = rel
		}
	}

	header := fmt.Sprintf("<#dump // %s:%d", relPath, line)
	fmt.Fprintln(out, d.colorize(colorGray, header))
}

// findFirstNonInternalFrame iterates through the call stack to find the first non-internal frame.
func (d *Dumper) findFirstNonInternalFrame(skip int) (string, int) {
	for i := initialCallerSkip; i < defaultMaxStackDepth; i++ {
		pc, file, line, ok := d.callerFn(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil || !strings.Contains(fn.Name(), "godump") || strings.HasSuffix(file, "_test.go") {
			if skip > 0 {
				skip--
				continue
			}

			return file, line
		}
	}
	return "", 0
}

// formatByteSliceAsHexDump formats a byte slice as a hex dump with ASCII representation.
func (d *Dumper) formatByteSliceAsHexDump(b []byte, indent int) string {
	var sb strings.Builder

	const lineLen = 16
	const asciiStartCol = 50
	const asciiMaxLen = 16

	fieldIndent := strings.Repeat(" ", indent*indentWidth)
	bodyIndent := fieldIndent

	// Header
	sb.WriteString(fmt.Sprintf("([]uint8) (len=%d cap=%d) {\n", len(b), cap(b)))

	for i := 0; i < len(b); i += lineLen {

		end := i + lineLen
		if end > len(b) {
			end = len(b)
		}
		line := b[i:end]

		visibleLen := 0

		// Offset
		offsetStr := fmt.Sprintf("%08x  ", i)
		sb.WriteString(bodyIndent)
		sb.WriteString(d.colorize(colorMeta, offsetStr))
		visibleLen += len(offsetStr)

		// Hex bytes
		for j := 0; j < lineLen; j++ {
			var hexStr string
			if j < len(line) {
				hexStr = fmt.Sprintf("%02x ", line[j])
			} else {
				hexStr = "   "
			}
			if j == 7 {
				hexStr += " "
			}
			sb.WriteString(d.colorize(colorCyan, hexStr))
			visibleLen += len(hexStr)
		}

		// Padding before ASCII
		padding := asciiStartCol - visibleLen
		if padding < 1 {
			padding = 1
		}
		sb.WriteString(strings.Repeat(" ", padding))

		// ASCII section
		sb.WriteString(d.colorize(colorGray, "| "))
		asciiCount := 0
		for _, c := range line {
			ch := "."
			if c >= 32 && c <= 126 {
				ch = string(c)
			}
			sb.WriteString(d.colorize(colorLime, ch))
			asciiCount++
		}
		if asciiCount < asciiMaxLen {
			sb.WriteString(strings.Repeat(" ", asciiMaxLen-asciiCount))
		}
		sb.WriteString(d.colorize(colorGray, " |") + "\n")
	}

	// Closing
	fieldIndent = fieldIndent[:len(fieldIndent)-indentWidth]
	sb.WriteString(fieldIndent + "}")
	return sb.String()
}

func (d *Dumper) writeDump(w io.Writer, state *dumpState, vs ...any) {
	for _, v := range vs {
		rv := reflect.ValueOf(v)
		rv = makeAddressable(rv)
		d.printValue(w, rv, 0, state)
		fmt.Fprintln(w)
	}
}

func (d *Dumper) getTypeString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", d.getTypeString(t.Key()), d.getTypeString(t.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", d.getTypeString(t.Elem()))
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), d.getTypeString(t.Elem()))
	case reflect.Ptr:
		return fmt.Sprintf("*%s", d.getTypeString(t.Elem()))
	default:
		return t.String()
	}
}

func (d *Dumper) printValue(w io.Writer, v reflect.Value, indent int, state *dumpState) {
	if !v.IsValid() {
		fmt.Fprint(w, d.colorize(colorGray, "<invalid>"))
		return
	}

	if isNil(v) {
		typeStr := d.getTypeString(v.Type())
		fmt.Fprintf(w, d.colorize(colorLime, typeStr)+d.colorize(colorGray, "(nil)"))
		return
	}

	if shouldTruncateAtDepth(v, indent, d.maxDepth) {
		fmt.Fprint(w, d.colorize(colorGray, "... (max depth)"))
		return
	}

	if s := d.asStringer(v); s != "" {
		fmt.Fprint(w, s)
		return
	}

	switch v.Kind() {
	case reflect.Chan:
		typ := d.colorizer(colorGray, d.getTypeString(v.Type()))
		fmt.Fprintf(w, "%s(%s)", d.colorize(colorGray, typ), d.colorize(colorCyan, fmt.Sprintf("%#x", v.Pointer())))
		return
	}

	if v.Kind() == reflect.Ptr && v.CanAddr() {
		ptr := v.Pointer()
		if id, ok := state.refs[ptr]; ok {
			fmt.Fprintf(w, d.colorize(colorRef, "↩︎ &%d"), id)
			return
		} else {
			state.refs[ptr] = state.nextRefID
			state.nextRefID++
		}
	}

	// We don't need to check any previous checks (validity, channel, nil,
	// addressable pointer) since they all work directly on the pointer type. We
	// can simply continue with the reference value from here and add a pointer
	// prefix to the output.
	ptrPrefix := ""
	for v.Kind() == reflect.Ptr {
		ptrPrefix += "*"
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Interface:
		d.printValue(w, v.Elem(), indent, state)
	case reflect.Struct:
		t := v.Type()
		fmt.Fprintf(w, "%s {", d.colorize(colorGray, fmt.Sprintf("#%s%s", ptrPrefix, d.getTypeString(v.Type()))))
		fmt.Fprintln(w)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := v.Field(i)
			if !d.shouldIncludeField(field.Name) {
				continue
			}

			symbol := "+"
			if field.PkgPath != "" {
				symbol = "-"
				fieldVal = forceExported(fieldVal)
			}
			indentPrint(w, indent+1, d.colorize(colorYellow, symbol)+field.Name)
			fmt.Fprint(w, "	=> ")
			if d.shouldRedactField(field.Name) {
				fmt.Fprint(w, d.redactedValue(fieldVal))
			} else {
				d.printValue(w, fieldVal, indent+1, state)
			}
			fmt.Fprintln(w)
		}
		indentPrint(w, indent, "")
		fmt.Fprint(w, "}")
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprint(w, d.colorize(colorCyan, fmt.Sprintf("%v", v.Complex())))
	case reflect.UnsafePointer:
		fmt.Fprint(w, d.colorize(colorGray, fmt.Sprintf("unsafe.Pointer(%#x)", v.Pointer())))
	case reflect.Map:
		fmt.Fprintf(w, "%s {", d.colorize(colorGray, fmt.Sprintf("#%s%s", ptrPrefix, d.getTypeString(v.Type()))))
		fmt.Fprintln(w)

		keys := v.MapKeys()
		for i, key := range keys {
			if i >= d.maxItems {
				indentPrint(w, indent+1, d.colorize(colorGray, "... (truncated)"))
				break
			}
			keyStr := fmt.Sprintf("%v", key.Interface())
			indentPrint(w, indent+1, fmt.Sprintf(" %s => ", d.colorize(colorMeta, keyStr)))
			d.printValue(w, v.MapIndex(key), indent+1, state)
			fmt.Fprintln(w)
		}
		indentPrint(w, indent, "")
		fmt.Fprint(w, "}")
	case reflect.Slice, reflect.Array:
		// []byte handling
		if v.Type().Elem().Kind() == reflect.Uint8 {
			if v.CanConvert(reflect.TypeOf([]byte{})) { // Check if it can be converted to []byte
				if data, ok := v.Convert(reflect.TypeOf([]byte{})).Interface().([]byte); ok {
					hexDump := d.formatByteSliceAsHexDump(data, indent+1)
					fmt.Fprint(w, d.colorize(colorLime, hexDump))
					break
				}
			}
		}

		// Default rendering for other slices/arrays
		fmt.Fprintf(w, "%s [", d.colorize(colorGray, fmt.Sprintf("#%s%s", ptrPrefix, d.getTypeString(v.Type()))))
		fmt.Fprintln(w)

		for i := 0; i < v.Len(); i++ {
			if i >= d.maxItems {
				indentPrint(w, indent+1, d.colorize(colorGray, "... (truncated)\n"))
				break
			}
			indentPrint(w, indent+1, fmt.Sprintf("%s => ", d.colorize(colorCyan, fmt.Sprintf("%d", i))))
			d.printValue(w, v.Index(i), indent+1, state)
			fmt.Fprintln(w)
		}
		indentPrint(w, indent, "")
		fmt.Fprint(w, "]")
	case reflect.String:
		str := escapeControl(v.String())
		if utf8.RuneCountInString(str) > d.maxStringLen {
			runes := []rune(str)
			str = string(runes[:d.maxStringLen]) + "…"
		}
		fmt.Fprint(w, d.colorize(colorYellow, `"`)+d.colorize(colorLime, str)+d.colorize(colorYellow, `"`))
	case reflect.Bool:
		if v.Bool() {
			fmt.Fprint(w, d.colorize(colorYellow, "true"))
		} else {
			fmt.Fprint(w, d.colorize(colorGray, "false"))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprint(w, d.colorize(colorCyan, fmt.Sprint(v.Int())))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprint(w, d.colorize(colorCyan, fmt.Sprint(v.Uint())))
	case reflect.Float32, reflect.Float64:
		fmt.Fprint(w, d.colorize(colorCyan, fmt.Sprintf("%f", v.Float())))
	case reflect.Func:
		fmt.Fprint(w, d.colorize(colorGray, v.Type().String()))
	}

	// These types should not have post types since they have a body and already
	// had their type written out.
	if contains([]reflect.Kind{
		reflect.Struct,
		reflect.UnsafePointer,
		reflect.Map,
		reflect.Slice,
		reflect.Array,
		reflect.Ptr,
		reflect.Interface,
	}, v.Kind()) {
		return
	}

	fmt.Fprint(w, d.colorizer(colorGray, fmt.Sprintf(" #%s%s", ptrPrefix, d.getTypeString(v.Type()))))
}

// asStringer checks if the value implements fmt.Stringer and returns its string representation.
func (d *Dumper) asStringer(v reflect.Value) string {
	if d.disableStringer {
		return ""
	}

	val := v
	if !val.CanInterface() {
		val = forceExported(val)
	}
	if val.CanInterface() {
		if s, ok := val.Interface().(fmt.Stringer); ok {
			rv := reflect.ValueOf(s)
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				return d.colorize(colorGray, val.Type().String()+"(nil)")
			}
			return d.colorize(colorLime, s.String()) + d.colorize(colorGray, " #"+d.getTypeString(val.Type()))
		}
	}
	return ""
}

// indentPrint prints indented text to the writer.
func indentPrint(w io.Writer, indent int, text string) {
	fmt.Fprint(w, strings.Repeat(" ", indent*indentWidth)+text)
}

// forceExported returns a value that is guaranteed to be exported, even if it is unexported.
func forceExported(v reflect.Value) reflect.Value {
	if v.CanInterface() {
		return v
	}
	if v.CanAddr() {
		return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	// Final fallback: return original value, even if unexported
	return v
}

// makeAddressable ensures the value is addressable, wrapping structs in pointers if necessary.
func makeAddressable(v reflect.Value) reflect.Value {
	// Already addressable? Do nothing
	if v.CanAddr() {
		return v
	}

	// If it's a struct and not addressable, wrap it in a pointer
	if v.Kind() == reflect.Struct {
		ptr := reflect.New(v.Type())
		ptr.Elem().Set(v)
		return ptr.Elem()
	}

	return v
}

// isNil checks if the value is nil based on its kind.
func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}

// replacer is used to escape control characters in strings.
var replacer = strings.NewReplacer(
	"\n", `\n`,
	"\t", `\t`,
	"\r", `\r`,
	"\v", `\v`,
	"\f", `\f`,
	"\x1b", `\x1b`,
)

// escapeControl escapes control characters in a string for safe display.
func escapeControl(s string) string {
	return replacer.Replace(s)
}

// detectColor checks environment variables to determine if color output should be enabled.
func detectColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	return true
}

// newColorizer picks the appropriate colorizer based on environment overrides.
func newColorizer() Colorizer {
	if detectColor() {
		return colorizeANSI
	}
	return colorizeUnstyled
}

// contains reports whether target exists in the candidates slice.
func contains(candidates []reflect.Kind, target reflect.Kind) bool {
	for _, candidate := range candidates {
		if candidate == target {
			return true
		}
	}

	return false
}

// shouldIncludeField returns true when the field survives include/exclude filtering (include takes precedence).
func (d *Dumper) shouldIncludeField(name string) bool {
	if len(d.includeFields) > 0 && !d.matchesAny(name, d.includeFields, FieldMatchExact) {
		return false
	}
	return !d.matchesAny(name, d.excludeFields, d.fieldMatchMode)
}

// shouldRedactField reports whether the field should be replaced with the redacted placeholder.
func (d *Dumper) shouldRedactField(name string) bool {
	return d.matchesAny(name, d.redactFields, d.redactMatchMode)
}

// matchesAny checks whether name matches any of the candidates using the provided mode.
func (d *Dumper) matchesAny(name string, candidates []string, mode FieldMatchMode) bool {
	if len(candidates) == 0 {
		return false
	}
	nameLower := strings.ToLower(name)
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		candidateLower := strings.ToLower(candidate)
		switch mode {
		case FieldMatchContains:
			if strings.Contains(nameLower, candidateLower) {
				return true
			}
		case FieldMatchPrefix:
			if strings.HasPrefix(nameLower, candidateLower) {
				return true
			}
		case FieldMatchSuffix:
			if strings.HasSuffix(nameLower, candidateLower) {
				return true
			}
		default:
			if strings.EqualFold(name, candidate) {
				return true
			}
		}
	}
	return false
}

func (d *Dumper) redactedValue(v reflect.Value) string {
	if !v.IsValid() {
		return d.colorize(colorRed, "<redacted>")
	}
	typeStr := d.getTypeString(v.Type())
	return d.colorize(colorRed, "<redacted>") + d.colorize(colorGray, " #"+typeStr)
}

// isComplexValue reports whether v unwraps to a struct/map/slice/array.
func isComplexValue(v reflect.Value) bool {
	_, ok := complexBaseKind(v)
	return ok
}

// complexBaseKind unwraps interfaces/pointers, rejects nil, and returns the underlying complex kind if present.
func complexBaseKind(v reflect.Value) (reflect.Kind, bool) {
	if !v.IsValid() {
		return 0, false
	}

	for {
		switch v.Kind() {
		case reflect.Interface:
			if v.IsNil() {
				return 0, false
			}
			v = v.Elem()
		case reflect.Ptr:
			if v.IsNil() {
				return 0, false
			}
			v = v.Elem()
		default:
			switch v.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
				return v.Kind(), true
			default:
				return 0, false
			}
		}
	}
}

// shouldTruncateAtDepth determines whether we should print a truncation placeholder at this depth for complex values.
func shouldTruncateAtDepth(v reflect.Value, indent, maxDepth int) bool {
	if indent < maxDepth {
		return false
	}

	kind, ok := complexBaseKind(v)
	if indent > maxDepth {
		return ok
	}

	if !ok {
		return false
	}

	switch kind {
	case reflect.Map, reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}
