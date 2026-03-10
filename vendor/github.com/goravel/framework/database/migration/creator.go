package migration

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/carbon"
)

type Creator struct {
}

func NewCreator() *Creator {
	return &Creator{}
}

// GetStub Get the migration stub file.
func (r *Creator) GetStub(table string, create bool) string {
	if table == "" {
		return Stubs{}.Empty()
	}

	if create {
		return Stubs{}.Create()
	}

	return Stubs{}.Update()
}

type StubData struct {
	FacadesPackage string
	FacadesImport  string
	Package        string
	SchemaFields   []string
	Signature      string
	StructName     string
	Table          string
}

// PopulateStub Populate the place-holders in the migration stub.
func (r *Creator) PopulateStub(stub string, data StubData) (string, error) {
	tmpl, err := template.New("stub").Parse(stub)
	if err != nil {
		return "", errors.TemplateFailedToParse.Args(err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", errors.TemplateFailedToExecute.Args(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", errors.TemplateFailedToFormatGoCode.Args(err)
	}

	return string(formatted), nil
}

// GetPath Get the full path to the migration.
func (r *Creator) GetPath(name string) string {
	pwd, _ := os.Getwd()

	return filepath.Join(pwd, support.Config.Paths.Migrations, name+".go")
}

// GetFileName Get the full path to the migration.
func (r *Creator) GetFileName(name string) string {
	return fmt.Sprintf("%s_%s", carbon.Now().ToShortDateTimeString(), name)
}
