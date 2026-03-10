package validation

import (
	"context"
	"strings"

	"github.com/gookit/validate"

	contractsvalidation "github.com/goravel/framework/contracts/validation"
)

func Rules(rules map[string]string) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(rules) > 0 {
			options["rules"] = rules
		}
	}
}

func Filters(filters map[string]string) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(filters) > 0 {
			options["filters"] = filters
		}
	}
}

func CustomFilters(filters []contractsvalidation.Filter) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(filters) > 0 {
			options["customFilters"] = filters
		}
	}
}

func CustomRules(rules []contractsvalidation.Rule) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(rules) > 0 {
			options["customRules"] = rules
		}
	}
}

func Messages(messages map[string]string) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(messages) > 0 {
			options["messages"] = messages
		}
	}
}

func Attributes(attributes map[string]string) contractsvalidation.Option {
	return func(options map[string]any) {
		if len(attributes) > 0 {
			options["attributes"] = attributes
		}
	}
}

func PrepareForValidation(prepare func(ctx context.Context, data contractsvalidation.Data) error) contractsvalidation.Option {
	return func(options map[string]any) {
		options["prepareForValidation"] = prepare
	}
}

func GenerateOptions(options []contractsvalidation.Option) map[string]any {
	realOptions := make(map[string]any)
	for _, option := range options {
		option(realOptions)
	}

	return realOptions
}

func AppendOptions(ctx context.Context, validation *validate.Validation, options map[string]any) {
	if options["rules"] != nil {
		rules := options["rules"].(map[string]string)
		for key, value := range rules {
			validation.StringRule(key, value)
		}
	}

	if options["filters"] != nil {
		filters, ok := options["filters"].(map[string]string)
		if ok {
			validation.FilterRules(filters)
		}
	}

	if options["messages"] != nil {
		messages := options["messages"].(map[string]string)
		for key, value := range messages {
			messages[key] = strings.ReplaceAll(value, ":attribute", "{field}")
		}
		validation.AddMessages(messages)
	}

	if options["attributes"] != nil && len(options["attributes"].(map[string]string)) > 0 {
		validation.AddTranslates(options["attributes"].(map[string]string))
	}

	if options["customRules"] != nil {
		customRules := options["customRules"].([]contractsvalidation.Rule)
		for _, customRule := range customRules {
			validation.AddMessages(map[string]string{
				customRule.Signature(): strings.ReplaceAll(customRule.Message(ctx), ":attribute", "{field}"),
			})
			validation.AddValidator(customRule.Signature(), func(val any, options ...any) bool {
				return customRule.Passes(ctx, validation, val, options...)
			})
		}
	}

	if options["customFilters"] != nil {
		customFilters := options["customFilters"].([]contractsvalidation.Filter)
		for _, customFilter := range customFilters {
			filterFunc := customFilter.Handle(ctx)
			if filterFunc == nil {
				continue
			}

			validation.AddFilter(customFilter.Signature(), filterFunc)
		}
	}

	validation.Trans().FieldMap()
}
