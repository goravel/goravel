package console

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/pterm/pterm"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type Make struct {
	name string
	root string
}

func NewMake(ctx console.Context, ttype, name, root string) (*Make, error) {
	if name == "" {
		var err error
		name, err = ctx.Ask(fmt.Sprintf("Enter the %s name", ttype), console.AskOption{
			Validate: func(s string) error {
				if s == "" {
					return errors.ConsoleEmptyFieldValue.Args(ttype)
				}

				return nil
			},
		})
		if err != nil {
			return nil, err
		}
	}

	m := &Make{
		name: name,
		root: root,
	}

	if !ctx.OptionBool("force") && file.Exists(m.GetFilePath()) {
		return nil, errors.ConsoleFileAlreadyExists.Args(ttype)
	}

	return m, nil
}

func (m *Make) GetName() string {
	return m.name
}

func (m *Make) GetFilePath() string {
	root := strings.Split(m.root, "/")
	paths := append(root, m.GetFolderPath(), str.Of(m.GetStructName()).Snake().String()+".go")
	path := filepath.Join(paths...)
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func (m *Make) GetSignature() string {
	return str.Of(filepath.Join(m.GetFolderPath(), m.GetStructName())).
		Replace(string(filepath.Separator), "_").Studly().String()
}

func (m *Make) GetStructName() string {
	name := strings.TrimSuffix(m.name, ".go")
	segments := strings.Split(name, "/")

	return str.Of(segments[len(segments)-1]).Studly().String()
}

func (m *Make) GetPackageImportPath() string {
	var paths []string
	if info, ok := debug.ReadBuildInfo(); ok {
		paths = append(paths, info.Main.Path)
	}

	if len(m.root) > 0 {
		paths = append(paths, strings.Split(m.root, string(filepath.Separator))...)
	}

	if folders := m.GetFolderPath(); len(folders) > 0 {
		paths = append(paths, strings.Split(folders, string(filepath.Separator))...)
	}

	return strings.Join(paths, "/")
}

func (m *Make) GetPackageName() string {
	name := strings.TrimSuffix(m.name, ".go")
	segments := strings.Split(name, "/")
	// Normalize path separators to forward slashes for cross-platform compatibility
	// Replace backslashes with forward slashes to handle Windows paths
	normalizedRoot := strings.ReplaceAll(m.root, "\\", "/")
	packageName := str.Of(normalizedRoot).Trim("/").AfterLast("/").String()

	if len(segments) > 1 {
		packageName = segments[len(segments)-2]
	}

	return packageName
}

func (m *Make) GetModuleName() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Path
	}

	// fallback to default
	return "goravel"
}

func (m *Make) GetFolderPath() string {
	name := strings.TrimSuffix(m.name, ".go")
	segments := strings.Split(name, "/")

	var folderPath string
	if len(segments) > 1 {
		folderPath = filepath.Join(segments[:len(segments)-1]...)
	}

	return folderPath
}

func ConfirmToProceed(ctx console.Context, env string) bool {
	if env != "production" {
		return true
	}
	if ctx.OptionBool("force") {
		return true
	}

	return ctx.Confirm("Are you sure you want to run this command?")
}

func TwoColumnDetail(first, second string, filler ...rune) string {
	margin := func(s string, left, right int) string {
		var builder strings.Builder
		if left > 0 {
			builder.WriteString(strings.Repeat(" ", left))
		}
		builder.WriteString(s)
		if right > 0 {
			builder.WriteString(strings.Repeat(" ", right))
		}
		return builder.String()
	}
	width := func(s string) int {
		return runewidth.StringWidth(pterm.RemoveColorFromString(s))
	}
	first = margin(first, 2, 1)
	if w := width(second); w > 0 {
		second = margin(second, 1, 2)
	} else {
		second = margin(second, 0, 2)
	}
	fillingText := ""
	if w := pterm.GetTerminalWidth() - width(first) - width(second); w > 0 {
		fillingText = color.Gray().Sprint(strings.Repeat(string(append(filler, '.')[0]), w))
	}

	return first + fillingText + second
}
