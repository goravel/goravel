package file

import (
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/convert"
)

func ClientOriginalExtension(file string) string {
	return strings.ReplaceAll(filepath.Ext(file), ".", "")
}

// DEPRECATED: Use Contains instead
func Contain(file string, search string) bool {
	return Contains(file, search)
}

func Contains(file string, search string) bool {
	if Exists(file) {
		data, err := GetContent(file)
		if err != nil {
			return false
		}

		// Normalize line endings to handle Windows (CRLF) vs Unix (LF) differences
		data = strings.ReplaceAll(data, "\r\n", "\n")
		search = strings.ReplaceAll(search, "\r\n", "\n")

		return strings.Contains(data, search)
	}

	return false
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer errors.Ignore(in.Close)

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer errors.Ignore(out.Close)

	_, err = io.Copy(out, in)

	return err
}

// Create a file with the given content
// Deprecated: Use PutContent instead
func Create(file string, content string) error {
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer errors.Ignore(f.Close)

	if _, err = f.WriteString(content); err != nil {
		return err
	}

	return nil
}

func Exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// Extension Supported types: https://github.com/gabriel-vasile/mimetype/blob/master/supported_mimes.md
func Extension(file string, originalWhenUnknown ...bool) (string, error) {
	getOriginal := false
	if len(originalWhenUnknown) > 0 {
		getOriginal = originalWhenUnknown[0]
	}

	mtype, err := mimetype.DetectFile(file)
	if err != nil && !getOriginal {
		return "", err
	}

	if mtype != nil && mtype.Extension() != "" {
		return strings.TrimPrefix(mtype.Extension(), "."), nil
	}

	if getOriginal {
		return ClientOriginalExtension(file), nil
	}

	return "", errors.UnknownFileExtension
}

func GetContent(file string) (string, error) {
	// Read the entire file
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return convert.UnsafeString(data), nil
}

func GetFrameworkContent(file string) (string, error) {
	return GetPackageContent("github.com/goravel/framework", file)
}

func GetPackageContent(pkgName, file string) (string, error) {
	pkg, err := build.Import(pkgName, "", build.FindOnly)
	if err != nil {
		return "", err
	}

	paths := strings.Split(file, "/")
	paths = append([]string{pkg.Dir}, paths...)

	return GetContent(filepath.Join(paths...))
}

func LastModified(file, timezone string) (time.Time, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return time.Time{}, err
	}

	l, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime().In(l), nil
}

func MimeType(file string) (string, error) {
	mtype, err := mimetype.DetectFile(file)
	if err != nil {
		return "", err
	}

	return mtype.String(), nil
}

func PutContent(file string, content string, options ...Option) error {
	// Default options
	opts := &option{
		mode:   os.ModePerm,
		append: false,
	}

	// Apply options
	for _, option := range options {
		option(opts)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(file), opts.mode); err != nil {
		return err
	}

	// Open file with appropriate flags
	flag := os.O_CREATE | os.O_WRONLY
	if opts.append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	// Open the file
	f, err := os.OpenFile(file, flag, opts.mode)
	if err != nil {
		return err
	}
	defer errors.Ignore(f.Close)

	// Write the content
	if _, err = f.WriteString(content); err != nil {
		return err
	}

	return nil
}

func Remove(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return os.RemoveAll(file)
}

func Size(file string) (int64, error) {
	fileInfo, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer errors.Ignore(fileInfo.Close)

	fi, err := fileInfo.Stat()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
