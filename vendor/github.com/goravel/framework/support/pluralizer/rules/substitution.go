package rules

import (
	"github.com/goravel/framework/contracts/support/pluralizer"
)

var _ pluralizer.Substitution = (*Substitution)(nil)

type Substitution struct {
	from string
	to   string
}

func NewSubstitution(from, to string) *Substitution {
	return &Substitution{
		from: from,
		to:   to,
	}
}

func (r *Substitution) From() string {
	return r.from
}

func (r *Substitution) To() string {
	return r.to
}

func GetFlippedSubstitutions(substitutions ...pluralizer.Substitution) []pluralizer.Substitution {
	flipped := make([]pluralizer.Substitution, len(substitutions))
	for i, sub := range substitutions {
		flipped[i] = NewSubstitution(sub.To(), sub.From())
	}

	return flipped
}
