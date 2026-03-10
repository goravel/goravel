package paths

import (
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/goravel/framework/contracts/packages"
	"github.com/goravel/framework/support"
)

type Paths struct {
	mainPath string
}

func NewPaths(mainPath string) *Paths {
	return &Paths{mainPath: mainPath}
}

func (r *Paths) App() packages.Path {
	return NewPath(support.Config.Paths.App, r.mainPath, false)
}

// Bootstrap returns the bootstrap package path, eg: goravel/bootstrap.
func (r *Paths) Bootstrap() packages.Path {
	return NewPath(support.Config.Paths.Bootstrap, r.mainPath, false)
}

// Config returns the config package path, eg: goravel/config.
func (r *Paths) Config() packages.Path {
	return NewPath(support.Config.Paths.Config, r.mainPath, false)
}

func (r *Paths) Database() packages.Path {
	return NewPath(support.Config.Paths.Database, r.mainPath, false)
}

// Facades returns the facades package path, eg: goravel/app/facades.
func (r *Paths) Facades() packages.Path {
	return NewPath(support.Config.Paths.Facades, r.mainPath, false)
}

func (r *Paths) Lang() packages.Path {
	return NewPath(support.Config.Paths.Lang, r.mainPath, false)
}

// Main returns the main package path, eg: github.com/goravel/goravel.
func (r *Paths) Main() packages.Path {
	return NewPath("", r.mainPath, false)
}

func (r *Paths) Migrations() packages.Path {
	return NewPath(support.Config.Paths.Migrations, r.mainPath, false)
}

func (r *Paths) Models() packages.Path {
	return NewPath(support.Config.Paths.Models, r.mainPath, false)
}

// Module returns the module path of the package, eg: github.com/goravel/framework/auth.
func (r *Paths) Module() packages.Path {
	var p string
	if info, ok := debug.ReadBuildInfo(); ok && strings.HasSuffix(info.Path, "setup") {
		p = path.Dir(info.Path)
	}

	return NewPath(p, r.mainPath, true)
}

func (r *Paths) Public() packages.Path {
	return NewPath(support.Config.Paths.Public, r.mainPath, false)
}

func (r *Paths) Resources() packages.Path {
	return NewPath(support.Config.Paths.Resources, r.mainPath, false)
}

// Routes returns the routes package path, eg: goravel/routes.
func (r *Paths) Routes() packages.Path {
	return NewPath(support.Config.Paths.Routes, r.mainPath, false)
}

func (r *Paths) Storage() packages.Path {
	return NewPath(support.Config.Paths.Storage, r.mainPath, false)
}

// Tests returns the tests package path, eg: goravel/tests.
func (r *Paths) Tests() packages.Path {
	return NewPath(support.Config.Paths.Tests, r.mainPath, false)
}

func (r *Paths) Views() packages.Path {
	return NewPath(support.Config.Paths.Views, r.mainPath, false)
}

type Path struct {
	main     string
	path     string
	isModule bool
}

func NewPath(path, main string, isModule bool) *Path {
	return &Path{path: path, main: main, isModule: isModule}
}

func (r *Path) Abs(paths ...string) string {
	paths = append(toSlice(r.path), paths...)

	return Abs(paths...)
}

// Package returns the sub-package name, or the main package name if no sub-package path is specified.
// For example, if r.path is "app/http/controllers", it returns "controllers".
// If r.path is empty, it returns the last component of r.main.
func (r *Path) Package() string {
	p := pkg(r.path)

	if p == "" {
		return pkg(r.main)
	}

	return p
}

// Import returns the sub-package import path, or the main package import path if no sub-package path is specified.
// For example, if r.path is "app/http/controllers" and r.main is "github.com/goravel/goravel",
// it returns "goravel/app/http/controllers". If r.path is empty, it returns "goravel".
// The path will be returned directly if it starts with "github.com/goravel/framework/", given it's a framework sub-package.
func (r *Path) Import() string {
	mainSlice := toSlice(r.main)
	mainImport := mainSlice[len(mainSlice)-1]

	if r.path != "" {
		if r.isModule {
			return r.path
		}

		pathSlice := toSlice(r.path)
		importSlice := append([]string{mainImport}, pathSlice...)

		return strings.Join(importSlice, "/")
	}

	return mainImport
}

func (r *Path) String(paths ...string) string {
	paths = append(toSlice(r.path), paths...)

	return filepath.Join(paths...)
}

func Abs(paths ...string) string {
	paths = append([]string{support.RelativePath}, paths...)
	path := filepath.Join(paths...)
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

// pkg extracts the last component of a file path string.
// For example, "app/http/controllers" returns "controllers".
func pkg(path string) string {
	s := toSlice(path)

	if len(s) == 0 {
		return ""
	}

	return s[len(s)-1]
}

// toSlice converts a file path string into a slice of its components,
// handling both forward slashes and backslashes, and trimming leading/trailing slashes.
// For example, "app/http/controllers" becomes []string{"app", "http", "controllers"}.
func toSlice(path string) []string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}
