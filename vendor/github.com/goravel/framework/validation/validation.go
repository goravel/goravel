package validation

import (
	"context"
	"net/url"
	"slices"

	"github.com/gookit/validate"

	validatecontract "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/errors"
)

type Validation struct {
	rules   []validatecontract.Rule
	filters []validatecontract.Filter
}

func NewValidation() *Validation {
	return &Validation{
		rules:   make([]validatecontract.Rule, 0),
		filters: make([]validatecontract.Filter, 0),
	}
}

func (r *Validation) Make(ctx context.Context, data any, rules map[string]string, options ...validatecontract.Option) (validatecontract.Validator, error) {
	if data == nil {
		return nil, errors.ValidationEmptyData
	}
	if len(rules) == 0 {
		return nil, errors.ValidationEmptyRules
	}

	var dataFace validate.DataFace
	var err error
	switch td := data.(type) {
	case validate.DataFace:
		dataFace = td
	case map[string]any:
		if len(td) == 0 {
			return nil, errors.ValidationEmptyData
		}
		dataFace = validate.FromMap(td)
	case url.Values:
		dataFace = validate.FromURLValues(td)
	case map[string][]string:
		dataFace = validate.FromURLValues(td)
	default:
		dataFace, err = validate.FromStruct(data)
		if err != nil {
			return nil, errors.ValidationDataInvalidType
		}
	}

	options = append(options, Rules(rules), CustomRules(r.rules), CustomFilters(r.filters))
	generateOptions := GenerateOptions(options)
	if generateOptions["prepareForValidation"] != nil {
		if err := generateOptions["prepareForValidation"].(func(ctx context.Context, data validatecontract.Data) error)(ctx, NewData(dataFace)); err != nil {
			return nil, err
		}
	}

	v := dataFace.Create()
	AppendOptions(ctx, v, generateOptions)

	return NewValidator(v, dataFace), nil
}

func (r *Validation) AddFilters(filters []validatecontract.Filter) error {
	existFilterNames := r.existFilterNames()
	for _, filter := range filters {
		if slices.Contains(existFilterNames, filter.Signature()) {
			return errors.ValidationDuplicateFilter.Args(filter.Signature())
		}
	}

	r.filters = append(r.filters, filters...)
	return nil
}

func (r *Validation) AddRules(rules []validatecontract.Rule) error {
	existRuleNames := r.existRuleNames()
	for _, rule := range rules {
		if slices.Contains(existRuleNames, rule.Signature()) {
			return errors.ValidationDuplicateRule.Args(rule.Signature())
		}
	}

	r.rules = append(r.rules, rules...)
	return nil
}

func (r *Validation) Rules() []validatecontract.Rule {
	return r.rules
}

func (r *Validation) Filters() []validatecontract.Filter {
	return r.filters
}

func (r *Validation) existRuleNames() []string {
	rules := []string{
		"required",
		"required_if",
		"requiredIf",
		"required_unless",
		"requiredUnless",
		"required_with",
		"requiredWith",
		"required_with_all",
		"requiredWithAll",
		"required_without",
		"requiredWithout",
		"required_without_all",
		"requiredWithoutAll",
		"safe",
		"int",
		"integer",
		"isInt",
		"uint",
		"isUint",
		"bool",
		"isBool",
		"string",
		"isString",
		"float",
		"isFloat",
		"slice",
		"isSlice",
		"in",
		"enum",
		"not_in",
		"notIn",
		"contains",
		"not_contains",
		"notContains",
		"string_contains",
		"stringContains",
		"starts_with",
		"startsWith",
		"ends_with",
		"endsWith",
		"range",
		"between",
		"max",
		"lte",
		"min",
		"gte",
		"eq",
		"equal",
		"isEqual",
		"ne",
		"notEq",
		"notEqual",
		"lt",
		"lessThan",
		"gt",
		"greaterThan",
		"int_eq",
		"intEq",
		"intEqual",
		"len",
		"length",
		"min_len",
		"minLen",
		"minLength",
		"max_len",
		"maxLen",
		"maxLength",
		"email",
		"isEmail",
		"regex",
		"regexp",
		"arr",
		"list",
		"array",
		"isArray",
		"map",
		"isMap",
		"strings",
		"isStrings",
		"ints",
		"isInts",
		"eq_field",
		"eqField",
		"ne_field",
		"neField",
		"gte_field",
		"gtField",
		"gt_field",
		"gteField",
		"lt_field",
		"ltField",
		"lte_field",
		"lteField",
		"file",
		"isFile",
		"image",
		"isImage",
		"mime",
		"mimeType",
		"inMimeTypes",
		"date",
		"isDate",
		"gt_date",
		"gtDate",
		"afterDate",
		"lt_date",
		"ltDate",
		"beforeDate",
		"gte_date",
		"gteDate",
		"afterOrEqualDate",
		"lte_date",
		"lteDate",
		"beforeOrEqualDate",
		"hasWhitespace",
		"ascii",
		"ASCII",
		"isASCII",
		"alpha",
		"isAlpha",
		"alpha_num",
		"alphaNum",
		"isAlphaNum",
		"alpha_dash",
		"alphaDash",
		"isAlphaDash",
		"multi_byte",
		"multiByte",
		"isMultiByte",
		"base64",
		"isBase64",
		"dns_name",
		"dnsName",
		"DNSName",
		"isDNSName",
		"data_uri",
		"dataURI",
		"isDataURI",
		"empty",
		"isEmpty",
		"hex_color",
		"hexColor",
		"isHexColor",
		"hexadecimal",
		"isHexadecimal",
		"json",
		"JSON",
		"isJSON",
		"lat",
		"latitude",
		"isLatitude",
		"lon",
		"longitude",
		"isLongitude",
		"mac",
		"isMAC",
		"num",
		"number",
		"isNumber",
		"cn_mobile",
		"cnMobile",
		"isCnMobile",
		"printableASCII",
		"isPrintableASCII",
		"rgbColor",
		"RGBColor",
		"isRGBColor",
		"full_url",
		"fullUrl",
		"isFullURL",
		"url",
		"URL",
		"isURL",
		"ip",
		"IP",
		"isIP",
		"ipv4",
		"isIPv4",
		"ipv6",
		"isIPv6",
		"cidr",
		"CIDR",
		"isCIDR",
		"CIDRv4",
		"isCIDRv4",
		"CIDRv6",
		"isCIDRv6",
		"uuid",
		"isUUID",
		"uuid3",
		"isUUID3",
		"uuid4",
		"isUUID4",
		"uuid5",
		"isUUID5",
		"filePath",
		"isFilePath",
		"unixPath",
		"isUnixPath",
		"winPath",
		"isWinPath",
		"isbn10",
		"ISBN10",
		"isISBN10",
		"isbn13",
		"ISBN13",
		"isISBN13",
	}
	for _, rule := range r.rules {
		rules = append(rules, rule.Signature())
	}

	return rules
}

func (r *Validation) existFilterNames() []string {
	filters := []string{
		"int",
		"toInt",
		"uint",
		"toUint",
		"int64",
		"toInt64",
		"float",
		"toFloat",
		"bool",
		"toBool",
		"trim",
		"trimSpace",
		"ltrim",
		"trimLeft",
		"rtrim",
		"trimRight",
		"int",
		"integer",
		"lower",
		"lowercase",
		"upper",
		"uppercase",
		"lcFirst",
		"lowerFirst",
		"ucFirst",
		"upperFirst",
		"ucWord",
		"upperWord",
		"camel",
		"camelCase",
		"snake",
		"snakeCase",
		"escapeJs",
		"escapeJS",
		"escapeHtml",
		"escapeHTML",
		"str2ints",
		"strToInts",
		"str2time",
		"strToTime",
		"str2arr",
		"str2array",
		"strToArray",
	}
	for _, filter := range r.filters {
		filters = append(filters, filter.Signature())
	}

	return filters
}
