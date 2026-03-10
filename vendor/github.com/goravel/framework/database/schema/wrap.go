package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/goravel/framework/support/collect"
)

var (
	escapeQuoteRegex = regexp.MustCompile(`(\\+)?'`)
	jsonPathRegex    = regexp.MustCompile(`(\[[^]]+])+$`)
	jsonKeyRegex     = regexp.MustCompile(`\[([^]]+)]`)
)

type Wrap struct {
	wrapValue func(string) string
	prefix    string
}

func NewWrap(prefix string) *Wrap {
	return &Wrap{
		prefix: prefix,
	}
}

func (r *Wrap) Column(column string) string {
	if strings.Contains(column, " as ") {
		return r.aliasedValue(column)
	}

	return r.Segments(strings.Split(column, "."))
}

func (r *Wrap) Columns(columns []string) []string {
	formatedColumns := make([]string, len(columns))
	for i, column := range columns {
		formatedColumns[i] = r.Column(column)
	}

	return formatedColumns
}

func (r *Wrap) Columnize(columns []string) string {
	columns = r.Columns(columns)

	return strings.Join(columns, ", ")
}

func (r *Wrap) GetPrefix() string {
	return r.prefix
}

func (r *Wrap) JsonFieldAndPath(column string) (string, string) {
	parts := strings.SplitN(column, "->", 2)
	field := r.Column(parts[0])

	var path string
	if len(parts) > 1 {
		path = ", " + r.JsonPath(parts[1])
	}

	return field, path
}

func (r *Wrap) JsonPath(value string) string {
	value = escapeQuoteRegex.ReplaceAllString(value, "''")
	segments := strings.Split(value, "->")
	for i := range segments {
		if parts := jsonPathRegex.FindString(segments[i]); parts != "" {
			if key := strings.TrimSuffix(segments[i], parts); len(key) > 0 {
				segments[i] = fmt.Sprintf(`"%s"%s`, key, parts)
				continue
			}
			segments[i] = parts
		}
		segments[i] = fmt.Sprintf(`"%s"`, segments[i])
	}

	jsonPath := strings.Join(segments, ".")
	if strings.HasPrefix(jsonPath, "[") {
		return fmt.Sprintf(`'$%s'`, jsonPath)
	}

	return fmt.Sprintf(`'$.%s'`, jsonPath)
}

func (r *Wrap) JsonPathAttributes(path []string, quoter ...string) []string {
	var quote = func(v string) string {
		if _, err := strconv.Atoi(v); err == nil {
			// if it's a number, we don't need to quote it
			return v
		}

		if len(quoter) > 0 {
			return quoter[0] + v + quoter[0]
		}

		return r.Quote(v)
	}

	var result []string
	for i := range path {
		if parts := jsonPathRegex.FindString(path[i]); parts != "" {
			key := strings.TrimSuffix(path[i], parts)
			result = append(result, quote(key))

			matches := jsonKeyRegex.FindAllStringSubmatch(parts, -1)
			for j := range matches {
				if len(matches[j]) > 1 && matches[j][1] != "" {
					result = append(result, quote(matches[j][1]))
				}
			}
			continue
		}
		result = append(result, quote(path[i]))
	}

	return result
}

func (r *Wrap) Not(query string, isNot bool) string {
	if isNot {
		return "not " + query
	}

	return query
}

func (r *Wrap) PrefixArray(prefix string, values []string) []string {
	return collect.Map(values, func(value string, _ int) string {
		return prefix + " " + value
	})
}

func (r *Wrap) Quote(value string) string {
	if value == "" {
		return value
	}

	return fmt.Sprintf("'%s'", value)
}

func (r *Wrap) Quotes(value []string) []string {
	return collect.Map(value, func(v string, _ int) string {
		return r.Quote(v)
	})
}

func (r *Wrap) Segments(segments []string) string {
	for i, segment := range segments {
		if i == 0 && len(segments) > 1 {
			segments[i] = r.Table(segment)
		} else {
			segments[i] = r.Value(segment)
		}
	}

	return strings.Join(segments, ".")
}

func (r *Wrap) SetValueWrapper(wrapper func(string) string) {
	r.wrapValue = wrapper
}

func (r *Wrap) Table(table string) string {
	if strings.Contains(table, " as ") {
		return r.aliasedTable(table)
	}
	if strings.Contains(table, ".") {
		lastDotIndex := strings.LastIndex(table, ".")

		return r.Value(table[:lastDotIndex]) + "." + r.Value(r.prefix+table[lastDotIndex+1:])
	}

	return r.Value(r.prefix + table)
}

func (r *Wrap) Value(value string) string {
	if value != "*" {
		if r.wrapValue != nil {
			return r.wrapValue(value)
		}

		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}

	return value
}

func (r *Wrap) aliasedTable(table string) string {
	segments := strings.Split(table, " as ")

	return r.Table(segments[0]) + " as " + r.Value(r.prefix+segments[1])
}

func (r *Wrap) aliasedValue(value string) string {
	segments := strings.Split(value, " as ")

	return r.Column(segments[0]) + " as " + r.Value(r.prefix+segments[1])
}
