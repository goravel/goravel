package pluralizer

type Ruleset interface {
	AddIrregular(substitutions ...Substitution) Ruleset
	AddUninflected(words ...string) Ruleset
	Regular() Transformations
	Uninflected() Patterns
	Irregular() Substitutions
}
