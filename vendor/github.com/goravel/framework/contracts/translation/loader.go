package translation

type Loader interface {
	// Load the messages for the given locale.
	Load(locale string, group string) (map[string]any, error)
}
