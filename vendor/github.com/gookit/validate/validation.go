package validate

import (
	"fmt"
	"reflect"
	"strings"
)

// some default value settings.
const (
	fieldTag  = "json"
	filterTag = "filter"
	labelTag  = "label"

	messageTag  = "message"
	validateTag = "validate"

	filterError   = "_filter"
	validateError = "_validate"

	// sniff Length, use for detect file mime type
	sniffLen = 512
	// 32 MB
	defaultMaxMemory int64 = 32 << 20

	// validator type
	validatorTypeBuiltin int8 = 1
	validatorTypeCustom  int8 = 2
)

// Validation definition
type Validation struct {
	// for optimize create instance. refer go-playground/validator
	// v *Validation
	// pool *sync.Pool

	// source input data
	data DataFace
	// all validated fields list
	// fields []string

	// save filtered/validated safe data
	safeData M
	// filtered clean data
	filteredData M
	// save user custom set default values
	defValues map[string]any

	// Errors for validate
	Errors Errors
	// CacheKey for cache rules
	// CacheKey string
	// StopOnError If true: An error occurs, it will cease to continue to verify
	StopOnError bool
	// SkipOnEmpty Skip check on field not exist or value is empty
	SkipOnEmpty bool
	// UpdateSource Whether to update source field value, useful for struct validate
	UpdateSource bool
	// CheckDefault Whether to validate the default value set by the user
	CheckDefault bool
	// CachingRules switch. default is False
	// CachingRules bool

	// mark has error occurs
	hasError bool
	// mark is filtered
	hasFiltered bool
	// mark is validated
	hasValidated bool
	// validate rules for the validation
	rules []*Rule

	// validators for the validation. map-value: 1=builtin, 2=custom
	validators map[string]int8
	// validator func meta info
	validatorMetas map[string]*funcMeta

	// current scene name
	scene string
	// scenes config.
	// {
	// 	"create": {"field0", "field1"}
	// 	"update": {"field0", "field2"}
	// }
	scenes SValues
	// should check fields in current scene.
	sceneFields map[string]uint8

	// filtering rules for the validation
	filterRules []*FilterRule
	// filter func reflect.Value map
	filterValues map[string]reflect.Value

	// translator instance
	trans *Translator
	// optional fields, useful for sub-struct field in struct data. eg: "Parent"
	//
	// key is field name, value is field vale is: init=0 empty=1 not-empty=2.
	optionals map[string]int8
}

// NewEmpty new validation instance, but not with data.
func NewEmpty(scene ...string) *Validation {
	return NewValidation(nil, scene...)
}

// NewValidation new validation instance
func NewValidation(data DataFace, scene ...string) *Validation {
	return newValidation(data).SetScene(scene...)
}

/*************************************************************
 * validation settings
 *************************************************************/

// ResetResult reset the validate result.
func (v *Validation) ResetResult() {
	v.Errors = Errors{}
	v.hasError = false
	v.hasFiltered = false
	v.hasValidated = false
	// result data
	v.safeData = make(map[string]any)
	v.filteredData = make(map[string]any)
}

// Reset the Validation instance.
//
// Will resets:
//   - validate result
//   - validate rules
//   - validate filterRules
//   - custom validators TODO
func (v *Validation) Reset() {
	v.ResetResult()

	// v.validators = make(map[string]int8)
	v.resetRules()
}

func (v *Validation) resetRules() {
	// reset rules
	v.rules = v.rules[:0]
	v.optionals = make(map[string]int8)
	v.filterRules = v.filterRules[:0]
}

// TODO Config(opt *Options) *Validation

// WithSelf config the Validation instance. TODO rename to WithConfig
func (v *Validation) WithSelf(fn func(v *Validation)) *Validation {
	fn(v)
	return v
}

// WithTrans with a custom translator
func (v *Validation) WithTrans(trans *Translator) *Validation {
	v.trans = trans
	return v
}

// WithScenarios is alias of the WithScenes()
func (v *Validation) WithScenarios(scenes SValues) *Validation {
	return v.WithScenes(scenes)
}

// WithScenes set scene config.
//
// Usage:
//
//	v.WithScenes(SValues{
//		"create": []string{"name", "email"},
//		"update": []string{"name"},
//	})
//	ok := v.AtScene("create").Validate()
func (v *Validation) WithScenes(scenes map[string][]string) *Validation {
	v.scenes = scenes
	return v
}

// AtScene setting current validate scene.
func (v *Validation) AtScene(scene string) *Validation {
	v.scene = scene
	return v
}

// InScene alias of the AtScene()
func (v *Validation) InScene(scene string) *Validation {
	return v.AtScene(scene)
}

// SetScene alias of the AtScene()
func (v *Validation) SetScene(scene ...string) *Validation {
	if len(scene) > 0 {
		v.AtScene(scene[0])
	}
	return v
}

/*************************************************************
 * add validators for validation
 *************************************************************/

// AddValidators to the Validation instance.
func (v *Validation) AddValidators(m map[string]any) *Validation {
	for name, checkFunc := range m {
		v.AddValidator(name, checkFunc)
	}
	return v
}

// AddValidator to the Validation instance. checkFunc must return a bool.
//
// Usage:
//
//	v.AddValidator("myFunc", func(data validate.DataFace, val any) bool {
//		// do validate val ...
//		return true
//	})
func (v *Validation) AddValidator(name string, checkFunc any) *Validation {
	fv := checkValidatorFunc(name, checkFunc)

	v.validators[name] = validatorTypeCustom
	// v.validatorValues[name] = fv
	v.validatorMetas[name] = newFuncMeta(name, false, fv)

	return v
}

// ValidatorMeta get by name. get validator from global or validation instance.
func (v *Validation) validatorMeta(name string) *funcMeta {
	// current validation
	if fm, ok := v.validatorMetas[name]; ok {
		return fm
	}

	// from global validators
	if fm, ok := validatorMetas[name]; ok {
		return fm
	}

	// if v.data is StructData instance.
	if v.data.Type() == sourceStruct {
		fv, ok := v.data.(*StructData).FuncValue(name)
		if ok {
			fm := newFuncMeta(name, false, fv)
			// storage it.
			v.validators[name] = validatorTypeCustom
			v.validatorMetas[name] = fm

			return fm
		}
	}
	return nil
}

// HasValidator check
func (v *Validation) HasValidator(name string) bool {
	name = ValidatorName(name)

	// current validation
	if _, ok := v.validatorMetas[name]; ok {
		return true
	}

	// global validators
	_, ok := validatorMetas[name]
	return ok
}

// Validators get all validator names
func (v *Validation) Validators(withGlobal bool) map[string]int8 {
	if withGlobal {
		mp := make(map[string]int8)
		for name, typ := range validators {
			mp[name] = typ
		}

		for name, typ := range v.validators {
			mp[name] = typ
		}
		return mp
	}

	return v.validators
}

/*************************************************************
 * Do filtering/sanitize
 *************************************************************/

// Sanitize data by filter rules
func (v *Validation) Sanitize() bool { return v.Filtering() }

// Filtering data by filter rules
func (v *Validation) Filtering() bool {
	if v.hasFiltered {
		return v.IsSuccess()
	}

	// apply rule to validate data.
	for _, rule := range v.filterRules {
		if err := rule.Apply(v); err != nil { // has error
			v.AddError(filterError, filterError, rule.fields[0]+": "+err.Error())
			break
		}
	}

	v.hasFiltered = true
	return v.IsSuccess()
}

/*************************************************************
 * errors messages
 *************************************************************/

// WithTranslates settings. you can be custom field translates.
//
// Usage:
//
//		v.WithTranslates(map[string]string{
//			"name": "Username",
//			"pwd": "Password",
//	 })
func (v *Validation) WithTranslates(m map[string]string) *Validation {
	v.trans.AddLabelMap(m)
	return v
}

// AddTranslates settings data. like WithTranslates()
func (v *Validation) AddTranslates(m map[string]string) {
	v.trans.AddLabelMap(m)
}

// WithMessages settings. you can custom validator error messages.
//
// Usage:
//
//		// key is "validator" or "field.validator"
//		v.WithMessages(map[string]string{
//			"require": "oh! {field} is required",
//			"range": "oh! {field} must be in the range %d - %d",
//	 })
func (v *Validation) WithMessages(m map[string]string) *Validation {
	v.trans.AddMessages(m)
	return v
}

// AddMessages settings data. like WithMessages()
func (v *Validation) AddMessages(m map[string]string) {
	v.trans.AddMessages(m)
}

// WithError add error of the validation
func (v *Validation) WithError(err error) *Validation {
	if err != nil {
		v.AddError(validateError, validateError, err.Error())
	}
	return v
}

// AddError message for a field
func (v *Validation) AddError(field, validator, msg string) {
	if !v.hasError {
		v.hasError = true
	}

	field = v.trans.FieldName(field)
	v.Errors.Add(field, validator, msg)
}

// AddErrorf add a formatted error message
func (v *Validation) AddErrorf(field, msgFormat string, args ...any) {
	v.AddError(field, validateError, fmt.Sprintf(msgFormat, args...))
}

// Trans get translator
func (v *Validation) Trans() *Translator {
	// if v.trans == nil {
	// 	v.trans = StdTranslator
	// }
	return v.trans
}

func (v *Validation) convArgTypeError(field, name string, argKind, wantKind reflect.Kind, argIdx int) {
	v.AddErrorf(field, "cannot convert %s to arg#%d(%s), validator '%s'", argKind, argIdx, wantKind, name)
}

/*************************************************************
 * getter methods
 *************************************************************/

// Raw value get by key
func (v *Validation) Raw(key string) (any, bool) {
	if v.data == nil { // check input data
		return nil, false
	}
	return v.data.Get(key)
}

// RawVal value get by key
func (v *Validation) RawVal(key string) any {
	if v.data == nil { // check input data
		return nil
	}
	val, _ := v.data.Get(key)
	return val
}

// try to get value by key.
//
// **NOTE:**
//
// If v.data is StructData, will return zero value check. Other dataSource will always return `zero=False`.
func (v *Validation) tryGet(key string) (val any, exist, zero bool) {
	if v.data == nil {
		return
	}

	// find from filtered data.
	if val1, ok := v.filteredData[key]; ok {
		return val1, true, false
	}

	// find from validated data. (such as has default value)
	if val2, ok := v.safeData[key]; ok {
		return val2, true, false
	}

	// TODO add cache data v.caches[key]
	// get from source data
	return v.data.TryGet(key)
}

// Get value by key.
func (v *Validation) Get(key string) (val any, exist bool) {
	val, exist, _ = v.tryGet(key)
	return
}

// GetWithDefault get field value by key.
//
// On not found, if it has default value, will return default-value.
func (v *Validation) GetWithDefault(key string) (val any, exist, isDefault bool) {
	var zero bool
	val, exist, zero = v.tryGet(key)
	if exist && !zero {
		return
	}

	// try read custom default value
	defVal, isDefault := v.defValues[key]
	if isDefault {
		val = defVal
	}
	return
}

// Filtered get filtered value by key
func (v *Validation) Filtered(key string) any {
	return v.filteredData[key]
}

// Safe get safe value by key
func (v *Validation) Safe(key string) (val any, ok bool) {
	if v.data == nil { // check input data
		return
	}
	val, ok = v.safeData[key]
	return
}

// SafeVal get safe value by key
func (v *Validation) SafeVal(key string) any {
	val, _ := v.Safe(key)
	return val
}

// GetSafe get safe value by key
func (v *Validation) GetSafe(key string) any {
	val, _ := v.Safe(key)
	return val
}

// BindStruct binding safe data to a struct. alias of BindSafeData
func (v *Validation) BindStruct(ptr any) error { return v.BindSafeData(ptr) }

// BindSafeData binding safe data to an struct.
func (v *Validation) BindSafeData(ptr any) error {
	if len(v.safeData) == 0 { // no safe data.
		return nil
	}

	// to json bytes
	bts, err := Marshal(v.safeData)
	if err != nil {
		return err
	}
	return Unmarshal(bts, ptr)
}

// Set value by key
func (v *Validation) Set(field string, val any) error {
	// check input data
	if v.data == nil {
		return ErrEmptyData
	}

	_, err := v.data.Set(field, val)
	return err
}

// only update set value by key for struct
func (v *Validation) updateValue(field string, val any) (any, error) {
	// data source is struct
	if v.data.Type() == sourceStruct {
		return v.data.Set(strings.TrimSuffix(field, ".*"), val)
	}

	// TODO dont update value for Form and Map data source
	return val, nil
}

// SetDefValue set a default value of given field
func (v *Validation) SetDefValue(field string, val any) {
	if v.defValues == nil {
		v.defValues = make(map[string]any)
	}
	v.defValues[field] = val
}

// GetDefValue get default value of the field
func (v *Validation) GetDefValue(field string) (any, bool) {
	defVal, ok := v.defValues[field]
	return defVal, ok
}

// SceneFields field names get
func (v *Validation) SceneFields() []string {
	return v.scenes[v.scene]
}

// scene field name map build
func (v *Validation) sceneFieldMap() (m map[string]uint8) {
	if v.scene == "" {
		return
	}

	if fields, ok := v.scenes[v.scene]; ok {
		m = make(map[string]uint8, len(fields))
		for _, field := range fields {
			m[field] = 1
		}
	}
	return
}

// Scene name get for current validation
func (v *Validation) Scene() string { return v.scene }

// IsOK for the validating
func (v *Validation) IsOK() bool { return !v.hasError }

// IsFail for the validating
func (v *Validation) IsFail() bool { return v.hasError }

// IsSuccess for the validating
func (v *Validation) IsSuccess() bool { return !v.hasError }

// SafeData get all validated safe data
func (v *Validation) SafeData() M { return v.safeData }

// FilteredData return filtered data.
func (v *Validation) FilteredData() M {
	return v.filteredData
}

/*************************************************************
 * helper methods
 *************************************************************/

// on stop on error
func (v *Validation) shouldStop() bool {
	return v.hasError && v.StopOnError
}

// check current field is in optional parent field.
//
// return: true - optional parent field value is empty.
func (v *Validation) isInOptional(field string) bool {
	for name, flag := range v.optionals {
		// check like: field="Parent.Child" name="Parent"
		if strings.HasPrefix(field, name+".") {
			if flag != 0 {
				return flag == 1 // 1=empty
			}

			pVal, exist, zero := v.tryGet(name)
			if !exist || zero {
				v.optionals[name] = 1
				return true // not check field.
			}
			if IsEmpty(pVal) {
				v.optionals[name] = 1
				return true // not check field.
			}

			v.optionals[name] = 2
			return false
		}
	}

	return false
}

func (v *Validation) isNotNeedToCheck(field string) bool {
	if len(v.sceneFields) == 0 {
		return false
	}

	fields := strings.Split(field, ".")
	for i := 0; i < len(fields); i++ {
		_, ok := v.sceneFields[strings.Join(fields[0:i], ".")]
		if ok {
			return false
		}
	}

	_, ok := v.sceneFields[field]
	return !ok
}
