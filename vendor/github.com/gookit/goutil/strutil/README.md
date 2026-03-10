# String Util

This is a go string operate util package.

- Github: https://github.com/gookit/goutil/strutil
- GoDoc: https://godoc.org/github.com/gookit/goutil/strutil

## Install

```bash
go get github.com/gookit/goutil/strutil
```

## Usage

```go
ss := strutil.ToArray("a,b,c", ",")
// Output: []string{"a", "b", "c"}

ints, err := strutil.ToIntSlice("1,2,3")
// Output: []int{1, 2, 3}
```

## Functions

```go
func AddSlashes(s string) string
func AnyToString(val interface{}, defaultAsErr bool) (str string, err error)
func B32Decode(str string) string
func B32Encode(str string) string
func B64Decode(str string) string
func B64Encode(str string) string
func Base64(str string) string
func Bool(s string) (bool, error)
func Byte2str(b []byte) string
func Byte2string(b []byte) string
func BytePos(s string, bt byte) int
func Camel(s string, sep ...string) string
func CamelCase(s string, sep ...string) string
func Compare(s1, s2, op string) bool
func Cut(s, sep string) (before string, after string, found bool)
func EscapeHTML(s string) string
func EscapeJS(s string) string
func FilterEmail(s string) string
func GenMd5(src interface{}) string
func HasAllSubs(s string, subs []string) bool
func HasOnePrefix(s string, prefixes []string) bool
func HasOneSub(s string, subs []string) bool
func HasPrefix(s string, prefix string) bool
func HasSuffix(s string, suffix string) bool
func Implode(sep string, ss ...string) string
func Indent(s, prefix string) string
func IndentBytes(b, prefix []byte) []byte
func Int(s string) (int, error)
func Int64(s string) int64
func Int64OrErr(s string) (int64, error)
func Int64OrPanic(s string) int64
func IntOrPanic(s string) int
func Ints(s string, sep ...string) []int
func IsAlphaNum(c uint8) bool
func IsAlphabet(char uint8) bool
func IsBlank(s string) bool
func IsBlankBytes(bs []byte) bool
func IsEmpty(s string) bool
func IsEndOf(s, suffix string) bool
func IsNotBlank(s string) bool
func IsNumChar(c byte) bool
func IsNumeric(s string) bool
func IsSpace(c byte) bool
func IsSpaceRune(r rune) bool
func IsStartOf(s, prefix string) bool
func IsStartsOf(s string, prefixes []string) bool
func IsSymbol(r rune) bool
func IsValidUtf8(s string) bool
func IsVersion(s string) bool
func Join(sep string, ss ...string) string
func JoinList(sep string, ss []string) string
func LTrim(s string, cutSet ...string) string
func Lower(s string) string
func LowerFirst(s string) string
func Lowercase(s string) string
func Ltrim(s string, cutSet ...string) string
func MD5(src interface{}) string
func Md5(src interface{}) string
func MicroTimeHexID() string
func MicroTimeID() string
func MustBool(s string) bool
func MustCut(s, sep string) (before string, after string)
func MustInt(s string) int
func MustInt64(s string) int64
func MustString(in interface{}) string
func MustToTime(s string, layouts ...string) time.Time
func NoCaseEq(s, t string) bool
func PadLeft(s, pad string, length int) string
func PadRight(s, pad string, length int) string
func Padding(s, pad string, length int, pos uint8) string
func PrettyJSON(v interface{}) (string, error)
func QuietBool(s string) bool
func QuietInt(s string) int
func QuietInt64(s string) int64
func QuietString(in interface{}) string
func Quote(s string) string
func RTrim(s string, cutSet ...string) string
func RandomBytes(length int) ([]byte, error)
func RandomChars(ln int) string
func RandomCharsV2(ln int) string
func RandomCharsV3(ln int) string
func RandomString(length int) (string, error)
func RenderTemplate(input string, data interface{}, fns template.FuncMap, isFile ...bool) string
func RenderText(input string, data interface{}, fns template.FuncMap, isFile ...bool) string
func Repeat(s string, times int) string
func RepeatBytes(char byte, times int) (chars []byte)
func RepeatRune(char rune, times int) (chars []rune)
func Replaces(str string, pairs map[string]string) string
func Rtrim(s string, cutSet ...string) string
func RuneCount(s string) int
func RuneIsLower(c rune) bool
func RuneIsUpper(c rune) bool
func RuneIsWord(c rune) bool
func RunePos(s string, ru rune) int
func RuneWidth(r rune) int
func Similarity(s, t string, rate float32) (float32, bool)
func SnakeCase(s string, sep ...string) string
func Split(s, sep string) (ss []string)
func SplitInlineComment(val string) (string, string)
func SplitN(s, sep string, n int) (ss []string)
func SplitNTrimmed(s, sep string, n int) (ss []string)
func SplitNValid(s, sep string, n int) (ss []string)
func SplitTrimmed(s, sep string) (ss []string)
func SplitValid(s, sep string) (ss []string)
func StrPos(s, sub string) int
func String(val interface{}) (string, error)
func StringOrErr(val interface{}) (string, error)
func Strings(s string, sep ...string) []string
func StripSlashes(s string) string
func Substr(s string, pos, length int) string
func TextSplit(s string, w int) []string
func TextTruncate(s string, w int, tail string) string
func TextWidth(s string) int
func TextWrap(s string, w int) string
func Title(s string) string
func ToArray(s string, sep ...string) []string
func ToBool(s string) (bool, error)
func ToBytes(s string) (b []byte)
func ToDuration(s string) (time.Duration, error)
func ToInt(s string) (int, error)
func ToInt64(s string) (int64, error)
func ToIntSlice(s string, sep ...string) (ints []int, err error)
func ToInts(s string, sep ...string) ([]int, error)
func ToSlice(s string, sep ...string) []string
func ToString(val interface{}) (string, error)
func ToStrings(s string, sep ...string) []string
func ToTime(s string, layouts ...string) (t time.Time, err error)
func Trim(s string, cutSet ...string) string
func TrimCut(s, sep string) (string, string)
func TrimLeft(s string, cutSet ...string) string
func TrimRight(s string, cutSet ...string) string
func URLDecode(s string) string
func URLEncode(s string) string
func Unquote(s string) string
func Upper(s string) string
func UpperFirst(s string) string
func UpperWord(s string) string
func Uppercase(s string) string
func Utf8Len(s string) int
func Utf8Split(s string, w int) []string
func Utf8Truncate(s string, w int, tail string) string
func Utf8Width(s string) (size int)
func Utf8len(s string) int
func VersionCompare(v1, v2, op string) bool
func WidthWrap(s string, w int) string
func WrapTag(s, tag string) string
```

## Code Check & Testing

```bash
gofmt -w -l ./
golint ./...
```

**Testing**:

```shell
go test -v ./strutil/...
```

**Test limit by regexp**:

```shell
go test -v -run ^TestSetByKeys ./strutil/...
```
