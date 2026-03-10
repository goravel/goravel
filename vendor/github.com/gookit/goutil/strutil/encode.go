package strutil

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"net/url"
	"strings"
	"text/template"
)

//
// -------------------- escape --------------------
//

// EscapeJS escape javascript string
func EscapeJS(s string) string {
	return template.JSEscapeString(s)
}

// EscapeHTML escape html string
func EscapeHTML(s string) string {
	return template.HTMLEscapeString(s)
}

// AddSlashes add slashes for the string.
func AddSlashes(s string) string {
	if ln := len(s); ln == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, char := range s {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}

	return buf.String()
}

// StripSlashes strip slashes for the string.
func StripSlashes(s string) string {
	ln := len(s)
	if ln == 0 {
		return ""
	}

	var skip bool
	var buf bytes.Buffer

	for i, char := range s {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < ln && s[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}

	return buf.String()
}

//
// -------------------- encode --------------------
//

// URLEncode encode url string.
func URLEncode(s string) string {
	if pos := strings.IndexRune(s, '?'); pos > -1 { // escape query data
		return s[0:pos+1] + url.QueryEscape(s[pos+1:])
	}
	return s
}

// URLDecode decode url string.
func URLDecode(s string) string {
	if pos := strings.IndexRune(s, '?'); pos > -1 { // un-escape query data
		qy, err := url.QueryUnescape(s[pos+1:])
		if err == nil {
			return s[0:pos+1] + qy
		}
	}

	return s
}

//
// -------------------- base encode --------------------
//

// base32 encoding with no padding
var (
	B32Std = base32.StdEncoding.WithPadding(base32.NoPadding)
	B32Hex = base32.HexEncoding.WithPadding(base32.NoPadding)
)

// B32Encode base32 encode
func B32Encode(str string) string {
	return B32Std.EncodeToString([]byte(str))
}

// B32Decode base32 decode
func B32Decode(str string) string {
	dec, _ := B32Std.DecodeString(str)
	return string(dec)
}

// B64Std base64 encoding with no padding
var B64Std = base64.StdEncoding.WithPadding(base64.NoPadding)

// B64Encode base64 encode
func B64Encode(str string) string {
	return B64Std.EncodeToString([]byte(str))
}

// B64EncodeBytes base64 encode
func B64EncodeBytes(src []byte) []byte {
	buf := make([]byte, B64Std.EncodedLen(len(src)))
	B64Std.Encode(buf, src)
	return buf
}

// B64Decode base64 decode
func B64Decode(str string) string {
	dec, _ := B64Std.DecodeString(str)
	return string(dec)
}

// B64DecodeBytes base64 decode
func B64DecodeBytes(str []byte) []byte {
	dbuf := make([]byte, B64Std.DecodedLen(len(str)))
	n, _ := B64Std.Decode(dbuf, str)
	return dbuf[:n]
}
