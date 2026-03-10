package validate

import (
	"bytes"
	"encoding/json"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

// Basic regular expressions for validating strings.
// (there are from package "asaskevich/govalidator")
const (
	Email        = `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	UUID3        = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4        = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5        = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID         = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	Int          = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float        = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	RGBColor     = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	FullWidth    = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	HalfWidth    = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	Base64       = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	Latitude     = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	Longitude    = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	DNSName      = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	FullURL      = `^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=?&-_%]+$`
	URLSchema    = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername  = `(\S+(:\S*)?@)`
	URLPath      = `((\/|\?|#)[^\s]*)`
	URLPort      = `(:(\d{1,5}))`
	URLIP        = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	URLSubdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	WinPath      = `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	UnixPath     = `^(/[^/\x00]*)+/?$`
)

// some string regexp. (it is from package "asaskevich/govalidator")
var (
	// rxUser           = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	// rxHostname       = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	// rxUserDot        = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	rxEmail     = regexp.MustCompile(Email)
	rxISBN10    = regexp.MustCompile(`^(?:\d{9}X|\d{10})$`)
	rxISBN13    = regexp.MustCompile(`^\d{13}$`)
	rxUUID3     = regexp.MustCompile(UUID3)
	rxUUID4     = regexp.MustCompile(UUID4)
	rxUUID5     = regexp.MustCompile(UUID5)
	rxUUID      = regexp.MustCompile(UUID)
	rxAlpha     = regexp.MustCompile("^[a-zA-Z]+$")
	rxAlphaNum  = regexp.MustCompile("^[a-zA-Z0-9]+$")
	rxAlphaDash = regexp.MustCompile(`^(?:[\w-]+)$`)
	rxNumber    = regexp.MustCompile("^[0-9]+$")
	rxInt       = regexp.MustCompile(Int)
	rxFloat     = regexp.MustCompile(Float)
	rxCnMobile  = regexp.MustCompile(`^1\d{10}$`)
	rxHexColor  = regexp.MustCompile(`^#?([\da-fA-F]{3}|[\da-fA-F]{6})$`)
	rxRGBColor  = regexp.MustCompile(RGBColor)
	rxASCII     = regexp.MustCompile("^[\x00-\x7F]+$")
	// --
	rxHexadecimal    = regexp.MustCompile(`^[\da-fA-F]+$`)
	rxPrintableASCII = regexp.MustCompile("^[\x20-\x7E]+$")
	rxMultiByte      = regexp.MustCompile("[^\x00-\x7F]")
	// rxFullWidth = regexp.MustCompile(FullWidth)
	// rxHalfWidth = regexp.MustCompile(HalfWidth)
	rxBase64    = regexp.MustCompile(Base64)
	rxDataURI   = regexp.MustCompile(`^data:.+/(.+);base64,(?:.+)`)
	rxLatitude  = regexp.MustCompile(Latitude)
	rxLongitude = regexp.MustCompile(Longitude)
	rxDNSName   = regexp.MustCompile(DNSName)
	rxFullURL   = regexp.MustCompile(FullURL)
	rxURLSchema = regexp.MustCompile(URLSchema)
	// rxSSN            = regexp.MustCompile(`^\d{3}[- ]?\d{2}[- ]?\d{4}$`)
	rxWinPath  = regexp.MustCompile(WinPath)
	rxUnixPath = regexp.MustCompile(UnixPath)
	// --
	rxHasLowerCase = regexp.MustCompile(".*[[:lower:]]")
	rxHasUpperCase = regexp.MustCompile(".*[[:upper:]]")
)

/*************************************************************
 * global validators
 *************************************************************/

type funcMeta struct {
	fv reflect.Value
	// validator name
	name string
	// readonly cache
	numIn  int
	numOut int
	// is an internal built-in validator
	builtin bool
	// the last arg is variadic param. like "... any"
	isVariadic bool
}

func (fm *funcMeta) checkArgNum(argNum int, name string) {
	// last arg is like "... any"
	if fm.isVariadic {
		if argNum+1 < fm.numIn {
			panicf("not enough parameters for validator '%s'!", name)
		}
	} else if argNum != fm.numIn {
		panicf(
			"the number of parameters given does not match the required. validator '%s', want %d, given %d",
			name,
			fm.numIn,
			argNum,
		)
	}
}

func newFuncMeta(name string, builtin bool, fv reflect.Value) *funcMeta {
	fm := &funcMeta{fv: fv, name: name, builtin: builtin}
	ft := fv.Type()

	fm.numIn = ft.NumIn()   // arg num of the func
	fm.numOut = ft.NumOut() // return arg num of the func
	fm.isVariadic = ft.IsVariadic()

	return fm
}

// ValidatorName get real validator name.
func ValidatorName(name string) string {
	if rName, ok := validatorAliases[name]; ok {
		return rName
	}
	return name
}

// AddValidators to the global validators map
func AddValidators(m map[string]any) {
	for name, checkFunc := range m {
		AddValidator(name, checkFunc)
	}
}

// AddValidator to the pkg. checkFunc must return a bool
//
// Usage:
//
//	v.AddValidator("myFunc", func(val any) bool {
//		// do validate val ...
//		return true
//	})
func AddValidator(name string, checkFunc any) {
	fv := checkValidatorFunc(name, checkFunc)

	validators[name] = validatorTypeCustom
	// validatorValues[name] = fv
	validatorMetas[name] = newFuncMeta(name, false, fv)
}

// Validators get all validator names
func Validators() map[string]int8 {
	return validators
}

/*************************************************************
 * context validators:
 *  - field value compare
 *  - requiredXXX
 *************************************************************/

// Required field val check
func (v *Validation) Required(field string, val any) bool {
	if v.isInOptional(field) {
		return true
	}

	if v.data != nil && v.data.Type() == sourceForm {
		// check is upload file
		if v.data.(*FormData).HasFile(field) {
			return true
		}
	}

	if v.isIgnoreableZeroNumeric(field) {
		return true
	}

	// check value
	return !IsEmpty(val)
}

// RequiredIf field under validation must be present and not empty,
// if the anotherField field is equal to any value.
//
// Usage:
//
//	v.AddRule("password", "requiredIf", "username", "tom")
func (v *Validation) RequiredIf(sourceField string, val any, kvs ...string) bool {
	if len(kvs) < 2 {
		return false
	}

	dstField, args := kvs[0], kvs[1:]
	if dstVal, has := v.Get(dstField); has {
		// up: only one check value, direct compare value
		if len(args) == 1 {
			rftDv := reflect.ValueOf(dstVal)
			wantVal, err := convTypeByBaseKind(args[0], rftDv.Kind())
			if err == nil && dstVal == wantVal {
				return val != nil && !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
			}
		} else if Enum(dstVal, args) {
			return val != nil && !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
		}
	}

	// default as True, skip check
	return true
}

// RequiredUnless field under validation must be present and not empty
// unless the dstField field is equal to any value.
//
//   - kvs format: [dstField, dstVal1, dstVal2 ...]
func (v *Validation) RequiredUnless(sourceField string, val any, kvs ...string) bool {
	if len(kvs) < 2 {
		return false
	}

	dstField, values := kvs[0], kvs[1:]
	if dstVal, has, _ := v.tryGet(dstField); has {
		if !Enum(dstVal, values) {
			return !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
		}
	}

	// fields in values
	return true
}

// RequiredWith field under validation must be present and not empty only
// if any of the other specified fields are present.
//
//   - fields format: [field1, field2 ...]
func (v *Validation) RequiredWith(sourceField string, val any, fields ...string) bool {
	if len(fields) == 0 {
		return false
	}

	for _, field := range fields {
		if _, has, zero := v.tryGet(field); has && !zero {
			return !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
		}
	}

	// all fields not exist
	return true
}

// RequiredWithAll field under validation must be present and not empty only if all the other specified fields are present.
func (v *Validation) RequiredWithAll(sourceField string, val any, fields ...string) bool {
	if len(fields) == 0 {
		return false
	}

	for _, field := range fields {
		if _, has, zero := v.tryGet(field); !has || zero {
			// if any field does not exist, not continue.
			return true
		}
	}

	// all fields exist
	return !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
}

// RequiredWithout field under validation must be present and not empty only when any of the other specified fields are not present.
func (v *Validation) RequiredWithout(sourceField string, val any, fields ...string) bool {
	if len(fields) == 0 {
		return false
	}

	for _, field := range fields {
		if _, has, zero := v.tryGet(field); !has || zero {
			return !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
		}
	}

	// all fields exist
	return true
}

// RequiredWithoutAll field under validation must be present and not empty only when any of the other specified fields are not present.
func (v *Validation) RequiredWithoutAll(sourceField string, val any, fields ...string) bool {
	if len(fields) == 0 {
		return false
	}

	for _, name := range fields {
		// if any field exists, not continue.
		if _, has, zero := v.tryGet(name); has && !zero {
			return true
		}
	}

	// all fields not exists, required
	return !IsEmpty(val) || v.isIgnoreableZeroNumeric(sourceField)
}

// EqField value should EQ the dst field value
func (v *Validation) EqField(val any, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return IsEqual(val, dstVal)
}

// NeField value should not equal the dst field value
func (v *Validation) NeField(val any, dstField string) bool {
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return !IsEqual(val, dstVal)
}

// GtField value should GT the dst field value
func (v *Validation) GtField(val any, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return valueCompare(val, dstVal, ">")
}

// GteField value should GTE the dst field value
func (v *Validation) GteField(val any, dstField string) bool {
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return valueCompare(val, dstVal, ">=")
}

// LtField value should LT the dst field value
func (v *Validation) LtField(val any, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return valueCompare(val, dstVal, "<")
}

// LteField value should LTE the dst field value(for int, string)
func (v *Validation) LteField(val any, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return valueCompare(val, dstVal, "<=")
}

/*
 ******************************************************************
 * context validators:
 *  - file validators
 ******************************************************************
 */

const fileValidators = "|isFile|isImage|inMimeTypes|"

var (
	imageMimeTypes = map[string]string{
		"bmp": "image/bmp",
		"gif": "image/gif",
		"ief": "image/ief",
		"jpg": "image/jpeg",
		// "jpe":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
		"webp": "image/webp",
	}
)

func isFileValidator(name string) bool {
	return strings.Contains(fileValidators, "|"+name+"|")
}

// IsFormFile check field is uploaded file. validator: isFile
func (v *Validation) IsFormFile(fd *FormData, field string) (ok bool) {
	field, _, _ = strings.Cut(field, ".*")
	if files := fd.GetFiles(field); len(files) > 0 {
		for i := range files {
			_, err := files[i].Open()
			if err != nil {
				return false
			}
		}
		return true
	}
	return false
}

// IsFormImage check field is uploaded image file. validator: isImage
//
// Usage:
//
//	v.AddRule("avatar", "image")
//	v.AddRule("avatar", "image", "jpg", "png", "gif") // set ext limit
//	v.AddRule("images.*", "image")
//	v.AddRule("images.*", "image", "jpg", "png", "gif") // set ext limit
func (v *Validation) IsFormImage(fd *FormData, field string, exts ...string) (ok bool) {
	field, _, expectArray := strings.Cut(field, ".*")
	if expectArray {
		for _, mime := range fd.FilesMimeType(field) {
			if !v.isImageMimeTypes(mime, exts...) {
				return false
			}
		}
		return true
	}

	return v.isImageMimeTypes(fd.FileMimeType(field), exts...)
}

func (v *Validation) isImageMimeTypes(mime string, exts ...string) (ok bool) {
	if mime == "" {
		return
	}

	var fileExt string
	for ext, imgMime := range imageMimeTypes {
		if imgMime == mime {
			fileExt = ext
			ok = true
			break
		}
	}

	// don't limit mime type
	if len(exts) == 0 {
		return ok // only check is an image
	}
	return Enum(fileExt, exts)
}

// InMimeTypes check field is uploaded file and mimetype is in the mimeTypes. validator: inMimeTypes
//
// Usage:
//
//	v.AddRule("video", "mimeTypes", "video/avi", "video/mpeg", "video/quicktime")
//	v.AddRule("videos.*", "mimeTypes", "video/avi", "video/mpeg", "video/quicktime")
func (v *Validation) InMimeTypes(fd *FormData, field, mimeType string, moreTypes ...string) bool {
	field, _, expectArray := strings.Cut(field, ".*")
	mimeTypes := append(moreTypes, mimeType) //nolint:gocritic
	if expectArray {
		for _, mime := range fd.FilesMimeType(field) {
			if !v.inMimeTypes(mime, mimeTypes) {
				return false
			}
		}
		return true
	}

	return v.inMimeTypes(fd.FileMimeType(field), mimeTypes)
}

func (v *Validation) inMimeTypes(mime string, mimeTypes []string) bool {
	if mime == "" {
		return false
	}
	return Enum(mime, mimeTypes)
}

func (v *Validation) isIgnoreableZeroNumeric(field string) bool {
	if v.data != nil && v.data.Type() == sourceMap {
		if val, ok := v.data.Get(field); ok {
			k := reflect.ValueOf(val).Kind()
			return k >= reflect.Int && k <= reflect.Float64
		}
	}
	return false
}

/*************************************************************
 * global: basic validators
 *************************************************************/

// IsEmpty of the value
func IsEmpty(val any) bool {
	if val == nil {
		return true
	}
	if s, ok := val.(string); ok {
		return s == ""
	}

	var rv reflect.Value

	// type check val is reflect.Value
	if v2, ok := val.(reflect.Value); ok {
		rv = v2
	} else {
		rv = reflect.ValueOf(val)
	}
	return ValueIsEmpty(rv)
}

// Contains check that the specified string, list(array, slice) or map contains the
// specified substring or element.
//
// Notice: list check value exist. map check key exist.
func Contains(s, sub any) bool {
	ok, found := includeElement(s, sub)

	// ok == false: 's' could not be applied builtin len()
	// found == false: 's' does not contain 'sub'
	return ok && found
}

// NotContains check that the specified string, list(array, slice) or map does NOT contain the
// specified substring or element.
//
// Notice: list check value exist. map check key exist.
func NotContains(s, sub any) bool {
	ok, found := includeElement(s, sub)

	// ok == false: could not be applied builtin len()
	// found == true: 's' contain 'sub'
	return ok && !found
}

/*************************************************************
 * global: type validators
 *************************************************************/

// IsUint check, allow: intX, uintX, string
func IsUint(val any) bool {
	switch typVal := val.(type) {
	case int:
		return typVal >= 0
	case int8:
		return typVal >= 0
	case int16:
		return typVal >= 0
	case int32:
		return typVal >= 0
	case int64:
		return typVal >= 0
	case uint, uint8, uint16, uint32, uint64:
		return true
	case string:
		_, err := strconv.ParseUint(typVal, 10, 32)
		return err == nil
	}
	return false
}

// IsBool check. allow: bool, string.
func IsBool(val any) bool {
	val = indirectValue(val)

	if _, ok := val.(bool); ok {
		return true
	}

	if typVal, ok := val.(string); ok {
		_, err := strutil.ToBool(typVal)
		return err == nil
	}
	return false
}

// IsFloat check. allow: floatX, string
func IsFloat(val any) bool {
	val = indirectValue(val)

	if val == nil {
		return false
	}

	switch rv := val.(type) {
	case float32, float64:
		return true
	case string:
		return rv != "" && rxFloat.MatchString(rv)
	}
	return false
}

// IsArray check value is array or slice.
func IsArray(val any, strict ...bool) (ok bool) {
	if val == nil {
		return false
	}

	rv := reflect.Indirect(reflect.ValueOf(val))

	// strict: must go array type.
	if len(strict) > 0 && strict[0] {
		return rv.Kind() == reflect.Array
	}

	// allow array, slice
	return rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice
}

// IsSlice check value is slice type
func IsSlice(val any) (ok bool) {
	if val == nil {
		return false
	}

	rv := reflect.Indirect(reflect.ValueOf(val))
	return rv.Kind() == reflect.Slice
}

// IsInts is int slice check
func IsInts(val any) bool {
	if val == nil {
		return false
	}

	switch val.(type) {
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
		return true
	}
	return false
}

// IsStrings is string slice check
func IsStrings(val any) (ok bool) {
	if val == nil {
		return false
	}

	_, ok = val.([]string)
	return
}

// IsMap check
func IsMap(val any) (ok bool) {
	if val == nil {
		return false
	}

	rv := reflect.Indirect(reflect.ValueOf(val))
	return rv.Kind() == reflect.Map
}

// IsInt check, and support length check
func IsInt(val any, minAndMax ...int64) (ok bool) {
	if val == nil {
		return false
	}
	val = indirectValue(val)

	intVal, err := valueToInt64(val, true)
	if err != nil {
		return false
	}

	argLn := len(minAndMax)
	if argLn == 0 { // only check type
		return true
	}

	// value check
	minVal := minAndMax[0]
	if argLn == 1 { // only min length check.
		return intVal >= minVal
	}

	// min and max length check
	return intVal >= minVal && intVal <= minAndMax[1]
}

// IsString check and support length check.
//
// Usage:
//
//	ok := IsString(val)
//	ok := IsString(val, 5) // with min len check
//	ok := IsString(val, 5, 12) // with min and max len check
func IsString(val any, minAndMaxLen ...int) (ok bool) {
	val = indirectValue(val)

	if val == nil {
		return false
	}

	argLn := len(minAndMaxLen)
	str, isStr := val.(string)

	// only check type
	if argLn == 0 {
		return isStr
	}

	if !isStr {
		return false
	}

	// length check
	strLen := len(str)
	minLen := minAndMaxLen[0]

	// only min length check.
	if argLn == 1 {
		return strLen >= minLen
	}

	// min and max length check
	return strLen >= minLen && strLen <= minAndMaxLen[1]
}

/*************************************************************
 * global: string validators
 *************************************************************/

// HasWhitespace check. eg "10"
func HasWhitespace(s string) bool {
	return s != "" && strings.ContainsRune(s, ' ')
}

// IsIntString check. eg "10"
func IsIntString(s string) bool { return s != "" && rxInt.MatchString(s) }

// IsASCII string.
func IsASCII(s string) bool { return s != "" && rxASCII.MatchString(s) }

// IsPrintableASCII string.
func IsPrintableASCII(s string) bool {
	return s != "" && rxPrintableASCII.MatchString(s)
}

// IsBase64 string.
func IsBase64(s string) bool { return s != "" && rxBase64.MatchString(s) }

// IsLatitude string.
func IsLatitude(s string) bool { return s != "" && rxLatitude.MatchString(s) }

// IsLongitude string.
func IsLongitude(s string) bool { return s != "" && rxLongitude.MatchString(s) }

// IsDNSName string.
func IsDNSName(s string) bool { return s != "" && rxDNSName.MatchString(s) }

// HasURLSchema string.
func HasURLSchema(s string) bool { return s != "" && rxURLSchema.MatchString(s) }

// IsFullURL string.
func IsFullURL(s string) bool { return s != "" && rxFullURL.MatchString(s) }

// IsURL string.
func IsURL(s string) bool {
	if s == "" {
		return false
	}

	_, err := url.Parse(s)
	return err == nil
}

// IsDataURI string.
//
// data:[<mime type>] ( [;charset=<charset>] ) [;base64],码内容
// eg. "data:image/gif;base64,R0lGODlhA..."
func IsDataURI(s string) bool { return s != "" && rxDataURI.MatchString(s) }

// IsMultiByte string.
func IsMultiByte(s string) bool { return s != "" && rxMultiByte.MatchString(s) }

// IsISBN10 string.
func IsISBN10(s string) bool { return s != "" && rxISBN10.MatchString(s) }

// IsISBN13 string.
func IsISBN13(s string) bool { return s != "" && rxISBN13.MatchString(s) }

// IsHexadecimal string.
func IsHexadecimal(s string) bool { return s != "" && rxHexadecimal.MatchString(s) }

// IsCnMobile string.
func IsCnMobile(s string) bool { return s != "" && rxCnMobile.MatchString(s) }

// IsHexColor string.
func IsHexColor(s string) bool { return s != "" && rxHexColor.MatchString(s) }

// IsRGBColor string.
func IsRGBColor(s string) bool { return s != "" && rxRGBColor.MatchString(s) }

// IsAlpha string.
func IsAlpha(s string) bool { return s != "" && rxAlpha.MatchString(s) }

// IsAlphaNum string.
func IsAlphaNum(s string) bool { return s != "" && rxAlphaNum.MatchString(s) }

// IsAlphaDash string.
func IsAlphaDash(s string) bool { return s != "" && rxAlphaDash.MatchString(s) }

// IsNumber string. should >= 0
func IsNumber(v any) bool {
	v = indirectValue(v)

	if v == nil {
		return false
	}

	if s, err := strutil.ToString(v); err == nil {
		return s != "" && rxNumber.MatchString(s)
	}
	return false
}

// IsNumeric is string/int number. should >= 0
func IsNumeric(v any) bool {
	v = indirectValue(v)

	if v == nil {
		return false
	}

	if s, err := strutil.ToString(v); err == nil {
		return s != "" && rxNumber.MatchString(s)
	}
	return false
}

// IsStringNumber is string number. should >= 0
func IsStringNumber(s string) bool { return s != "" && rxNumber.MatchString(s) }

// IsEmail check
func IsEmail(s string) bool { return s != "" && rxEmail.MatchString(s) }

// IsUUID string
func IsUUID(s string) bool { return s != "" && rxUUID.MatchString(s) }

// IsUUID3 string
func IsUUID3(s string) bool { return s != "" && rxUUID3.MatchString(s) }

// IsUUID4 string
func IsUUID4(s string) bool { return s != "" && rxUUID4.MatchString(s) }

// IsUUID5 string
func IsUUID5(s string) bool { return s != "" && rxUUID5.MatchString(s) }

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func IsIP(s string) bool { return s != "" && net.ParseIP(s) != nil }

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
func IsIPv4(s string) bool {
	if s == "" {
		return false
	}

	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func IsIPv6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
func IsMAC(s string) bool {
	if s == "" {
		return false
	}
	_, err := net.ParseMAC(s)
	return err == nil
}

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func IsCIDRv4(s string) bool {
	if s == "" {
		return false
	}
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() != nil
}

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func IsCIDRv6(s string) bool {
	if s == "" {
		return false
	}

	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() == nil
}

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func IsCIDR(s string) bool {
	if s == "" {
		return false
	}

	_, _, err := net.ParseCIDR(s)
	return err == nil
}

// IsJSON check if the string is valid JSON (note: uses json.Unmarshal).
func IsJSON(s string) bool {
	if s == "" {
		return false
	}

	var js json.RawMessage
	return Unmarshal([]byte(s), &js) == nil
}

// HasLowerCase check string has lower case
func HasLowerCase(s string) bool {
	if s == "" {
		return false
	}
	return rxHasLowerCase.MatchString(s)
}

// HasUpperCase check string has upper case
func HasUpperCase(s string) bool {
	return s != "" && rxHasUpperCase.MatchString(s)
}

// StartsWith check string is starts with sub-string
func StartsWith(s, sub string) bool { return s != "" && strings.HasPrefix(s, sub) }

// EndsWith check string is ends with sub-string
func EndsWith(s, sub string) bool { return s != "" && strings.HasSuffix(s, sub) }

// StringContains check string is containing sub-string
func StringContains(s, sub string) bool { return s != "" && strings.Contains(s, sub) }

// Regexp match value string
func Regexp(str string, pattern string) bool {
	ok, _ := regexp.MatchString(pattern, str)
	return ok
}

/*************************************************************
 * global: filesystem validators
 *************************************************************/

// PathExists reports whether the named file or directory exists.
func PathExists(path string) bool { return fsutil.PathExists(path) }

// IsFilePath path is a local filepath
func IsFilePath(path string) bool { return fsutil.IsFile(path) }

// IsDirPath path is a local dir path
func IsDirPath(path string) bool { return fsutil.IsDir(path) }

// IsWinPath string
func IsWinPath(s string) bool {
	return s != "" && rxWinPath.MatchString(s)
}

// IsUnixPath string
func IsUnixPath(s string) bool {
	return s != "" && rxUnixPath.MatchString(s)
}

/*************************************************************
 * global: compare validators
 *************************************************************/

// IsEqual check two value is equals. Don't compare func, struct
//
// Support:
//
//	bool, int(X), uint(X), string, float(X) AND slice, array, map
func IsEqual(val, wantVal any) bool {
	// check is nil
	if val == nil || wantVal == nil {
		return val == wantVal
	}

	sv := removeValuePtr(reflect.ValueOf(val))
	wv := removeValuePtr(reflect.ValueOf(wantVal))

	// don't compare func, struct
	if sv.Kind() == reflect.Func || sv.Kind() == reflect.Struct {
		return false
	}
	if wv.Kind() == reflect.Func || wv.Kind() == reflect.Struct {
		return false
	}

	// compare basic type: bool, int(X), uint(X), string, float(X)
	equal, err := eq(sv, wv)

	// is not a basic type, eg: slice, array, map ...
	if err != nil {
		expBt, ok := val.([]byte)
		if !ok {
			return reflect.DeepEqual(val, wantVal)
		}

		actBt, ok := wantVal.([]byte)
		if !ok {
			return false
		}
		if expBt == nil || actBt == nil {
			return expBt == nil && actBt == nil
		}

		return bytes.Equal(expBt, actBt)
	}

	return equal
}

// NotEqual check
func NotEqual(val, wantVal any) bool { return !IsEqual(val, wantVal) }

// IntEqual check
func IntEqual(val any, wantVal int64) bool {
	// intVal, isInt := IntVal(val)
	intVal, err := mathutil.Int64(val)
	if err != nil {
		return false
	}

	return intVal == wantVal
}

// Gt check value greater dst value.
//
// only check for: int(X), uint(X), float(X), string.
func Gt(val, min any) bool { return valueCompare(val, min, ">") }

// Gte check value greater or equal dst value
// only check for: int(X), uint(X), float(X), string.
func Gte(val, min any) bool { return valueCompare(val, min, ">=") }

// Min check value greater or equal dst value, alias Gte()
// only check for: int(X), uint(X), float(X), string.
func Min(val, min any) bool { return valueCompare(val, min, ">=") }

// Lt less than dst value.
// only check for: int(X), uint(X), float(X).
func Lt(val, max any) bool { return valueCompare(val, max, "<") }

// Lte less than or equal dst value.
// only check for: int(X), uint(X), float(X).
func Lte(val, max any) bool { return valueCompare(val, max, "<=") }

// Max less than or equal dst value, alias Lte()
// only check for: int(X), uint(X), float(X).
func Max(val, max any) bool { return valueCompare(val, max, "<=") }

// Between int value in the given range.
// only check for: int(X), uint(X).
func Between(val any, min, max int64) bool {
	val = indirectValue(val)

	intVal, err := mathutil.Int64(val)
	if err != nil {
		return false
	}

	return intVal >= min && intVal <= max
}

/*************************************************************
 * global: array, slice, map validators
 *************************************************************/

// Enum value(int(X),string) should be in the given enum(strings, ints, uints).
func Enum(val, enum any) bool {
	if val == nil || enum == nil {
		return false
	}

	v, err := convToBasicType(val)
	if err != nil {
		return false
	}

	// if is string value
	if strVal, ok := v.(string); ok {
		if ss, ok := enum.([]string); ok {
			for _, strItem := range ss {
				if strVal == strItem { // exists
					return true
				}
			}
		}
		return false
	}

	// as int64 value
	intVal := v.(int64)
	if int64s, err := arrutil.ToInt64s(enum); err == nil {
		for _, i64 := range int64s {
			if intVal == i64 {
				return true
			}
		}
	}
	return false
}

// NotIn value should be not in the given enum(strings, ints, uints).
func NotIn(val, enum any) bool { return !Enum(val, enum) }

/*************************************************************
 * global: length validators
 *************************************************************/

// Length equal check for string, array, slice, map
func Length(val any, wantLen int) bool {
	ln := CalcLength(val)
	return ln != -1 && ln == wantLen
}

// MinLength check for string, array, slice, map
func MinLength(val any, minLen int) bool {
	ln := CalcLength(val)
	return ln != -1 && ln >= minLen
}

// MaxLength check for string, array, slice, map
func MaxLength(val any, maxLen int) bool {
	ln := CalcLength(val)
	return ln != -1 && ln <= maxLen
}

// ByteLength check string's length
func ByteLength(str string, minLen int, maxLen ...int) bool {
	strLen := len(str)

	// only min length check.
	if len(maxLen) == 0 {
		return strLen >= minLen
	}

	// min and max length check
	return strLen >= minLen && strLen <= maxLen[0]
}

// RuneLength check string's length (including multibyte strings)
func RuneLength(val any, minLen int, maxLen ...int) bool {
	str, isString := val.(string)
	if !isString {
		return false
	}

	// strLen := len([]rune(str))
	strLen := utf8.RuneCountInString(str)

	// only min length check.
	if len(maxLen) == 0 {
		return strLen >= minLen
	}

	// min and max length check
	return strLen >= minLen && strLen <= maxLen[0]
}

// StringLength check string's length (including multibyte strings)
func StringLength(val any, minLen int, maxLen ...int) bool {
	return RuneLength(val, minLen, maxLen...)
}

/*************************************************************
 * global: date/time validators
 *************************************************************/

// IsDate check value is an date string.
func IsDate(srcDate string, layouts ...string) bool {
	_, err := strutil.ToTime(srcDate, layouts...)
	return err == nil
}

// DateFormat check
func DateFormat(s string, layout string) bool {
	_, err := time.Parse(layout, s)
	return err == nil
}

// DateEquals check.
// Usage:
// 	DateEquals(val, "2017-05-12")
// func DateEquals(srcDate, dstDate string) bool {
// 	return false
// }

// BeforeDate check
func BeforeDate(srcDate, dstDate string) bool {
	st, err := strutil.ToTime(srcDate)
	if err != nil {
		return false
	}

	dt, err := strutil.ToTime(dstDate)
	if err != nil {
		return false
	}

	return st.Before(dt)
}

// BeforeOrEqualDate check
func BeforeOrEqualDate(srcDate, dstDate string) bool {
	st, err := strutil.ToTime(srcDate)
	if err != nil {
		return false
	}

	dt, err := strutil.ToTime(dstDate)
	if err != nil {
		return false
	}

	return st.Before(dt) || st.Equal(dt)
}

// AfterOrEqualDate check
func AfterOrEqualDate(srcDate, dstDate string) bool {
	st, err := strutil.ToTime(srcDate)
	if err != nil {
		return false
	}

	dt, err := strutil.ToTime(dstDate)
	if err != nil {
		return false
	}

	return st.After(dt) || st.Equal(dt)
}

// AfterDate check
func AfterDate(srcDate, dstDate string) bool {
	st, err := strutil.ToTime(srcDate)
	if err != nil {
		return false
	}

	dt, err := strutil.ToTime(dstDate)
	if err != nil {
		return false
	}

	return st.After(dt)
}
