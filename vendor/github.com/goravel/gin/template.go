package gin

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
)

type Delims struct {
	Left  string
	Right string
}

type RenderOptions struct {
	Delims  *Delims
	FuncMap template.FuncMap
}

func NewTemplate(options RenderOptions) (*render.HTMLProduction, error) {
	instance := template.New("")
	if options.Delims != nil {
		instance.Delims(options.Delims.Left, options.Delims.Right)
	}
	if options.FuncMap != nil {
		instance.Funcs(options.FuncMap)
	}

	dir := path.Resource("views")
	if !file.Exists(dir) {
		return nil, nil
	}

	var files []string
	if err := filepath.Walk(dir, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, fullPath)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	tmpl := template.Must(instance.ParseFiles(files...))

	return &render.HTMLProduction{Template: tmpl}, nil
}

// DefaultTemplate creates a TemplateRender instance with default options.
func DefaultTemplate() (*render.HTMLProduction, error) {
	return NewTemplate(RenderOptions{})
}
