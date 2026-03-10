package inflector

import (
	"strings"

	"github.com/goravel/framework/contracts/support/pluralizer"
)

type Inflector struct {
	language pluralizer.Language
}

func New(language pluralizer.Language) pluralizer.Inflector {
	return &Inflector{
		language: language,
	}
}

func (r *Inflector) Language() pluralizer.Language {
	return r.language
}

func (r *Inflector) Plural(word string) string {
	return r.inflect(word, r.language.PluralRuleset())
}

func (r *Inflector) SetLanguage(language pluralizer.Language) pluralizer.Inflector {
	r.language = language

	return r
}

func (r *Inflector) Singular(word string) string {
	return r.inflect(word, r.language.SingularRuleset())
}

func (r *Inflector) inflect(word string, ruleset pluralizer.Ruleset) string {
	if word == "" {
		return ""
	}

	for _, pattern := range ruleset.Uninflected() {
		if pattern.Matches(word) {
			return word
		}
	}

	// Check if word is already in target form (To)
	for _, substitution := range ruleset.Irregular() {
		if strings.EqualFold(word, substitution.To()) {
			return word
		}
	}

	// Check if word is in source form (From) and convert to target form (To)
	for _, substitution := range ruleset.Irregular() {
		if strings.EqualFold(word, substitution.From()) {
			return MatchCase(substitution.To(), word)
		}
	}

	for _, transformation := range ruleset.Regular() {
		if result := transformation.Apply(word); result != "" {
			return MatchCase(result, word)
		}
	}

	return word
}
