package validate

import (
	"reflect"
	"strings"
)

var (
	// DefaultFieldName for value validate.
	DefaultFieldName = "input"
	// provide value validate
	emptyV = newValValidation()
)

// apply validator to each sub-element of the val(slice, map)
// TODO func Each(val any, rule string)

// Var validating the value by given rule.
// alias of the Val()
func Var(val any, rule string) error {
	return Val(val, rule)
}

// Val quick validating the value by given rule.
// returns error on fail, return nil on check ok.
//
// Usage:
//
//	validate.Val("xyz@mail.com", "required|email")
//
// refer the Validation.StringRule() for parse rule string.
func Val(val any, rule string) error {
	rule = strings.TrimSpace(rule)
	// input empty rule, skip validate
	if rule == "" {
		return nil
	}

	field := DefaultFieldName
	rules := stringSplit(strings.Trim(rule, "|:"), "|")

	es := make(Errors)
	var r *Rule
	var realName string
	for _, validator := range rules {
		validator = strings.Trim(validator, ":")
		if validator == "" {
			continue
		}

		// validator has args. eg: "min:12"
		if strings.ContainsRune(validator, ':') {
			list := stringSplit(validator, ":")
			// reassign value
			validator = list[0]
			realName = ValidatorName(validator)
			switch realName {
			// eg 'regex:\d{4,6}' dont need split args. args is "\d{4,6}"
			case "regexp":
				// v.AddRule(field, validator, list[1])
				r = buildRule(field, validator, realName, []any{list[1]})
				// some special validator. need merge args to one.
			case "enum", "notIn":
				arg := parseArgString(list[1])
				// ev.AddRule(field, validator, arg)
				r = buildRule(field, validator, realName, []any{arg})
			default:
				args := parseArgString(list[1])
				r = buildRule(field, validator, realName, strings2Args(args))
			}
		} else {
			realName = ValidatorName(validator)
			r = buildRule(field, validator, realName, nil)
		}

		// validate value use validator.
		if !r.valueValidate(field, realName, val, emptyV) {
			es.Add(field, validator, r.errorMessage(field, r.validator, emptyV))
			break
		}
	}

	return es.ErrOrNil()
}

// add one Rule
func buildRule(fields, validator, realName string, args []any) *Rule {
	rule := NewRule(fields, validator, args...)

	// init some settings
	rule.realName = realName
	rule.skipEmpty = gOpt.SkipOnEmpty
	// validator name is not "required"
	rule.nameNotRequired = !strings.HasPrefix(realName, "required")

	return rule
}

// create a without context validator's instance.
// see newValidation()
func newValValidation() *Validation {
	v := &Validation{
		trans: NewTranslator(),
		// validator names
		validators: make(map[string]int8, 2),
	}

	// init build in context validator
	v.validatorMetas = make(map[string]*funcMeta, 2)
	ctxValidatorMap := map[string]reflect.Value{
		"required": reflect.ValueOf(v.Required),
	}

	// collect func meta info
	for n, fv := range ctxValidatorMap {
		v.validators[n] = validatorTypeBuiltin
		v.validatorMetas[n] = newFuncMeta(n, true, fv)
	}

	return v
}
