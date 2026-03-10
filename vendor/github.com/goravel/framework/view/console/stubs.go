package console

type Stubs struct {
}

func (r Stubs) View() string {
	return `{{ define "DummyDefinition" }}
<h1>Welcome</h1>
{{ end }}
`
}
