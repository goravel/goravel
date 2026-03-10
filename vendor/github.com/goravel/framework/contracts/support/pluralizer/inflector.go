package pluralizer

type Inflector interface {
	Language() Language
	Plural(word string) string
	SetLanguage(language Language) Inflector
	Singular(word string) string
}
