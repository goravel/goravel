// Package english provides a comprehensive set of rules for English pluralization
// and singularization. It's the result of a deep dive into English morphology,
// designed to be highly accurate by respecting the complex hierarchy of its rules.
//
// The core philosophy is that inflection isn't a single step, but a cascade: a word
// is first checked against a list of truly invariable words (uninflected), then
// against a list of unique irregular words, and only then is it processed by the
// general pattern-based rules.
//
// A key principle in this library is the careful distinction between nouns that are
// sometimes uncountable (like "work") and those that are almost always invariable
// (like "information"). To allow for correct pluralization of words like "work"
// into "works" or "permission" into "permissions", the uninflected list is
// intentionally kept strict and concise.
//
// Sources and references used for curation:
//
//   - English Plurals (Wikipedia): For the foundational rules of regular,
//     irregular, and loanword pluralization.
//     https://en.wikipedia.org/wiki/English_plurals
//
//   - Mass Nouns / Uncountable Nouns (Wikipedia): To understand the principles
//     behind uncountable nouns and curate the uninflected list.
//     https://en.wikipedia.org/wiki/Mass_noun
//
//   - Plurale & Singulare Tantum (Wikipedia): For handling plural-only and
//     singular-only nouns like "scissors" and "news".
//     https://en.wikipedia.org/wiki/Plurale_tantum
//     https://en.wikipedia.org/wiki/Singulare_tantum
//
//   - Grammarist - Uncountable Nouns: Provided a practical list that helped
//     refine the distinction between contextually and strictly uncountable nouns.
//     https://grammarist.com/grammar/uncountable-nouns/
//
//   - Wiktionary, the free dictionary: Used extensively for cross-referencing
//     the plural forms and usage notes of individual, ambiguous words.
//     https://en.wiktionary.org/
//
//   - Doctrine Inflector (GitHub): For analysis of a robust, widely-used, and
//     battle-tested inflection library.
//     https://github.com/doctrine/inflector
package english

import (
	"github.com/goravel/framework/contracts/support/pluralizer"
	"github.com/goravel/framework/support/pluralizer/rules"
)

var _ pluralizer.Language = (*Language)(nil)

type Language struct {
	singularRuleset pluralizer.Ruleset
	pluralRuleset   pluralizer.Ruleset
}

func New() *Language {
	return &Language{
		singularRuleset: newEnglishSingularRuleset(),
		pluralRuleset:   newEnglishPluralRuleset(),
	}
}

func (r *Language) Name() string {
	return "english"
}

func (r *Language) SingularRuleset() pluralizer.Ruleset {
	return r.singularRuleset
}

func (r *Language) PluralRuleset() pluralizer.Ruleset {
	return r.pluralRuleset
}

func newEnglishPluralRuleset() pluralizer.Ruleset {
	uninflected := getUninflectedDefault()
	uninflected = append(uninflected,
		rules.NewPattern(`(?i)media$`),
		rules.NewPattern(`(?i)people$`),
		rules.NewPattern(`(?i)trivia$`),
		rules.NewPattern(`(?i)\w+ware$`),
	)

	irregular := pluralizer.Substitutions{
		rules.NewSubstitution("atlas", "atlases"),
		rules.NewSubstitution("brother", "brothers"),
		rules.NewSubstitution("brother-in-law", "brothers-in-law"),
		rules.NewSubstitution("cafe", "cafes"),
		rules.NewSubstitution("chateau", "chateaux"),
		rules.NewSubstitution("child", "children"),
		rules.NewSubstitution("cookie", "cookies"),
		rules.NewSubstitution("corpus", "corpora"),
		rules.NewSubstitution("criterion", "criteria"),
		rules.NewSubstitution("curriculum", "curricula"),
		rules.NewSubstitution("daughter-in-law", "daughters-in-law"),
		rules.NewSubstitution("demo", "demos"),
		rules.NewSubstitution("die", "dice"),
		rules.NewSubstitution("domino", "dominoes"),
		rules.NewSubstitution("echo", "echoes"),
		rules.NewSubstitution("father-in-law", "fathers-in-law"),
		rules.NewSubstitution("foe", "foes"),
		rules.NewSubstitution("foot", "feet"),
		rules.NewSubstitution("fungus", "fungi"),
		rules.NewSubstitution("genie", "genies"),
		rules.NewSubstitution("genus", "genera"),
		rules.NewSubstitution("goose", "geese"),
		rules.NewSubstitution("graffito", "graffiti"),
		rules.NewSubstitution("hippopotamus", "hippopotami"),
		rules.NewSubstitution("iris", "irises"),
		rules.NewSubstitution("larva", "larvae"),
		rules.NewSubstitution("leaf", "leaves"),
		rules.NewSubstitution("lens", "lenses"),
		rules.NewSubstitution("loaf", "loaves"),
		rules.NewSubstitution("louse", "lice"),
		rules.NewSubstitution("man", "men"),
		rules.NewSubstitution("memorandum", "memoranda"),
		rules.NewSubstitution("mongoose", "mongooses"),
		rules.NewSubstitution("mother-in-law", "mothers-in-law"),
		rules.NewSubstitution("motto", "mottoes"),
		rules.NewSubstitution("mouse", "mice"),
		rules.NewSubstitution("move", "moves"),
		rules.NewSubstitution("mythos", "mythoi"),
		rules.NewSubstitution("nucleus", "nuclei"),
		rules.NewSubstitution("oasis", "oases"),
		rules.NewSubstitution("octopus", "octopuses"),
		rules.NewSubstitution("opus", "opuses"),
		rules.NewSubstitution("ox", "oxen"),
		rules.NewSubstitution("passer-by", "passers-by"),
		rules.NewSubstitution("passerby", "passersby"),
		rules.NewSubstitution("phenomenon", "phenomena"),
		rules.NewSubstitution("person", "people"),
		rules.NewSubstitution("plateau", "plateaux"),
		rules.NewSubstitution("runner-up", "runners-up"),
		rules.NewSubstitution("sex", "sexes"),
		rules.NewSubstitution("sister-in-law", "sisters-in-law"),
		rules.NewSubstitution("soliloquy", "soliloquies"),
		rules.NewSubstitution("son-in-law", "sons-in-law"),
		rules.NewSubstitution("syllabus", "syllabi"),
		rules.NewSubstitution("testis", "testes"),
		rules.NewSubstitution("thief", "thieves"),
		rules.NewSubstitution("tooth", "teeth"),
		rules.NewSubstitution("tornado", "tornadoes"),
		rules.NewSubstitution("volcano", "volcanoes"),
		rules.NewSubstitution("woman", "women"),
		rules.NewSubstitution("zombie", "zombies"),
	}

	regular := pluralizer.Transformations{
		rules.NewTransformation(`(?i)(quiz)$`, `${1}zes`),
		rules.NewTransformation(`(?i)(matr|vert|ind)(ix|ex)$`, `${1}ices`),
		rules.NewTransformation(`(?i)(x|ch|ss|sh|z)$`, `${1}es`),
		rules.NewTransformation(`(?i)([^aeiouy]|qu)y$`, `${1}ies`),
		rules.NewTransformation(`(?i)(hive|gulf)$`, `${1}s`),
		rules.NewTransformation(`(?i)(?:([^f])fe|([lr])f)$`, `${1}${2}ves`),
		rules.NewTransformation(`(?i)sis$`, `ses`),
		rules.NewTransformation(`(?i)([ti])um$`, `${1}a`),
		rules.NewTransformation(`(?i)(tax)on$`, `${1}a`),
		rules.NewTransformation(`(?i)(buffal|her|potat|tomat|volcan)o$`, `${1}oes`),
		rules.NewTransformation(`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|vir)us$`, `${1}i`),
		rules.NewTransformation(`(?i)us$`, `uses`),
		rules.NewTransformation(`(?i)(alias|status)$`, `${1}es`),
		rules.NewTransformation(`(?i)(analys|ax|cris|test|thes)is$`, `${1}es`),
		rules.NewTransformation(`(?i)s$`, `s`),
		rules.NewTransformation(`(?i)$`, `s`),
	}

	return rules.NewRuleset(regular, uninflected, irregular)
}

func newEnglishSingularRuleset() pluralizer.Ruleset {
	uninflected := getUninflectedDefault()
	uninflected = append(uninflected,
		rules.NewPattern(`(?i).*ss$`),
		rules.NewPattern(`(?i)athletics$`),
		rules.NewPattern(`(?i)data$`),
		rules.NewPattern(`(?i)electronics$`),
		rules.NewPattern(`(?i)genetics$`),
		rules.NewPattern(`(?i)graphics$`),
		rules.NewPattern(`(?i)jeans$`),
		rules.NewPattern(`(?i)mathematics$`),
		rules.NewPattern(`(?i)news$`),
		rules.NewPattern(`(?i)pliers$`),
		rules.NewPattern(`(?i)politics$`),
		rules.NewPattern(`(?i)scissors$`),
		rules.NewPattern(`(?i)shorts$`),
		rules.NewPattern(`(?i)trousers$`),
		rules.NewPattern(`(?i)trivia$`),
	)

	irregular := pluralizer.Substitutions{}
	for _, sub := range newEnglishPluralRuleset().Irregular() {
		irregular = append(irregular, rules.NewSubstitution(sub.To(), sub.From()))
	}

	regular := pluralizer.Transformations{
		rules.NewTransformation(`(?i)(s)tatuses$`, `${1}tatus`),
		rules.NewTransformation(`(?i)(quiz)zes$`, `${1}`),
		rules.NewTransformation(`(?i)(matr)ices$`, `${1}ix`),
		rules.NewTransformation(`(?i)(vert|ind)ices$`, `${1}ex`),
		rules.NewTransformation(`(?i)^(ox)en`, `${1}`),
		rules.NewTransformation(`(?i)(alias|status)(es)?$`, `${1}`),
		rules.NewTransformation(`(?i)(buffal|her|potat|tomat|volcan)oes$`, `${1}o`),
		rules.NewTransformation(`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$`, `${1}us`),
		rules.NewTransformation(`(?i)([ftw]ax)es$`, `${1}`),
		rules.NewTransformation(`(?i)(analys|ax|cris|test|thes)es$`, `${1}is`),
		rules.NewTransformation(`(?i)(shoe|slave)s$`, `${1}`),
		rules.NewTransformation(`(?i)(o)es$`, `${1}`),
		rules.NewTransformation(`(?i)ouses$`, `ouse`),
		rules.NewTransformation(`(?i)([^a])uses$`, `${1}us`),
		rules.NewTransformation(`(?i)([ml])ice$`, `${1}ouse`),
		rules.NewTransformation(`(?i)(x|ch|ss|sh|z)es$`, `${1}`),
		rules.NewTransformation(`(?i)(m)ovies$`, `${1}ovie`),
		rules.NewTransformation(`(?i)(s)eries$`, `${1}eries`),
		rules.NewTransformation(`(?i)([^aeiouy]|qu)ies$`, `${1}y`),
		rules.NewTransformation(`(?i)([lr])ves$`, `${1}f`),
		rules.NewTransformation(`(?i)(tive)s$`, `${1}`),
		rules.NewTransformation(`(?i)(hive)s$`, `${1}`),
		rules.NewTransformation(`(?i)([^fo])ves$`, `${1}fe`),
		rules.NewTransformation(`(?i)(^analy)ses$`, `${1}sis`),
		rules.NewTransformation(`(?i)(analy|diagno|^ba|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$`, `${1}${2}sis`),
		rules.NewTransformation(`(?i)(tax)a$`, `${1}on`),
		rules.NewTransformation(`(?i)(c)riteria$`, `${1}riterion`),
		rules.NewTransformation(`(?i)([ti])a$`, `${1}um`),
		rules.NewTransformation(`(?i)eaus$`, `eau`),
		rules.NewTransformation(`(?i)s$`, ``),
	}

	return rules.NewRuleset(regular, uninflected, irregular)
}

// getUninflectedDefault returns a list of nouns that are truly uninflected.
// This list is carefully curated to contain only pure mass nouns and words with
// identical singular/plural forms, preventing incorrect behavior for words that
// can be both countable and uncountable (e.g., "permission" -> "permissions").
func getUninflectedDefault() []pluralizer.Pattern {
	return []pluralizer.Pattern{
		// Truly uncountable nouns
		rules.NewPattern(`(?i)advice$`),
		rules.NewPattern(`(?i)baggage$`),
		rules.NewPattern(`(?i)butter$`),
		rules.NewPattern(`(?i)clothing$`),
		rules.NewPattern(`(?i)coal$`),
		rules.NewPattern(`(?i)debris$`),
		rules.NewPattern(`(?i)education$`),
		rules.NewPattern(`(?i)equipment$`),
		rules.NewPattern(`(?i)evidence$`),
		rules.NewPattern(`(?i)feedback$`),
		rules.NewPattern(`(?i)food$`),
		rules.NewPattern(`(?i)furniture$`),
		rules.NewPattern(`(?i)homework$`),
		rules.NewPattern(`(?i)information$`),
		rules.NewPattern(`(?i)knowledge$`),
		rules.NewPattern(`(?i)leather$`),
		rules.NewPattern(`(?i)luggage$`),
		rules.NewPattern(`(?i)money$`),
		rules.NewPattern(`(?i)music$`),
		rules.NewPattern(`(?i)plankton$`),
		rules.NewPattern(`(?i)progress$`),
		rules.NewPattern(`(?i)rain$`),
		rules.NewPattern(`(?i)research$`),
		rules.NewPattern(`(?i)rice$`),
		rules.NewPattern(`(?i)sand$`),
		rules.NewPattern(`(?i)spam$`),
		rules.NewPattern(`(?i)sugar$`),
		rules.NewPattern(`(?i)traffic$`),
		rules.NewPattern(`(?i)water$`),
		rules.NewPattern(`(?i)weather$`),
		rules.NewPattern(`(?i)wheat$`),
		rules.NewPattern(`(?i)wood$`),
		rules.NewPattern(`(?i)wool$`),

		// Nouns with identical singular and plural forms (zero-plurals)
		rules.NewPattern(`(?i)aircraft$`),
		rules.NewPattern(`(?i)bison$`),
		rules.NewPattern(`(?i)buffalo$`),
		rules.NewPattern(`(?i)chassis$`),
		rules.NewPattern(`(?i)cod$`),
		rules.NewPattern(`(?i)corps$`),
		rules.NewPattern(`(?i)deer$`),
		rules.NewPattern(`(?i)fish$`),
		rules.NewPattern(`(?i)flounder$`),
		rules.NewPattern(`(?i)jedi$`),
		rules.NewPattern(`(?i)mackerel$`),
		rules.NewPattern(`(?i)moose$`),
		rules.NewPattern(`(?i)offspring$`),
		rules.NewPattern(`(?i)pokemon$`),
		rules.NewPattern(`(?i)salmon$`),
		rules.NewPattern(`(?i)series$`),
		rules.NewPattern(`(?i)sheep$`),
		rules.NewPattern(`(?i)shrimp$`),
		rules.NewPattern(`(?i)species$`),
		rules.NewPattern(`(?i)swine$`),
		rules.NewPattern(`(?i)trout$`),
		rules.NewPattern(`(?i)tuna$`),

		// Plural-only nouns (pluralia tantum)
		rules.NewPattern(`(?i)belongings$`),
		rules.NewPattern(`(?i)binoculars$`),
		rules.NewPattern(`(?i)cattle$`),
		rules.NewPattern(`(?i)clothes$`),
		rules.NewPattern(`(?i)congratulations$`),
		rules.NewPattern(`(?i)jeans$`),
		rules.NewPattern(`(?i)pants$`),
		rules.NewPattern(`(?i)pliers$`),
		rules.NewPattern(`(?i)police$`),
		rules.NewPattern(`(?i)scissors$`),
		rules.NewPattern(`(?i)shorts$`),
		rules.NewPattern(`(?i)thanks$`),
		rules.NewPattern(`(?i)trousers$`),

		// Singular-only nouns that end in -s (singularia tantum)
		rules.NewPattern(`(?i)athletics$`),
		rules.NewPattern(`(?i)billiards$`),
		rules.NewPattern(`(?i)diabetes$`),
		rules.NewPattern(`(?i)economics$`),
		rules.NewPattern(`(?i)ethics$`),
		rules.NewPattern(`(?i)gallows$`),
		rules.NewPattern(`(?i)gymnastics$`),
		rules.NewPattern(`(?i)innings$`),
		rules.NewPattern(`(?i)linguistics$`),
		rules.NewPattern(`(?i)mathematics$`),
		rules.NewPattern(`(?i)measles$`),
		rules.NewPattern(`(?i)mumps$`),
		rules.NewPattern(`(?i)news$`),
		rules.NewPattern(`(?i)nexus$`),
		rules.NewPattern(`(?i)physics$`),
		rules.NewPattern(`(?i)politics$`),
		rules.NewPattern(`(?i)rabies$`),
		rules.NewPattern(`(?i)rickets$`),
		rules.NewPattern(`(?i)shingles$`),
		rules.NewPattern(`(?i)statistics$`),

		// Other special cases
		rules.NewPattern(`(?i)audio$`),
		rules.NewPattern(`(?i)data$`),
		rules.NewPattern(`(?i)emoji$`),
		rules.NewPattern(`(?i)metadata$`),
		rules.NewPattern(`(?i)sms$`),
		rules.NewPattern(`(?i)staff$`),
	}
}
