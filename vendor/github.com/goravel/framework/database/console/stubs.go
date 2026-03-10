package console

type Stubs struct {
}

func (r Stubs) Model() string {
	return `package {{.PackageName}}

{{if or .Imports (and (not .Embeds) (not .Fields)) -}}
import (
{{- range $path, $_ := .Imports }}
	"{{$path}}"
{{- end }}
{{- if and (not .Embeds) (not .Fields) }}
	"github.com/goravel/framework/database/orm"
{{- end }}
)
{{- end }}

type {{.StructName}} struct {
{{- if and (not .Embeds) (not .Fields) }}
	orm.Model
{{- end }}
{{- range .Embeds }}
	{{.}}
{{- end }}
{{- range .Fields }}
	{{.}}
{{- end }}
}

{{- if .TableNameMethod }}

{{.TableNameMethod}}
{{- end }}
`
}

func (r Stubs) Observer() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/database/orm"
)


type DummyObserver struct{}

func (u *DummyObserver) Created(event orm.Event) error {
	return nil
}

func (u *DummyObserver) Updated(event orm.Event) error {
	return nil
}

func (u *DummyObserver) Deleted(event orm.Event) error {
	return nil
}

func (u *DummyObserver) ForceDeleted(event orm.Event) error {
	return nil
}
`
}

func (r Stubs) Seeder() string {
	return `package DummyPackage
	
type DummySeeder struct {
}

// Signature The name and signature of the seeder.
func (s *DummySeeder) Signature() string {
	return "DummySignature"
}

// Run executes the seeder logic.
func (s *DummySeeder) Run() error {
	return nil
}
`
}

func (r Stubs) Factory() string {
	return `package DummyPackage

type DummyFactory struct {
}

// Definition Define the model's default state.
func (f *DummyFactory) Definition() map[string]any {
     return nil
}
`
}
