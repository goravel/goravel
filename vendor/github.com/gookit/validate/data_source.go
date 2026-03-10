package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/strutil"
)

const (
	// from user setting, unmarshal JSON
	sourceMap uint8 = iota + 1
	// from URL.Values, PostForm. contains Files data
	sourceForm
	// from user setting
	sourceStruct
)

// 0: top level field
// 1: field at anonymous struct
// 2: field at non-anonymous struct
const (
	fieldAtTopStruct int8 = iota
	fieldAtAnonymous
	fieldAtSubStruct
)

// data (Un)marshal func
var (
	Marshal   MarshalFunc   = json.Marshal
	Unmarshal UnmarshalFunc = json.Unmarshal
)

type (
	// MarshalFunc define
	MarshalFunc func(v any) ([]byte, error)
	// UnmarshalFunc define
	UnmarshalFunc func(data []byte, ptr any) error
)

// DataFace data source interface definition
//
// Current has three data source:
//   - map
//   - form
//   - struct
type DataFace interface {
	Type() uint8
	Src() any
	Get(key string) (val any, exist bool)
	// TryGet value by key.
	// if source data is struct, will return zero check
	TryGet(key string) (val any, exist, zero bool)
	Set(field string, val any) (any, error)
	// Create validation instance create func
	Create(err ...error) *Validation
	Validation(err ...error) *Validation
}

/*************************************************************
 * Map Data
 *************************************************************/

// MapData definition
type MapData struct {
	// Map the source map data
	Map map[string]any
	// from reflect Map
	value reflect.Value
	// bodyJSON from the original JSON bytes/string.
	// available for FromJSONBytes(), FormJSON().
	bodyJSON []byte
	// TODO map field value cache by key path
	// cache map[string]any
}

/*************************************************************
 * Map data operate
 *************************************************************/

// Src get
func (d *MapData) Src() any { return d.Map }

// Type get
func (d *MapData) Type() uint8 { return sourceMap }

// Set value by key
func (d *MapData) Set(field string, val any) (any, error) {
	d.Map[field] = val
	return val, nil
}

// Get value by key. support get value by path.
func (d *MapData) Get(field string) (any, bool) {
	// if fv, ok := d.fields[field]; ok {
	// 	return fv, true
	// }
	return maputil.GetByPath(field, d.Map)
}

// TryGet value by key
func (d *MapData) TryGet(field string) (val any, exist, zero bool) {
	val, exist = maputil.GetByPath(field, d.Map)
	return
}

// Create a Validation from data
func (d *MapData) Create(err ...error) *Validation { return d.Validation(err...) }

// Validation creates from data
func (d *MapData) Validation(err ...error) *Validation {
	if len(err) > 0 {
		return NewValidation(d).WithError(err[0])
	}
	return NewValidation(d)
}

// BindJSON binds v to the JSON data in the request body.
// It calls json.Unmarshal and sets the value of v.
func (d *MapData) BindJSON(ptr any) error {
	if len(d.bodyJSON) == 0 {
		return nil
	}
	return Unmarshal(d.bodyJSON, ptr)
}

/*************************************************************
 * Struct Data
 *************************************************************/

// ConfigValidationFace definition. you can do something on create Validation.
type ConfigValidationFace interface {
	ConfigValidation(v *Validation)
}

// FieldTranslatorFace definition. you can custom field translates.
// Usage:
//
//	type User struct {
//		Name string `json:"name" validate:"required|minLen:5"`
//	}
//
//	func (u *User) Translates() map[string]string {
//		return MS{
//			"Name": "Username",
//		}
//	}
type FieldTranslatorFace interface {
	Translates() map[string]string
}

// CustomMessagesFace definition. you can custom validator error messages.
// Usage:
//
//	type User struct {
//		Name string `json:"name" validate:"required|minLen:5"`
//	}
//
//	func (u *User) Messages() map[string]string {
//		return MS{
//			"Name.required": "oh! Username is required",
//		}
//	}
type CustomMessagesFace interface {
	Messages() map[string]string
}

// StructData definition.
//
// more struct tags define please see GlobalOption
type StructData struct {
	// source struct data, from user setting
	src any
	// max depth for parse sub-struct. TODO WIP ...
	// depth int

	// reflect.Value of the source struct
	value reflect.Value
	// reflect.Type of the source struct
	valueTyp reflect.Type
	// field names in the src struct
	// 0:common field 1:anonymous field 2:nonAnonymous field
	fieldNames map[string]int8
	// cache field reflect value info. key is path. eg: top.sub
	fieldValues map[string]reflect.Value
	// TODO field reflect values cache
	// fieldRftValues map[string]any
	// FilterTag name in the struct tags.
	//
	// see GlobalOption.FilterTag
	FilterTag string
	// ValidateTag name in the struct tags.
	//
	// see GlobalOption.ValidateTag
	ValidateTag string
}

// StructOption definition
// type StructOption struct {
// 	// ValidateTag in the struct tags.
// 	ValidateTag string
// 	// MethodName  string
// }

var (
	cmFaceType = reflect.TypeOf(new(CustomMessagesFace)).Elem()
	ftFaceType = reflect.TypeOf(new(FieldTranslatorFace)).Elem()
	cvFaceType = reflect.TypeOf(new(ConfigValidationFace)).Elem()
	timeType   = reflect.TypeOf(time.Time{})
)

// Src get
func (d *StructData) Src() any { return d.src }

// Type get
func (d *StructData) Type() uint8 { return sourceStruct }

// Validation creates a Validation from the StructData
func (d *StructData) Validation(err ...error) *Validation { return d.Create(err...) }

// Create from the StructData
//
//nolint:forcetypeassert
func (d *StructData) Create(err ...error) *Validation {
	v := NewValidation(d)
	if len(err) > 0 && err[0] != nil {
		return v.WithError(err[0])
	}

	// collect field filter/validate rules from struct tags
	d.parseRulesFromTag(v)

	// has custom config func
	if d.valueTyp.Implements(cvFaceType) {
		fv := d.value.MethodByName("ConfigValidation")
		fv.Call([]reflect.Value{reflect.ValueOf(v)})
	}

	// collect custom field translates config
	if d.valueTyp.Implements(ftFaceType) {
		fv := d.value.MethodByName("Translates")
		vs := fv.Call(nil)
		v.WithTranslates(vs[0].Interface().(map[string]string))
	}

	// collect custom error messages config
	// if reflect.PtrTo(d.valueTyp).Implements(cmFaceType) {
	if d.valueTyp.Implements(cmFaceType) {
		fv := d.value.MethodByName("Messages")
		vs := fv.Call(nil)
		v.WithMessages(vs[0].Interface().(map[string]string))
	}

	// for struct, default update source value
	v.UpdateSource = true
	return v
}

// parse and collect rules from struct tags.
func (d *StructData) parseRulesFromTag(v *Validation) {
	if d.ValidateTag == "" {
		d.ValidateTag = gOpt.ValidateTag
	}

	if d.FilterTag == "" {
		d.FilterTag = gOpt.FilterTag
	}

	fOutMap := make(map[string]string)
	var recursiveFunc func(vv reflect.Value, vt reflect.Type, preStrName string, parentIsAnonymous bool)

	vv := d.value
	vt := d.valueTyp
	// preStrName - the parent field name.
	recursiveFunc = func(vv reflect.Value, vt reflect.Type, parentFName string, parentIsAnonymous bool) {
		for i := 0; i < vt.NumField(); i++ {
			fv := vt.Field(i)
			// skip don't exported field
			name := fv.Name
			if name[0] >= 'a' && name[0] <= 'z' {
				if !gOpt.ValidatePrivateFields {
					continue
				}
			}

			if parentFName == "" {
				d.fieldNames[name] = fieldAtTopStruct
			} else {
				name = parentFName + "." + name
				if parentIsAnonymous {
					d.fieldNames[name] = fieldAtAnonymous
				} else {
					d.fieldNames[name] = fieldAtSubStruct
				}
			}

			// validate rule
			vRule := fv.Tag.Get(d.ValidateTag)
			if vRule != "" {
				v.StringRule(name, vRule)
			}

			// filter rule
			fRule := fv.Tag.Get(d.FilterTag)
			if fRule != "" {
				v.FilterRule(name, fRule)
			}

			// load field output name by FieldTag. eg: `json:"user_name"`
			outName := ""
			if gOpt.FieldTag != "" {
				outName = fv.Tag.Get(gOpt.FieldTag)
				outName = strings.SplitN(outName, ",", 2)[0]
			}

			// add pre field display name to fName
			if outName != "" {
				if parentFName != "" {
					if pOutName, ok := fOutMap[parentFName]; ok {
						outName = pOutName + "." + outName
					}
				}

				fOutMap[name] = outName
			}

			// load field translate name
			// preferred to use label tag name. eg: `label:"display name"`
			// and then use field output name. eg: `json:"user_name"`
			if gOpt.LabelTag != "" {
				v.trans.addLabelName(name, fv.Tag.Get(gOpt.LabelTag))
			}

			// load custom error messages.
			// eg: `message:"required:name is required|minLen:name min len is %d"`
			if gOpt.MessageTag != "" {
				errMsg := fv.Tag.Get(gOpt.MessageTag)
				if errMsg != "" {
					d.loadMessagesFromTag(v.trans, name, vRule, errMsg)
				}
			}

			ft := removeTypePtr(vt.Field(i).Type)

			// collect rules from sub-struct and from arrays/slices elements
			if ft != timeType && removeValuePtr(vv).IsValid() {

				// feat: only collect sub-struct rule on current field has rule.
				if vRule == "" && gOpt.CheckSubOnParentMarked {
					continue
				}

				fValue := removeValuePtr(vv).Field(i)

				switch ft.Kind() {
				case reflect.Struct:
					recursiveFunc(fValue, ft, name, fv.Anonymous)

				case reflect.Array, reflect.Slice:
					fValue = removeValuePtr(fValue)

					// Check if the reflect.Value is valid and not a nil pointer
					if !fValue.IsValid() || (ft.Kind() == reflect.Slice && fValue.IsNil()) {
						continue
					}
					// perf: skip parse on elements is simple kind
					if reflects.IsSimpleKind(ft.Elem().Kind()) {
						continue
					}

					for j := 0; j < fValue.Len(); j++ {
						elemValue := removeValuePtr(fValue.Index(j))
						elemType := removeTypePtr(elemValue.Type())

						arrayName := fmt.Sprintf("%s.%d", name, j)
						if outName != "" {
							fOutMap[arrayName] = fmt.Sprintf("%s.%d", outName, j)
						}
						if elemType.Kind() == reflect.Struct {
							recursiveFunc(elemValue, elemType, arrayName, fv.Anonymous)
						}
					}

				case reflect.Map:
					fValue = removeValuePtr(fValue)

					// Check if the reflect.Value is valid and not a nil pointer
					if !fValue.IsValid() || fValue.IsNil() {
						continue
					}

					for _, key := range fValue.MapKeys() {
						key = removeValuePtr(key)
						elemValue := removeValuePtr(fValue.MapIndex(key))
						elemType := removeTypePtr(elemValue.Type())

						format := "%s."
						kind := key.Kind()
						val := key.Interface()
						switch {
						case kind == reflect.String:
							format += "%s"
							val = strings.ReplaceAll(key.String(), "\"", "")
						case kind >= reflect.Int && kind <= reflect.Uint64:
							format += "%d"
						case kind >= reflect.Float32 && kind <= reflect.Complex128:
							format += "%f"
						default:
							format += "%#v"
						}

						arrayName := fmt.Sprintf(format, name, val)
						if outName != "" {
							fOutMap[arrayName] = fmt.Sprintf(format, outName, val)
						}
						if elemType.Kind() == reflect.Struct {
							recursiveFunc(elemValue, elemType, arrayName, fv.Anonymous)
						}
					}
				default:
					// do nothing
				}
			}
		}
	}

	recursiveFunc(removeValuePtr(vv), vt, "", false)

	if len(fOutMap) > 0 {
		v.Trans().AddFieldMap(fOutMap)
	}
}

// eg: `message:"required:name is required|minLen:name min len is %d"`
func (d *StructData) loadMessagesFromTag(trans *Translator, field, vRule, vMsg string) {
	var msgKey, vName string
	var vNames []string

	// only one message, use for first validator.
	// eg: `message:"name is required"`
	if !strings.ContainsRune(vMsg, '|') {
		// eg: `message:"required:name is required"`
		if strings.ContainsRune(vMsg, ':') {
			nodes := strings.SplitN(vMsg, ":", 2)
			vName = strings.TrimSpace(nodes[0])
			vNames = []string{vName}
			// first is validator name
			vMsg = strings.TrimSpace(nodes[1])
		}

		if vName == "" {
			// eg `validate:"required|date"`
			vNames = []string{vRule}
			if strings.ContainsRune(vRule, '|') {
				vNames = strings.Split(vRule, "|")
			}

			for i, node := range vNames {
				// has params for validator: "minLen:5"
				if strings.ContainsRune(node, ':') {
					tmp := strings.SplitN(node, ":", 2)
					vNames[i] = tmp[0]
				}
			}
		}

		// if rName, has := validatorAliases[validator]; has {
		// 	msgKey = field + "." + rName
		// } else {

		for _, name := range vNames {
			msgKey = field + "." + name
			trans.AddMessage(msgKey, vMsg)
		}

		return
	}

	// multi message for validators
	// eg: `message:"required:name is required | minLen:name min len is %d"`
	for _, validatorWithMsg := range strings.Split(vMsg, "|") {
		// validatorWithMsg eg: "required:name is required"
		nodes := strings.SplitN(validatorWithMsg, ":", 2)

		validator := nodes[0]
		msgKey = field + "." + validator

		trans.AddMessage(msgKey, strings.TrimSpace(nodes[1]))
	}
}

/*
 **************************************************************
 * Struct data operate
 **************************************************************
 */

// Get value by field name. support get sub-value by path.
func (d *StructData) Get(field string) (val any, exist bool) {
	val, exist, _ = d.TryGet(field)
	return
}

// TryGet value by field name. support get sub-value by path.
func (d *StructData) TryGet(field string) (val any, exist, zero bool) {
	field = strutil.UpperFirst(field)
	// try read from cache
	if fv, ok := d.fieldValues[field]; ok {
		return fv.Interface(), true, fv.IsZero()
	}

	// var isPtr bool
	var fv reflect.Value
	// want to get sub struct field.
	if strings.IndexByte(field, '.') > 0 {
		fieldNodes := strings.Split(field, ".")
		topLevelField, ok := d.valueTyp.FieldByName(fieldNodes[0])
		if !ok {
			return
		}

		kind := removeTypePtr(topLevelField.Type).Kind()
		if kind != reflect.Struct && kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map {
			return
		}

		fv = removeValuePtr(d.value.FieldByName(fieldNodes[0]))
		if !fv.IsValid() {
			return
		}

		fieldNodes = fieldNodes[1:]
		// lastIndex := len(fieldNodes) - 1

		// last key is wildcard, return all sub-value
		if len(fieldNodes) == 1 && fieldNodes[0] == maputil.Wildcard {
			return fv.Interface(), true, fv.IsZero()
		}

		kind = fv.Type().Kind()
		for _, fieldNode := range fieldNodes {
			// fieldNode = strings.ReplaceAll(fieldNode, "\"", "") // for strings as keys
			switch kind {
			case reflect.Array, reflect.Slice:
				index, _ := strconv.Atoi(fieldNode)
				fv = fv.Index(index)
			case reflect.Map:
				fv = fv.MapIndex(reflect.ValueOf(fieldNode))
			case reflect.Struct:
				fv = fv.FieldByName(fieldNode)
			default: // no sub-value, should never have happened
				return
			}

			// isPtr = fv.Kind() == reflect.Pointer
			fv = removeValuePtr(fv)
			if !fv.IsValid() {
				return
			}

			kind = fv.Type().Kind()
			// if IsZero(fv) || (fv.Kind() == reflect.Ptr && fv.IsNil()) {
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				return
			}
		}

		d.fieldNames[field] = fieldAtSubStruct
	} else {
		// field at top struct
		fv = d.value.FieldByName(field)
		if !fv.IsValid() { // field not exists
			return
		}

		// is it a pointer
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() { // fix: top-field is nil
				return
			}
			// fv = removeValuePtr(fv)
		}
	}

	// check can interface
	if fv.CanInterface() {
		// TIP: if is zero value, as not exist.
		// - bool as exists.
		// if fv.Kind() != reflect.Bool && IsZero(fv) {
		// 	return fv.Interface(), false
		// 	// return nil, false
		// }

		// cache field value info
		d.fieldValues[field] = fv
		// isZero := fv.IsZero()
		// if isPtr {
		// 	isZero = fv.Elem().IsZero()
		// }

		return fv.Interface(), true, fv.IsZero()
	}
	return
}

// Set value by field name.
//
// Notice: `StructData.src` the incoming struct must be a pointer to set the value
func (d *StructData) Set(field string, val any) (newVal any, err error) {
	field = strutil.UpperFirst(field)
	if !d.HasField(field) { // field not found
		return nil, ErrNoField
	}

	// find from cache
	fv, ok := d.fieldValues[field]
	if !ok {
		f := d.fieldNames[field]
		switch f {
		case fieldAtTopStruct:
			fv = d.value.FieldByName(field)
		case fieldAtAnonymous:
		case fieldAtSubStruct:
			fieldNodes := strings.Split(field, ".")
			if len(fieldNodes) < 2 {
				return nil, ErrInvalidData
			}

			fv = d.value.FieldByName(fieldNodes[0])
			fieldNodes = fieldNodes[1:]

			for _, fieldNode := range fieldNodes {
				switch fv.Type().Kind() {
				case reflect.Array, reflect.Slice:
					index, err := strconv.Atoi(fieldNode)
					if err != nil {
						return nil, ErrInvalidData
					}

					fv = fv.Index(index)
				case reflect.Map:
					fv = fv.MapIndex(reflect.ValueOf(fieldNode))
				default:
					fv = fv.FieldByName(fieldNode)
				}
			}
		default:
			return nil, ErrNoField
		}

		// add cache
		d.fieldValues[field] = fv
	}

	// check whether the value of v can be changed.
	fv = removeValuePtr(fv)
	if !fv.CanSet() {
		return nil, ErrUnaddressableField
	}

	// Notice: need convert value type
	// - check whether you can direct convert type
	rftVal := removeValuePtr(reflect.ValueOf(val))
	if rftVal.Type().ConvertibleTo(fv.Type()) {
		fv.Set(rftVal.Convert(fv.Type()))
		return val, nil
	}

	// try manual convert type
	newRv, err1 := reflects.ValueByKind(val, fv.Kind())
	if err1 != nil {
		return nil, err1
	}

	// update field value
	fv.Set(newRv)

	// fix: custom func return ptr type, fv.Kind() always is real type.
	if rftVal.Kind() != fv.Kind() {
		newVal = newRv.Interface()
	}
	return
}

// FuncValue get func value in the src struct
func (d *StructData) FuncValue(name string) (reflect.Value, bool) {
	fv := d.value.MethodByName(strutil.UpperFirst(name))
	return fv, fv.IsValid()
}

// HasField in the src struct
func (d *StructData) HasField(field string) bool {
	if _, ok := d.fieldNames[field]; ok {
		return true
	}

	// has field, cache it
	if _, ok := d.valueTyp.FieldByName(field); ok {
		d.fieldNames[field] = fieldAtTopStruct
		return true
	}
	return false
}

/*************************************************************
 * Form Data
 *************************************************************/

// FormData obtained from the request body or url query parameters or user custom setting.
type FormData struct {
	// Form holds any basic key-value string data
	// This includes all fields from urlencoded form,
	// and the form fields only (not files) from a multipart form
	Form url.Values
	// Files holds files from a multipart form only.
	// For any other type of request, it will always
	// be empty. Files only supports one file per key,
	// since this is by far the most common use. If you
	// need to have more than one file per key, parse the
	// files manually using r.MultipartForm.File.
	Files map[string][]*multipart.FileHeader
}

func newFormData() *FormData {
	return &FormData{
		Form:  make(map[string][]string),
		Files: make(map[string][]*multipart.FileHeader),
	}
}

//
/*************************************************************
 * Form data operate
 *************************************************************/
//

// Src data get
func (d *FormData) Src() any { return d.Form }

// Type get
func (d *FormData) Type() uint8 { return sourceForm }

// Create a Validation from data
func (d *FormData) Create(err ...error) *Validation {
	return d.Validation(err...)
}

// Validation creates from data
func (d *FormData) Validation(err ...error) *Validation {
	if len(err) > 0 && err[0] != nil {
		return NewValidation(d).WithError(err[0])
	}
	return NewValidation(d)
}

// Add adds the value to key. It appends to any existing values associated with a key.
func (d *FormData) Add(key, value string) { d.Form.Add(key, value) }

// AddValues to Data.Form
func (d *FormData) AddValues(values url.Values) {
	for key, vals := range values {
		for _, val := range vals {
			d.Form.Add(key, val)
		}
	}
}

// AddFiles adds the multipart form files to data
func (d *FormData) AddFiles(filesMap map[string][]*multipart.FileHeader) {
	for key, files := range filesMap {
		for _, file := range files {
			d.AddFile(key, file)
		}
	}
}

// AddFile adds the multipart form file to data with the given key.
func (d *FormData) AddFile(key string, files ...*multipart.FileHeader) {
	if _, ok := d.Files[key]; !ok {
		d.Files[key] = []*multipart.FileHeader{}
	}
	d.Files[key] = append(d.Files[key], files...)
}

// Del deletes the values associated with a key.
func (d *FormData) Del(key string) { d.Form.Del(key) }

// DelFile deletes the file associated with a key (if any).
// If there is no file associated with a key, it does nothing.
func (d *FormData) DelFile(key string) { delete(d.Files, key) }

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux") sorted by key.
// Any files in d will be ignored because there is no direct way to convert a file to a
// URL encoded value.
func (d *FormData) Encode() string { return d.Form.Encode() }

// Set sets the key to value. It replaces any existing values.
func (d *FormData) Set(field string, val any) (newVal any, err error) {
	newVal = val
	switch tpVal := val.(type) {
	case string:
		d.Form.Set(field, tpVal)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		newVal = strutil.MustString(val)
		d.Form.Set(field, newVal.(string))
	default:
		err = fmt.Errorf("set value failure for field: %s", field)
	}
	return
}

// TryGet value by key
func (d *FormData) TryGet(key string) (val any, exist, zero bool) {
	val, exist = d.Get(key)
	return
}

// Get value by key
func (d *FormData) Get(key string) (any, bool) {
	// get form value
	key, _, expectArray := strings.Cut(key, ".*")
	if vs, ok := d.Form[key]; ok && len(vs) > 0 {
		if len(vs) > 1 || expectArray {
			return vs, true
		}
		return vs[0], true
	}

	// get uploaded file
	if fh, ok := d.Files[key]; ok && len(fh) > 0 {
		if len(fh) > 1 || expectArray {
			return fh, true
		}
		return fh[0], true
	}
	return nil, false
}

// String value gets by key
func (d *FormData) String(key string) string { return d.Form.Get(key) }

// Strings value gets by key
func (d *FormData) Strings(key string) []string { return d.Form[key] }

// GetFile returns the multipart form file associated with a key, if any, as a *multipart.FileHeader.
// If there is no file associated with a key, it returns nil. If you just want the body of the
// file, use GetFileBytes.
func (d *FormData) GetFile(key string) *multipart.FileHeader {
	if fh, ok := d.Files[key]; ok && len(fh) > 0 {
		return fh[0]
	}
	return nil
}

// GetFiles returns the multipart form files associated with a key, if any, as a []*multipart.FileHeader.
func (d *FormData) GetFiles(key string) []*multipart.FileHeader { return d.Files[key] }

// Has key in the Data
func (d *FormData) Has(key string) bool {
	if vs, ok := d.Form[key]; ok && len(vs) > 0 {
		return true
	}

	if _, ok := d.Files[key]; ok {
		return true
	}
	return false
}

// HasField returns true iff data.Form[key] exists. When parsing a request body, the key
// is considered to be in existence if it was provided in the request body, even if its value
// is empty.
func (d *FormData) HasField(key string) bool {
	_, found := d.Form[key]
	return found
}

// HasFile returns true iff data.Files[key] exists. When parsing a request body, the key
// is considered to be in existence if it was provided in the request body, even if the file
// is empty.
func (d *FormData) HasFile(key string) bool {
	key, _, _ = strings.Cut(key, ".*")
	_, found := d.Files[key]
	return found
}

// Int returns the first element in data[key] converted to an int.
func (d *FormData) Int(key string) int {
	if val := d.String(key); val != "" {
		iVal, _ := strconv.Atoi(val)
		return iVal
	}
	return 0
}

// Int64 returns the first element in data[key] converted to an int64.
func (d *FormData) Int64(key string) int64 {
	if val := d.String(key); val != "" {
		i64, _ := strconv.ParseInt(val, 10, 0)
		return i64
	}
	return 0
}

// Float returns the first element in data[key] converted to a float.
func (d *FormData) Float(key string) float64 {
	if val := d.String(key); val != "" {
		result, _ := strconv.ParseFloat(val, 64)
		return result
	}
	return 0
}

// Bool returns the first element in data[key] converted to a bool.
func (d *FormData) Bool(key string) bool {
	if val := d.String(key); val != "" {
		return strutil.SafeBool(val)
	}
	return false
}

// FileBytes returns the body of the file associated with key. If there is no
// file associated with key, it returns nil (not an error). It may return an error if
// there was a problem reading the file. If you need to know whether or not the file
// exists (i.e. whether it was provided in the request), use the FileExists method.
func (d *FormData) FileBytes(field string) ([]byte, error) {
	fh := d.GetFile(field)
	if fh == nil {
		return nil, nil
	}

	file, err := fh.Open()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

// FileMimeType get File Mime Type name. eg "image/png"
func (d *FormData) FileMimeType(field string) (mime string) {
	return d.fileMimeType(d.GetFile(field))
}

// FilesMimeType get all File Mime Type names by field.
func (d *FormData) FilesMimeType(field string) (mimes []string) {
	for _, file := range d.GetFiles(field) {
		mimes = append(mimes, d.fileMimeType(file))
	}
	return
}

func (d *FormData) fileMimeType(fh *multipart.FileHeader) (mime string) {
	if fh == nil {
		return
	}

	if file, err := fh.Open(); err == nil {
		var buf [sniffLen]byte
		n, _ := io.ReadFull(file, buf[:])
		mime = http.DetectContentType(buf[:n])
		// strip optional parameters (e.g., "charset") from MIME type
		// "text/plain; charset=utf-8" → "text/plain"
		mime = strings.TrimSpace(strings.Split(mime, ";")[0])
	}
	return
}
