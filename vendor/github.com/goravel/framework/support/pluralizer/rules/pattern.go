package rules

import (
	"github.com/goravel/framework/contracts/support/pluralizer"
	"regexp"
)

var _ pluralizer.Pattern = (*Pattern)(nil)

type Pattern struct {
	pattern *regexp.Regexp
}

func NewPattern(pattern string) *Pattern {
	return &Pattern{
		pattern: regexp.MustCompile(pattern),
	}
}

func (r *Pattern) Matches(word string) bool {
	return r.pattern.MatchString(word)
}
