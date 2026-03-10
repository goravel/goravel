package rules

import "github.com/goravel/framework/contracts/support/pluralizer"

var _ pluralizer.Ruleset = (*Ruleset)(nil)

type Ruleset struct {
	regular     pluralizer.Transformations
	uninflected pluralizer.Patterns
	irregular   pluralizer.Substitutions
}

func NewRuleset(regular pluralizer.Transformations, uninflected pluralizer.Patterns, irregular pluralizer.Substitutions) *Ruleset {
	return &Ruleset{
		regular:     regular,
		uninflected: uninflected,
		irregular:   irregular,
	}
}

func (r *Ruleset) AddIrregular(substitutions ...pluralizer.Substitution) pluralizer.Ruleset {
	r.irregular = append(substitutions, r.irregular...)
	return r
}

func (r *Ruleset) AddUninflected(words ...string) pluralizer.Ruleset {
	if len(words) == 0 {
		return r
	}

	patterns := make([]pluralizer.Pattern, len(words))
	for i, word := range words {
		patterns[i] = NewPattern(word)
	}

	r.uninflected = append(patterns, r.uninflected...)
	return r
}

func (r *Ruleset) Regular() pluralizer.Transformations {
	return r.regular
}

func (r *Ruleset) Uninflected() pluralizer.Patterns {
	return r.uninflected
}

func (r *Ruleset) Irregular() pluralizer.Substitutions {
	return r.irregular
}
