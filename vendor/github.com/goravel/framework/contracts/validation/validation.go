package validation

import "context"

type Option func(map[string]any)

type Validation interface {
	// Make create a new validator instance.
	Make(ctx context.Context, data any, rules map[string]string, options ...Option) (Validator, error)
	// AddFilters add the custom filters.
	AddFilters([]Filter) error
	// AddRules add the custom rules.
	AddRules([]Rule) error
	// Rules get the custom rules.
	Rules() []Rule
	// Filters get the custom filters.
	Filters() []Filter
}

type Validator interface {
	// Bind the data to the validation.
	Bind(ptr any) error
	// Errors get the validation errors.
	Errors() Errors
	// Fails determine if the validation fails.
	Fails() bool
}

type Errors interface {
	// One gets the first error message for a given field.
	One(key ...string) string
	// Get gets all the error messages for a given field.
	Get(key string) map[string]string
	// All gets all the error messages.
	All() map[string]map[string]string
	// Has checks if there are any error messages for a given field.
	Has(key string) bool
}

type Data interface {
	// Get the value from the given key.
	Get(key string) (val any, exist bool)
	// Set the value for a given key.
	Set(key string, val any) error
}

type Rule interface {
	// Signature set the unique signature of the rule.
	Signature() string
	// Passes determine if the validation rule passes.
	Passes(ctx context.Context, data Data, val any, options ...any) bool
	// Message gets the validation error message.
	Message(ctx context.Context) string
}

type Filter interface {
	// Signature sets the unique signature of the filter.
	Signature() string

	// Handle defines the filter function to apply.
	//
	// The Handle method should return a function that processes an input and
	// returns a transformed value. The function can either return the
	// transformed value alone or a tuple of the transformed value and an error.
	// The input to the filter function is flexible: the first input is the value
	// of the key on which the filter is applied, and the rest of the inputs are
	// the arguments passed to the filter.
	//
	// Example usages:
	//
	// 1. Return only the transformed value:
	//    func (val string) int {
	//        // conversion logic
	//        return 1
	//    }
	//
	// 2. Return the transformed value and an error:
	//    func (val int) (int, error) {
	//        // conversion logic with error handling
	//        return 1, nil
	//    }
	//
	// 3. Take additional arguments:
	//    func (val string, def ...string) string {
	//        if val == "" && len(def) > 0 {
	//            return def[0]
	//        }
	//        return val
	//    }
	//
	Handle(ctx context.Context) any
}
