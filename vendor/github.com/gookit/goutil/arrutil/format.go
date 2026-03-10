package arrutil

import (
	"io"
	"reflect"

	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
)

// ArrFormatter struct
type ArrFormatter struct {
	comdef.BaseFormatter
	// Prefix string for each element
	Prefix string
	// Indent string for format each element
	Indent string
	// ClosePrefix on before end char: ]
	ClosePrefix string
}

// NewFormatter instance
func NewFormatter(arr any) *ArrFormatter {
	f := &ArrFormatter{}
	f.Src = arr
	return f
}

// FormatIndent array data to string.
func FormatIndent(arr any, indent string) string {
	return NewFormatter(arr).WithIndent(indent).Format()
}

// WithFn for config self
func (f *ArrFormatter) WithFn(fn func(f *ArrFormatter)) *ArrFormatter {
	fn(f)
	return f
}

// WithIndent string
func (f *ArrFormatter) WithIndent(indent string) *ArrFormatter {
	f.Indent = indent
	return f
}

// FormatTo to custom buffer
func (f *ArrFormatter) FormatTo(w io.Writer) {
	f.SetOutput(w)
	f.doFormat()
}

// Format to string
func (f *ArrFormatter) String() string {
	return f.Format()
}

// Format to string
func (f *ArrFormatter) Format() string {
	f.doFormat()
	return f.BsWriter().String()
}

// Format to string
//
//goland:noinspection GoUnhandledErrorResult
func (f *ArrFormatter) doFormat() {
	if f.Src == nil {
		return
	}

	rv, ok := f.Src.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(f.Src)
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return
	}

	writer := f.BsWriter()
	arrLn := rv.Len()
	if arrLn == 0 {
		writer.WriteString("[]")
		return
	}

	// if f.AfterReset {
	// 	defer f.Reset()
	// }

	// sb.Grow(arrLn * 4)
	writer.WriteByte('[')

	indentLn := len(f.Indent)
	if indentLn > 0 {
		writer.WriteByte('\n')
	}

	for i := 0; i < arrLn; i++ {
		if indentLn > 0 {
			writer.WriteString(f.Indent)
		}
		writer.WriteString(strutil.QuietString(rv.Index(i).Interface()))

		if i < arrLn-1 {
			writer.WriteByte(',')

			// no indent, with space
			if indentLn == 0 {
				writer.WriteByte(' ')
			}
		}
		if indentLn > 0 {
			writer.WriteByte('\n')
		}
	}

	if f.ClosePrefix != "" {
		writer.WriteString(f.ClosePrefix)
	}
	writer.WriteByte(']')
}
