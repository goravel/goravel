package filter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

// Filtration definition. Sanitization Sanitizing Sanitize
type Filtration struct {
	err error
	// raw data
	data map[string]any
	// mark has apply filters
	filtered bool
	// filtered and clean data
	cleanData map[string]any
	// filter rules
	filterRules []*Rule
}

// New a Filtration
func New(data map[string]any) *Filtration {
	return &Filtration{
		data: data,
		// init map
		cleanData: make(map[string]any),
	}
}

// LoadData set raw data for filtering.
func (f *Filtration) LoadData(data map[string]any) {
	f.data = data
}

// ResetData reset raw and filtered data
func (f *Filtration) ResetData(resetRaw bool) {
	f.err = nil
	f.filtered = false

	// reset data.
	if resetRaw {
		f.data = make(map[string]any)
	}

	f.cleanData = make(map[string]any)
}

// ResetRules reset rules and filtered data
func (f *Filtration) ResetRules() {
	f.err = nil
	f.filtered = false

	// clear rules
	f.filterRules = f.filterRules[:0]

	// clear cleanData
	f.cleanData = make(map[string]any)
}

// Clear all data and rules
func (f *Filtration) Clear() {
	f.data = make(map[string]any)
	f.ResetRules()
}

/*************************************************************
 * add rules and filtering data
 *************************************************************/

// AddRule add filter(s) rule.
//
// Usage:
//
//	f.AddRule("name", "trim")
//	f.AddRule("age", "int")
//	f.AddRule("age", "trim|int")
func (f *Filtration) AddRule(field string, rule any) *Rule {
	fields := strutil.Split(field, ",")
	if len(fields) == 0 {
		panic("filter: invalid fields parameters, cannot be empty")
	}

	r := newRule(fields)

	if strRule, ok := rule.(string); ok {
		strRule = strings.TrimSpace(strRule)
		rules := strutil.Split(strings.Trim(strRule, "|:"), "|")

		if len(rules) == 0 {
			panic("filter: invalid 'rule' params, cannot be empty")
		}

		r.AddFilters(rules...)
	} else if fn, ok := rule.(func(any) (any, error)); ok {
		r.SetFilterFunc(fn)
	} else {
		panic("filter: 'rule' params cannot be empty and type allow: string, func")
	}

	f.filterRules = append(f.filterRules, r)
	return r
}

// AddRules add multi rules.
//
// Usage:
//
//	f.AddRules(map[string]string{
//		"name": "trim|lower",
//		"age": "trim|int",
//	})
func (f *Filtration) AddRules(rules map[string]string) *Filtration {
	for field, rule := range rules {
		f.AddRule(field, rule)
	}
	return f
}

// Sanitize is alias of the Filtering()
func (f *Filtration) Sanitize() error {
	return f.Filtering()
}

// Filtering apply all filter rules, filtering data
func (f *Filtration) Filtering() error {
	if f.filtered || f.err != nil {
		return f.err
	}

	// apply rule to validate data.
	for _, rule := range f.filterRules {
		if err := rule.Apply(f); err != nil { // has error
			f.err = err
			break
		}
	}

	f.filtered = true
	return f.err
}

// IsOK of to apply filters
func (f *Filtration) IsOK() bool {
	return f.err == nil
}

// Err get error
func (f *Filtration) Err() error {
	return f.err
}

/*************************************************************
 * get raw/filtered data value
 *************************************************************/

// Raw get raw value by key
func (f *Filtration) Raw(key string) (any, bool) {
	return maputil.GetByPath(key, f.data)
}

// Safe get filtered value by key
func (f *Filtration) Safe(key string) (any, bool) {
	return maputil.GetByPath(key, f.cleanData)
}

// SafeVal get filtered value by key
func (f *Filtration) SafeVal(key string) any {
	val, _ := maputil.GetByPath(key, f.cleanData)
	return val
}

// Get value by key
func (f *Filtration) Get(key string) (any, bool) {
	val, ok := maputil.GetByPath(key, f.cleanData)
	if !ok {
		val, ok = maputil.GetByPath(key, f.data)
	}

	return val, ok
}

// MustGet value by key
func (f *Filtration) MustGet(key string) any {
	val, _ := f.Get(key)
	return val
}

// Int get a int value from filtered data.
func (f *Filtration) Int(key string) int {
	if val, ok := f.Safe(key); ok {
		return MustInt(val)
	}
	return 0
}

// Int64 get a int value from filtered data.
func (f *Filtration) Int64(key string) int64 {
	if val, ok := f.Safe(key); ok {
		return MustInt64(val)
	}
	return 0
}

// Bool value get from the filtered data.
func (f *Filtration) Bool(key string) bool {
	if val, ok := f.Safe(key); ok {
		return val.(bool)
	}

	return false
}

// String get a string value from filtered data.
func (f *Filtration) String(key string) string {
	val, ok := f.Safe(key)
	if !ok {
		return ""
	}

	// is string.
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprint(val)
}

// BindStruct bind the filtered data to struct.
func (f *Filtration) BindStruct(ptr any) error {
	bts, err := json.Marshal(f.cleanData)
	if err != nil {
		return err
	}

	return json.Unmarshal(bts, ptr)
}

// RawData get raw data
func (f *Filtration) RawData() map[string]any {
	return f.data
}

// CleanData get filtered data
func (f *Filtration) CleanData() map[string]any {
	return f.cleanData
}

/*************************************************************
 * filtering rule
 *************************************************************/

// Rule definition
type Rule struct {
	// fields to filter
	fields []string
	// filter name list
	filters []string
	// filter args. { index: "args" }
	filterArgs map[int]string
	// user custom filter func
	filterFunc func(val any) (any, error)
	// default value for the rule
	defaultVal any
}

func newRule(fields []string) *Rule {
	return &Rule{
		fields: fields,
		// init map
		filterArgs: make(map[int]string),
	}
}

// SetDefaultVal set default value for the rule
func (r *Rule) SetDefaultVal(defaultVal any) *Rule {
	r.defaultVal = defaultVal
	return r
}

// SetFilterFunc user custom filter func
func (r *Rule) SetFilterFunc(fn func(val any) (any, error)) *Rule {
	r.filterFunc = fn
	return r
}

// AddFilters add multi filter(s).
//
// Usage:
//
//	r.AddFilters("int", "str2arr:,")
func (r *Rule) AddFilters(filters ...string) *Rule {
	for _, filterName := range filters {
		pos := strings.IndexRune(filterName, ':')
		if pos > 0 { // has filter args
			name := filterName[:pos]
			index := len(r.filters)
			r.filters = append(r.filters, name)
			r.filterArgs[index] = filterName[pos+1:]
		} else {
			r.filters = append(r.filters, filterName)
		}
	}

	return r
}

// Apply rule for the rule fields
func (r *Rule) Apply(f *Filtration) (err error) {
	// validate field
	for _, field := range r.Fields() {
		// get field value.
		val, has := f.Get(field)
		if !has { // no field
			if r.defaultVal == nil {
				continue
			}

			// has default value
			val = r.defaultVal
		}

		// custom filter func
		if r.filterFunc != nil {
			val, err = r.filterFunc(val)
			if err != nil {
				return err
			}

			// save filtered value.
			f.cleanData[field] = val
			continue
		}

		// call built-in filters
		for i, name := range r.filters {
			args := parseArgString(r.filterArgs[i])
			val, err = Apply(name, val, args)
			if err != nil {
				return err
			}
		}

		// save filtered value.
		f.cleanData[field] = val
	}

	return
}

// Fields name get
func (r *Rule) Fields() []string {
	return r.fields
}
