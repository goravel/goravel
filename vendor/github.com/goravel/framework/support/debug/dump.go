package debug

import (
	"fmt"
	"io"
	"os"

	"github.com/goforj/godump"
)

var (
	writer io.Writer = os.Stdout
	dumper           = godump.NewDumper(godump.WithWriter(writer), godump.WithSkipStackFrames(1))
	osExit           = os.Exit
)

// DD is used to display detailed information about variables and then exit the program
func DD(v ...any) {
	dumper.Dump(v...)
	osExit(0)
}

// Dump is used to display detailed information about variables
func Dump(v ...any) {
	dumper.Dump(v...)
}

// DumpHTML is used to display detailed information about variables in HTML format
func DumpHTML(v ...any) {
	output := dumper.DumpHTML(v...)
	_, _ = fmt.Fprintln(writer, output)
}

// DumpJSON is used to display detailed information about variables in JSON format
func DumpJSON(v ...any) {
	dumper.DumpJSON(v...)
}

// FDump is used to display detailed information about variables to the specified io.Writer
func FDump(w io.Writer, v ...any) {
	godump.NewDumper(godump.WithWriter(w), godump.WithSkipStackFrames(1)).Dump(v...)
}

// FDumpHTML is used to display detailed information about variables in HTML format to the specified io.Writer
func FDumpHTML(w io.Writer, v ...any) {
	output := dumper.DumpHTML(v...)
	_, _ = fmt.Fprintln(w, output)
}

// FDumpJSON is used to display detailed information about variables in JSON format to the specified io.Writer
func FDumpJSON(w io.Writer, v ...any) {
	output := dumper.DumpJSONStr(v...)
	_, _ = fmt.Fprintln(w, output)
}

// SDump is used to display detailed information about variables as a string
func SDump(v ...any) string {
	return dumper.DumpStr(v...)
}

// SDumpHTML is used to display detailed information about variables in HTML format as a string
func SDumpHTML(v ...any) string {
	return dumper.DumpHTML(v...)
}

// SDumpJSON is used to display detailed information about variables in JSON format as a string
func SDumpJSON(v ...any) string {
	return dumper.DumpJSONStr(v...)
}
