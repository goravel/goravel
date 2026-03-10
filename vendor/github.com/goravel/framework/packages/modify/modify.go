package modify

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/spf13/cast"
	"golang.org/x/tools/imports"

	"github.com/goravel/framework/contracts/packages/match"
	"github.com/goravel/framework/contracts/packages/modify"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
	supportfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
	"github.com/goravel/framework/support/str"
)

func RegisterMigration(pkg, migration string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return AddMigration(pkg, migration)
	})
}

func RegisterProvider(pkg, provider string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return AddProvider(pkg, provider)
	})
}

func RegisterRoute(pkg, route string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return AddRoute(pkg, route)
	})
}

func Call(fn func(options []modify.Option) error) modify.Apply {
	return &callModifier{
		fn: fn,
	}
}

func File(path string) modify.File {
	return &file{path: path}
}

func GoFile(file string) modify.GoFile {
	return &goFile{file: file}
}

func UnregisterMigration(pkg, migration string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return RemoveMigration(pkg, migration)
	})
}

func UnregisterProvider(pkg, provider string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return RemoveProvider(pkg, provider)
	})
}

func UnregisterRoute(pkg, route string) modify.Apply {
	return Call(func(_ []modify.Option) error {
		return RemoveRoute(pkg, route)
	})
}

func When(fn func(options map[string]any) bool, applies ...modify.Apply) modify.Apply {
	return &whenModifier{
		fn:      fn,
		applies: applies,
	}
}

func WhenDriver(driver string, applies ...modify.Apply) modify.Apply {
	return When(func(options map[string]any) bool {
		return driver == options["driver"]
	}, applies...)
}

func WhenFacade(facade string, applies ...modify.Apply) modify.Apply {
	return When(func(options map[string]any) bool {
		return facade == options["facade"]
	}, applies...)
}

func WhenFileContains(file, content string, applies ...modify.Apply) modify.Apply {
	return When(func(_ map[string]any) bool {
		return supportfile.Contains(file, content)
	}, applies...)
}

func WhenFileExists(file string, applies ...modify.Apply) modify.Apply {
	return When(func(_ map[string]any) bool {
		return supportfile.Exists(file)
	}, applies...)
}

func WhenFileNotContains(file, content string, applies ...modify.Apply) modify.Apply {
	return When(func(_ map[string]any) bool {
		return !supportfile.Contains(file, content)
	}, applies...)
}

func WhenFileNotExists(file string, applies ...modify.Apply) modify.Apply {
	return When(func(_ map[string]any) bool {
		return !supportfile.Exists(file)
	}, applies...)
}

func WhenNoFacades(facades []string, applies ...modify.Apply) modify.Apply {
	return When(func(options map[string]any) bool {
		var exist bool
		for _, facade := range facades {
			if facade == options["facade"] {
				continue
			}

			if supportfile.Exists(facadeToFilepath(facade)) {
				exist = true
				break
			}
		}

		return !exist
	}, applies...)
}

func generateOptions(options []modify.Option) map[string]any {
	result := make(map[string]any)
	for _, option := range options {
		option(result)
	}
	return result
}

type callModifier struct {
	fn func(options []modify.Option) error
}

func (r *callModifier) Apply(options ...modify.Option) error {
	return r.fn(options)
}

type file struct {
	path string
}

func (r *file) Append(content string) modify.Apply {
	return &appendFile{
		content: content,
		path:    r.path,
	}
}

func (r *file) Overwrite(content string) modify.Apply {
	return &overwriteFile{
		content: content,
		path:    r.path,
	}
}

func (r *file) Remove() modify.Apply {
	return &removeFile{
		path: r.path,
	}
}

type appendFile struct {
	content string
	path    string
}

func (r *appendFile) Apply(options ...modify.Option) error {
	return supportfile.PutContent(r.path, r.content, supportfile.WithAppend())
}

type overwriteFile struct {
	content string
	path    string
}

func (r *overwriteFile) Apply(options ...modify.Option) error {
	generatedOptions := generateOptions(options)

	if supportfile.Exists(r.path) && !cast.ToBool(generatedOptions["force"]) {
		color.Warningln(errors.ConsoleFileAlreadyExists.Args(r.path))
		return nil
	}

	return supportfile.PutContent(r.path, r.content)
}

type removeFile struct {
	path string
}

func (r *removeFile) Apply(...modify.Option) error {
	return supportfile.Remove(r.path)
}

type goFile struct {
	file      string
	modifiers []modify.GoNode
}

func (r goFile) Apply(...modify.Option) error {
	source, err := supportfile.GetContent(r.file)
	if err != nil {
		return err
	}

	df, err := decorator.Parse(source)
	if err != nil {
		return err
	}

	for i := range r.modifiers {
		if err = r.modifiers[i].Apply(df); err != nil {
			return errors.PackageModifyGoFileFail.Args(r.file, err)
		}
	}

	var buf bytes.Buffer
	err = decorator.Fprint(&buf, df)
	if err != nil {
		return err
	}

	return supportfile.PutContent(r.file, buf.String())
}

func (r goFile) Find(matchers []match.GoNode) modify.GoNode {
	modifier := &goNode{
		matchers: matchers,
		goFile:   &r,
	}
	r.modifiers = append(r.modifiers, modifier)

	return modifier
}

func (r goFile) FindOrCreate(matchers []match.GoNode, fn func(node dst.Node) error) modify.GoNode {
	modifier := &goNode{
		createFunc: fn,
		matchers:   matchers,
		goFile:     &r,
	}
	r.modifiers = append(r.modifiers, modifier)

	return modifier
}

func (r goFile) Format() modify.Apply {
	return Call(func(options []modify.Option) error {
		if err := r.Apply(options...); err != nil {
			return err
		}

		content, err := supportfile.GetContent(r.file)
		if err != nil {
			return err
		}

		// Clear unused imports
		opts := &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  false,
		}
		formatted, err := imports.Process(filepath.Base(r.file), []byte(content), opts)
		if err != nil {
			return err
		}

		// Format source code
		formatted, err = format.Source(formatted)
		if err != nil {
			return err
		}

		return supportfile.PutContent(r.file, string(formatted))
	})
}

type goNode struct {
	actions    []modify.Action
	createFunc func(node dst.Node) error
	goFile     *goFile
	matchers   []match.GoNode
}

func (r *goNode) Apply(node dst.Node) (err error) {
	var (
		current      int
		matched      bool
		matchedNodes = make(map[dst.Node]bool)
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// match the node and apply the action
	dstutil.Apply(node, func(cursor *dstutil.Cursor) bool {
		// if already modified, skip the rest of the nodes
		if matched {
			return false
		}

		if r.matchers[current].MatchCursor(cursor) {
			matchedNodes[cursor.Node()] = true
			current++
			if current == len(r.matchers) {
				// apply the actions after all matchers are matched
				for _, action := range r.actions {
					action(cursor)
				}
				matched = true

				return false
			}
		}

		return true
	}, func(cursor *dstutil.Cursor) bool {
		if nd := cursor.Node(); nd != nil && matchedNodes[nd] {
			return false
		}

		return true
	})

	if !matched {
		if r.createFunc != nil {
			if err := r.createFunc(node); err != nil {
				return err
			}

			r.createFunc = nil // prevent infinite recursion

			return r.Apply(node) // try to apply again after creation
		}

		count := len(r.matchers)
		return errors.PackageMatchGoNodeFail.Args(count-current, count)
	}

	return nil
}

func (r *goNode) Modify(actions ...modify.Action) modify.GoFile {
	r.actions = actions

	return r.goFile
}

type whenModifier struct {
	fn      func(options map[string]any) bool
	applies []modify.Apply
}

func (r whenModifier) Apply(options ...modify.Option) error {
	if !r.fn(generateOptions(options)) {
		return nil
	}

	for _, apply := range r.applies {
		if err := apply.Apply(options...); err != nil {
			return err
		}
	}

	return nil
}

func facadeToFilepath(facade string) string {
	return path.Facade(str.Of(facade).Snake().String() + ".go")
}
