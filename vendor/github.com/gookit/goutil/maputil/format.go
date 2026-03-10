package maputil

import (
	"io"
	"reflect"

	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
)

// MapFormatter struct
type MapFormatter struct {
	comdef.BaseFormatter
	// Prefix string for each element
	Prefix string
	// Indent string for each element
	Indent string
	// ClosePrefix string for last "}"
	ClosePrefix string
	// AfterReset after reset on call Format().
	// AfterReset bool
}

// NewFormatter instance
func NewFormatter(mp any) *MapFormatter {
	f := &MapFormatter{}
	f.Src = mp

	return f
}

// WithFn for config self
func (f *MapFormatter) WithFn(fn func(f *MapFormatter)) *MapFormatter {
	fn(f)
	return f
}

// WithIndent string
func (f *MapFormatter) WithIndent(indent string) *MapFormatter {
	f.Indent = indent
	return f
}

// FormatTo to custom buffer
func (f *MapFormatter) FormatTo(w io.Writer) {
	f.SetOutput(w)
	f.doFormat()
}

// Format to string
func (f *MapFormatter) String() string {
	return f.Format()
}

// Format to string
func (f *MapFormatter) Format() string {
	f.doFormat()
	return f.BsWriter().String()
}

// Format map data to string.
//
//goland:noinspection GoUnhandledErrorResult
func (f *MapFormatter) doFormat() {
	if f.Src == nil {
		return
	}

	rv, ok := f.Src.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(f.Src)
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Map {
		return
	}

	buf := f.BsWriter()
	ln := rv.Len()
	if ln == 0 {
		buf.WriteString("{}")
		return
	}

	// buf.Grow(ln * 16)
	buf.WriteByte('{')

	indentLn := len(f.Indent)
	if indentLn > 0 {
		buf.WriteByte('\n')
	}

	for i, key := range rv.MapKeys() {
		strK := strutil.SafeString(key.Interface())
		if indentLn > 0 {
			buf.WriteString(f.Indent)
		}

		buf.WriteString(strK)
		buf.WriteByte(':')

		strV := strutil.SafeString(rv.MapIndex(key).Interface())
		buf.WriteString(strV)
		if i < ln-1 {
			buf.WriteByte(',')

			// no indent, with space
			if indentLn == 0 {
				buf.WriteByte(' ')
			}
		}

		// with newline
		if indentLn > 0 {
			buf.WriteByte('\n')
		}
	}

	if f.ClosePrefix != "" {
		buf.WriteString(f.ClosePrefix)
	}

	buf.WriteByte('}')
}
