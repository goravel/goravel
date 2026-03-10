package pluralizer

import (
	"github.com/goravel/framework/contracts/support/pluralizer"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/pluralizer/english"
	"github.com/goravel/framework/support/pluralizer/inflector"
	"github.com/goravel/framework/support/pluralizer/rules"
)

var (
	instance         pluralizer.Inflector
	defaultLanguage  = LanguageEnglish
	inflectorFactory = map[string]pluralizer.Inflector{
		"english": inflector.New(english.New()),
	}
)

func init() {
	instance = inflectorFactory[defaultLanguage]
}

func UseLanguage(lang string) error {
	if factory, exists := inflectorFactory[lang]; exists {
		instance = factory
		return nil
	}
	return errors.PluralizerLanguageNotFound.Args(lang)
}

func GetLanguage() pluralizer.Language {
	return instance.Language()
}

func RegisterLanguage(language pluralizer.Language) error {
	if language == nil || language.Name() == "" {
		return errors.PluralizerEmptyLanguageName
	}

	inflectorFactory[language.Name()] = inflector.New(language)
	return nil
}

func getLanguageInstance(lang string) (pluralizer.Language, pluralizer.Inflector, bool) {
	factory, exists := inflectorFactory[lang]
	if !exists {
		return nil, nil, false
	}

	language := factory.Language()
	return language, factory, true
}

func RegisterIrregular(lang string, substitutions ...pluralizer.Substitution) error {
	if len(substitutions) == 0 {
		return errors.PluralizerNoSubstitutionsGiven
	}

	language, factory, exists := getLanguageInstance(lang)
	if !exists {
		return errors.PluralizerLanguageNotFound.Args(lang)
	}

	language.PluralRuleset().AddIrregular(substitutions...)

	flipped := rules.GetFlippedSubstitutions(substitutions...)
	language.SingularRuleset().AddIrregular(flipped...)

	factory.SetLanguage(language)
	return nil
}

func RegisterUninflected(lang string, words ...string) error {
	if len(words) == 0 {
		return errors.PluralizerNoWordsGiven
	}

	language, factory, exists := getLanguageInstance(lang)
	if !exists {
		return errors.PluralizerLanguageNotFound.Args(lang)
	}

	language.PluralRuleset().AddUninflected(words...)
	language.SingularRuleset().AddUninflected(words...)

	factory.SetLanguage(language)
	return nil
}

func RegisterPluralUninflected(lang string, words ...string) error {
	if len(words) == 0 {
		return errors.PluralizerNoWordsGiven
	}

	language, factory, exists := getLanguageInstance(lang)
	if !exists {
		return errors.PluralizerLanguageNotFound.Args(lang)
	}

	language.PluralRuleset().AddUninflected(words...)
	factory.SetLanguage(language)
	return nil
}

func RegisterSingularUninflected(lang string, words ...string) error {
	if len(words) == 0 {
		return errors.PluralizerNoWordsGiven
	}

	language, factory, exists := getLanguageInstance(lang)
	if !exists {
		return errors.PluralizerLanguageNotFound.Args(lang)
	}

	language.SingularRuleset().AddUninflected(words...)
	factory.SetLanguage(language)
	return nil
}

func Plural(word string) string {
	return instance.Plural(word)
}

func Singular(word string) string {
	return instance.Singular(word)
}
