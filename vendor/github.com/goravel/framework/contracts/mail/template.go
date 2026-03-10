package mail

type Template interface {
	// Render renders a template with the given data
	Render(path string, data any) (string, error)
}
