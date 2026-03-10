package validate

import (
	"reflect"
	"strings"

	"github.com/gookit/filter"
)

/*************************************************************
 * Global filters
 *************************************************************/

var (
	filterValues map[string]reflect.Value
)

// AddFilters add global filters
func AddFilters(m map[string]any) {
	for name, filterFunc := range m {
		AddFilter(name, filterFunc)
	}
}

// AddFilter add global filter to the pkg.
func AddFilter(name string, filterFunc any) {
	if filterValues == nil {
		filterValues = make(map[string]reflect.Value)
	}

	filterValues[name] = checkFilterFunc(name, filterFunc)
}

/*************************************************************
 * filters for current validation
 *************************************************************/

// AddFilters to the Validation
func (v *Validation) AddFilters(m map[string]any) {
	for name, filterFunc := range m {
		v.AddFilter(name, filterFunc)
	}
}

// AddFilter to the Validation.
func (v *Validation) AddFilter(name string, filterFunc any) {
	if v.filterValues == nil {
		v.filterValues = make(map[string]reflect.Value)
	}

	// v.filterFuncs[name] = filterFunc
	v.filterValues[name] = checkFilterFunc(name, filterFunc)
}

// FilterFuncValue get filter by name
func (v *Validation) FilterFuncValue(name string) reflect.Value {
	if fv, ok := v.filterValues[name]; ok {
		return fv
	}

	if fv, ok := filterValues[name]; ok {
		return fv
	}

	return emptyValue
}

// FilterRule add filter rule.
//
// Usage:
//
//	v.FilterRule("name", "trim|lower")
//	v.FilterRule("age", "int")
func (v *Validation) FilterRule(field string, rule string) *FilterRule {
	rule = strings.TrimSpace(rule)
	rules := stringSplit(strings.Trim(rule, "|:"), "|")
	fields := stringSplit(field, ",")

	if len(fields) == 0 || len(rules) == 0 {
		panicf("no enough arguments or contains invalid argument for add filter rule")
	}

	r := newFilterRule(fields)
	r.AddFilters(rules...)
	v.filterRules = append(v.filterRules, r)

	return r
}

// FilterRules add multi filter rules.
func (v *Validation) FilterRules(rules map[string]string) *Validation {
	for field, rule := range rules {
		v.FilterRule(field, rule)
	}
	return v
}

/*************************************************************
 * filtering rule
 *************************************************************/

// FilterRule definition
type FilterRule struct {
	// fields to filter
	fields []string
	// filter name list
	filters []string
	// filter args. { index: "args" }
	filterArgs map[int]string
}

func newFilterRule(fields []string) *FilterRule {
	return &FilterRule{
		fields: fields,
		// init map
		filterArgs: make(map[int]string),
	}
}

// AddFilters add filter(s).
//
// Usage:
//
//	r.AddFilters("int", "str2arr:,")
func (r *FilterRule) AddFilters(filters ...string) *FilterRule {
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
func (r *FilterRule) Apply(v *Validation) (err error) {
	// filter field value
	for _, field := range r.Fields() {
		val, exist, zero := v.tryGet(field)
		if !exist || zero {
			defVal, ok := v.GetDefValue(field)
			// there is also no custom default value
			if !ok {
				continue
			}

			// update source data field value
			newVal, err := v.updateValue(field, defVal)
			if err != nil {
				return err
			}

			// re-set value
			val = newVal

			// dont need check default value
			if !v.CheckDefault {
				v.safeData[field] = newVal // save validated value.
				continue
			}
		}

		// call filters
		for i, name := range r.filters {
			fv := v.FilterFuncValue(name)
			args := parseArgString(r.filterArgs[i])
			if !fv.IsValid() { // is built int filters
				val, err = filter.Apply(name, val, args)
			} else {
				val, err = callCustomFilter(fv, val, args)
			}
			if err != nil {
				return err
			}
		}

		// update source data field value
		newVal, err := v.updateValue(field, val)
		if err != nil {
			return err
		}

		// save filtered value.
		v.filteredData[field] = newVal
	}
	return
}

// Fields name get
func (r *FilterRule) Fields() []string {
	return r.fields
}

func callCustomFilter(fv reflect.Value, val any, args []string) (any, error) {
	var rs []reflect.Value
	if len(args) > 0 {
		rs = CallByValue(fv, buildArgs(val, strings2Args(args))...)
	} else {
		rs = CallByValue(fv, val)
	}

	// return new val.
	if rl := len(rs); rl > 0 {
		val = rs[0].Interface() // `func(val) (newVal)`

		// filter func report error. `func(val) (newVal, error)`
		if rl == 2 {
			err := rs[1].Interface()
			if err != nil {
				return nil, err.(error)
			}
		}
	}

	return val, nil
}
