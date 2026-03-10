// Package validate is a generic go data validate, filtering library.
//
// Source code and other details for the project are available at GitHub:
//
//	https://github.com/gookit/validate
package validate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/gookit/goutil/reflects"
)

// M is a short name for map[string]any
type M map[string]any

// MS is a short name for map[string]string
type MS map[string]string

// SValues simple values
type SValues map[string][]string

// One get one item's value string
func (ms MS) One() string {
	for _, msg := range ms {
		return msg
	}
	return ""
}

// String convert map[string]string to string
func (ms MS) String() string {
	if len(ms) == 0 {
		return ""
	}

	ss := make([]string, 0, len(ms))
	for name, msg := range ms {
		ss = append(ss, " "+name+": "+msg)
	}

	return strings.Join(ss, "\n")
}

// GlobalOption settings for validate
type GlobalOption struct {
	// FilterTag name in the struct tags.
	//
	// default: filter
	FilterTag string
	// ValidateTag in the struct tags.
	//
	// default: validate
	ValidateTag string
	// FieldTag the output field name in the struct tags.
	// it as placeholder on error message.
	//
	// default: json
	FieldTag string
	// LabelTag the display name in the struct tags.
	// use for define field translate name on error.
	//
	// default: label
	LabelTag string
	// MessageTag define error message for the field.
	//
	// default: message
	MessageTag string
	// DefaultTag define default value for the field.
	//
	// tag: default TODO
	DefaultTag string
	// StopOnError If true: An error occurs, it will cease to continue to verify. default is True.
	StopOnError bool
	// SkipOnEmpty Skip check on field not exist or value is empty. default is True.
	SkipOnEmpty bool
	// UpdateSource Whether to update source field value, useful for struct validate
	UpdateSource bool
	// CheckDefault Whether to validate the default value set by the user
	CheckDefault bool
	// CheckZero Whether validate the default zero value. (intX,uintX: 0, string: "")
	CheckZero bool
	// ErrKeyFmt config. TODO
	//
	// allow:
	// - 0 use struct field name as key. (for compatible)
	// - 1 use FieldTag defined name as key.
	ErrKeyFmt int8
	// CheckSubOnParentMarked True: only collect sub-struct rule on current field has rule.
	CheckSubOnParentMarked bool
	// ValidatePrivateFields Whether to validate private fields or not, especially when inheriting other other structs.
	//
	//  type foo struct {
	//	  Field int `json:"field" validate:"required"`
	//  }
	//  type bar struct {
	//    foo // <-- validate this field
	//    Field2 int `json:"field2" validate:"required"`
	//  }
	//
	// default: false
	ValidatePrivateFields bool

	// RestoreRequestBody Whether to restore the request body after reading it.
	// default: false
	RestoreRequestBody bool
}

// global options
var gOpt = newGlobalOption()

// Config global options
func Config(fn func(opt *GlobalOption)) {
	fn(gOpt)
}

// ResetOption reset global option
func ResetOption() {
	*gOpt = *newGlobalOption()
}

// Option get global options
func Option() GlobalOption {
	return *gOpt
}

func newGlobalOption() *GlobalOption {
	return &GlobalOption{
		StopOnError: true,
		SkipOnEmpty: true,
		// tag name in struct tags
		FieldTag: fieldTag,
		// label tag - display name in struct tags
		LabelTag: labelTag,
		// tag name in struct tags
		FilterTag:  filterTag,
		MessageTag: messageTag,
		// tag name in struct tags
		ValidateTag: validateTag,
	}
}

// pool for validation instance
// var vPool = &sync.Pool{
// 	New: func() any {
// 		return newEmpty()
// 	},
// }

func newValidation(data DataFace) *Validation {
	// v := vPool.Get().(*Validation)
	// // reset some runtime data
	// v.ResetResult()
	// v.trans = NewTranslator()

	v := newEmpty()
	v.data = data
	return v
}

func newEmpty() *Validation {
	v := &Validation{
		Errors: make(Errors),
		// create message translator
		// trans: StdTranslator,
		trans: NewTranslator(),
		// validated data
		safeData:  make(map[string]any),
		optionals: make(map[string]int8),
		// validator names
		validators: make(map[string]int8, 16),
		// filtered data
		filteredData: make(map[string]any),
		// default config
		StopOnError: gOpt.StopOnError,
		SkipOnEmpty: gOpt.SkipOnEmpty,
	}

	// init build in context validator
	ctxValidatorMap := map[string]reflect.Value{
		"required":           reflect.ValueOf(v.Required),
		"requiredIf":         reflect.ValueOf(v.RequiredIf),
		"requiredUnless":     reflect.ValueOf(v.RequiredUnless),
		"requiredWith":       reflect.ValueOf(v.RequiredWith),
		"requiredWithAll":    reflect.ValueOf(v.RequiredWithAll),
		"requiredWithout":    reflect.ValueOf(v.RequiredWithout),
		"requiredWithoutAll": reflect.ValueOf(v.RequiredWithoutAll),
		// field compare
		"eqField":  reflect.ValueOf(v.EqField),
		"neField":  reflect.ValueOf(v.NeField),
		"gtField":  reflect.ValueOf(v.GtField),
		"gteField": reflect.ValueOf(v.GteField),
		"ltField":  reflect.ValueOf(v.LtField),
		"lteField": reflect.ValueOf(v.LteField),
		// file upload check
		"isFile":      reflect.ValueOf(v.IsFormFile),
		"isImage":     reflect.ValueOf(v.IsFormImage),
		"inMimeTypes": reflect.ValueOf(v.InMimeTypes),
	}

	v.validatorMetas = make(map[string]*funcMeta, len(ctxValidatorMap))

	// make and collect meta info
	for n, fv := range ctxValidatorMap {
		v.validators[n] = validatorTypeBuiltin
		v.validatorMetas[n] = newFuncMeta(n, true, fv)
	}

	return v
}

/*************************************************************
 * quick create Validation
 *************************************************************/

// New create a Validation instance
//
// data type support:
//   - DataFace
//   - M/map[string]any
//   - SValues/url.Values/map[string][]string
//   - struct ptr
func New(data any, scene ...string) *Validation {
	switch td := data.(type) {
	case DataFace:
		return NewValidation(td, scene...)
	case M:
		return FromMap(td).Create().SetScene(scene...)
	case map[string]any:
		return FromMap(td).Create().SetScene(scene...)
	case SValues:
		return FromURLValues(url.Values(td)).Create().SetScene(scene...)
	case url.Values:
		return FromURLValues(td).Create().SetScene(scene...)
	case map[string][]string:
		return FromURLValues(td).Create().SetScene(scene...)
	}

	return Struct(data, scene...)
}

// NewWithOptions new Validation with options TODO
// func NewWithOptions(data any, fn func(opt *GlobalOption)) *Validation {
// 	fn(gOpt)
// 	return New(data)
// }

// Map validation create
func Map(m map[string]any, scene ...string) *Validation {
	return FromMap(m).Create().SetScene(scene...)
}

// MapWithRules validation create and with rules
// func MapWithRules(m map[string]any, rules MS) *Validation {
// 	return FromMap(m).Create().StringRules(rules)
// }

// JSON create validation from JSON string.
func JSON(s string, scene ...string) *Validation {
	return mustNewValidation(FromJSON(s)).SetScene(scene...)
}

// Struct validation create
func Struct(s any, scene ...string) *Validation {
	return mustNewValidation(FromStruct(s)).SetScene(scene...)
}

// Request validation create
func Request(r *http.Request) *Validation {
	return mustNewValidation(FromRequest(r))
}

func mustNewValidation(d DataFace, err error) *Validation {
	if d == nil {
		if err != nil {
			return NewValidation(d).WithError(err)
		}
		return NewValidation(d)
	}

	return d.Create(err)
}

/*************************************************************
 * create data-source instance
 *************************************************************/

// FromMap build data instance.
func FromMap(m map[string]any) *MapData {
	data := &MapData{}
	if m != nil {
		data.Map = m
		data.value = reflect.ValueOf(m)
	}
	return data
}

// FromJSON string build data instance.
func FromJSON(s string) (*MapData, error) {
	return FromJSONBytes([]byte(s))
}

// FromJSONBytes string build data instance.
func FromJSONBytes(bs []byte) (*MapData, error) {
	mp := map[string]any{}
	if err := json.Unmarshal(bs, &mp); err != nil {
		return nil, err
	}

	data := &MapData{
		Map:   mp,
		value: reflect.ValueOf(mp),
		// save JSON bytes
		bodyJSON: bs,
	}

	return data, nil
}

// FromStruct create a Data from struct
func FromStruct(s any) (*StructData, error) {
	data := &StructData{
		ValidateTag: gOpt.ValidateTag,
		// init map
		fieldNames:  make(map[string]int8),
		fieldValues: make(map[string]reflect.Value),
	}

	if s == nil {
		return data, ErrInvalidData
	}

	val := reflects.Elem(reflect.ValueOf(s))
	typ := val.Type()

	if val.Kind() != reflect.Struct || typ == timeType {
		return data, ErrInvalidData
	}

	data.src = s
	data.value = val
	data.valueTyp = typ

	return data, nil
}

var jsonContent = regexp.MustCompile(`(?i)application/((\w|\.|-)+\+)?json(-seq)?`)

// FromRequest collect data from request instance
func FromRequest(r *http.Request, maxMemoryLimit ...int64) (DataFace, error) {
	// nobody. like GET DELETE ....
	if r.Method != "POST" && r.Method != "PUT" && r.Method != "PATCH" {
		return FromURLValues(r.URL.Query()), nil
	}

	cType := r.Header.Get("Content-Type")

	// contains file uploaded form
	// strings.HasPrefix(mediaType, "multipart/")
	if strings.Contains(cType, "multipart/form-data") {
		maxMemory := defaultMaxMemory
		if len(maxMemoryLimit) > 0 {
			maxMemory = maxMemoryLimit[0]
		}

		if err := r.ParseMultipartForm(maxMemory); err != nil {
			return nil, err
		}

		// collect from values
		data := FromURLValues(r.MultipartForm.Value)
		// collect uploaded files
		data.AddFiles(r.MultipartForm.File)
		// add queries data
		data.AddValues(r.URL.Query())
		return data, nil
	}

	// basic POST form. content type: application/x-www-form-urlencoded
	if strings.Contains(cType, "form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}

		data := FromURLValues(r.PostForm)
		// add queries data
		data.AddValues(r.URL.Query())
		return data, nil
	}

	// JSON body request
	if jsonContent.MatchString(cType) {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		// restore request body
		if gOpt.RestoreRequestBody {
			r.Body = io.NopCloser(bytes.NewBuffer(bs))
		}
		return FromJSONBytes(bs)
	}

	return nil, ErrEmptyData
}

// FromURLValues build data instance.
func FromURLValues(values url.Values) *FormData {
	data := newFormData()
	for key, vals := range values {
		for _, val := range vals {
			data.Add(key, val)
		}
	}

	return data
}

// FromQuery build data instance.
//
// Usage:
//
//	validate.FromQuery(r.URL.Query()).Create()
func FromQuery(values url.Values) *FormData {
	return FromURLValues(values)
}
