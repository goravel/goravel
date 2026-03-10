package pluralizer

type Language interface {
	Name() string
	SingularRuleset() Ruleset
	PluralRuleset() Ruleset
}

type Transformation interface {
	Apply(word string) string
}

type Pattern interface {
	Matches(word string) bool
}

type Substitution interface {
	From() string
	To() string
}

type Transformations []Transformation

type Patterns []Pattern

type Substitutions []Substitution
