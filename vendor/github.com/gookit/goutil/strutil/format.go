package strutil

import (
	"regexp"
	"strings"
	"unicode"
)

/*************************************************************
 * change string case
 *************************************************************/

// methods aliases
var (
	UpWords = UpperWord
	LoFirst = LowerFirst
	UpFirst = UpperFirst

	Snake = SnakeCase
)

// Title alias of the strings.ToTitle()
func Title(s string) string { return strings.ToTitle(s) }

// Lower alias of the strings.ToLower()
func Lower(s string) string { return strings.ToLower(s) }

// Lowercase alias of the strings.ToLower()
func Lowercase(s string) string { return strings.ToLower(s) }

// Upper alias of the strings.ToUpper()
func Upper(s string) string { return strings.ToUpper(s) }

// Uppercase alias of the strings.ToUpper()
func Uppercase(s string) string { return strings.ToUpper(s) }

// UpperWord Change the first character of each word to uppercase
func UpperWord(s string) string {
	if len(s) == 0 {
		return s
	}

	if len(s) == 1 {
		return strings.ToUpper(s)
	}

	inWord := true
	buf := make([]byte, 0, len(s))

	i := 0
	rs := []rune(s)
	if RuneIsLower(rs[i]) {
		buf = append(buf, []byte(string(unicode.ToUpper(rs[i])))...)
	} else {
		buf = append(buf, []byte(string(rs[i]))...)
	}

	for j := i + 1; j < len(rs); j++ {
		if !RuneIsWord(rs[i]) && RuneIsWord(rs[j]) {
			inWord = false
		}

		if RuneIsLower(rs[j]) && !inWord {
			buf = append(buf, []byte(string(unicode.ToUpper(rs[j])))...)
			inWord = true
		} else {
			buf = append(buf, []byte(string(rs[j]))...)
		}

		if RuneIsWord(rs[j]) {
			inWord = true
		}

		i++
	}

	return string(buf)
}

// LowerFirst lower first char
func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	rs := []rune(s)
	f := rs[0]

	if 'A' <= f && f <= 'Z' {
		return string(unicode.ToLower(f)) + string(rs[1:])
	}
	return s
}

// UpperFirst upper first char
func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	rs := []rune(s)
	f := rs[0]

	if 'a' <= f && f <= 'z' {
		return string(unicode.ToUpper(f)) + string(rs[1:])
	}
	return s
}

// SnakeCase convert. eg "RangePrice" -> "range_price"
func SnakeCase(s string, sep ...string) string {
	sepChar := "_"
	if len(sep) > 0 {
		sepChar = sep[0]
	}

	str := toSnakeReg.ReplaceAllStringFunc(s, func(s string) string {
		return sepChar + LowerFirst(s)
	})

	return strings.TrimLeft(str, sepChar)
}

// Camel alias of the CamelCase
func Camel(s string, sep ...string) string { return CamelCase(s, sep...) }

// CamelCase convert string to camel case.
//
// Support:
//
//	"range_price" -> "rangePrice"
//	"range price" -> "rangePrice"
//	"range-price" -> "rangePrice"
func CamelCase(s string, sep ...string) string {
	sepChar := "_"
	if len(sep) > 0 {
		sepChar = sep[0]
	}

	// Not contains sep char
	if !strings.Contains(s, sepChar) {
		return s
	}

	// Get regexp instance
	rgx, ok := toCamelRegs[sepChar]
	if !ok {
		rgx = regexp.MustCompile(regexp.QuoteMeta(sepChar) + "+[a-zA-Z]")
	}

	return rgx.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.TrimLeft(s, sepChar)
		return UpperFirst(s)
	})
}

//
// Indent format multi line text
// from package: github.com/kr/text
//

// Indent inserts prefix at the beginning of each non-empty line of s. The
// end-of-line marker is NL.
func Indent(s, prefix string) string {
	return string(IndentBytes([]byte(s), []byte(prefix)))
}

// IndentBytes inserts prefix at the beginning of each non-empty line of b.
// The end-of-line marker is NL.
func IndentBytes(b, prefix []byte) []byte {
	if len(b) == 0 {
		return b
	}

	bol := true
	res := make([]byte, 0, len(b)+len(prefix)*4)

	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	return res
}
