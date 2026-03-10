package filesystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/errors"
	supportfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type Local struct {
	config config.Config
	root   string
	url    string
}

func NewLocal(config config.Config, disk string) (*Local, error) {
	return &Local{
		config: config,
		root:   config.GetString(fmt.Sprintf("filesystems.disks.%s.root", disk)),
		url:    config.GetString(fmt.Sprintf("filesystems.disks.%s.url", disk)),
	}, nil
}

func (r *Local) AllDirectories(path string) ([]string, error) {
	var directories []string
	err := filepath.Walk(r.fullPath(path), func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			realPath := strings.ReplaceAll(fullPath, r.fullPath(path), "")
			realPath = strings.TrimPrefix(realPath, string(filepath.Separator))
			if realPath != "" {
				directories = append(directories, realPath+string(filepath.Separator))
			}
		}

		return nil
	})

	return directories, err
}

func (r *Local) AllFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(r.fullPath(path), func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, strings.ReplaceAll(fullPath, r.fullPath(path)+string(filepath.Separator), ""))
		}

		return nil
	})

	return files, err
}

func (r *Local) Copy(originFile, targetFile string) error {
	content, err := r.Get(originFile)
	if err != nil {
		return err
	}

	return r.Put(targetFile, content)
}

func (r *Local) Delete(files ...string) error {
	for _, file := range files {
		fileInfo, err := os.Stat(r.fullPath(file))
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.FilesystemDeleteDirectory
		}
	}

	for _, file := range files {
		if err := os.Remove(r.fullPath(file)); err != nil {
			return err
		}
	}

	return nil
}

func (r *Local) DeleteDirectory(directory string) error {
	return os.RemoveAll(r.fullPath(directory))
}

func (r *Local) Directories(path string) ([]string, error) {
	var directories []string
	fileInfo, _ := os.ReadDir(r.fullPath(path))
	for _, f := range fileInfo {
		if f.IsDir() {
			directories = append(directories, f.Name()+string(filepath.Separator))
		}
	}

	return directories, nil
}

func (r *Local) Exists(file string) bool {
	_, err := os.Stat(r.fullPath(file))
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func (r *Local) Files(path string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(r.fullPath(path))
	if err != nil {
		return nil, err
	}
	for _, f := range fileInfo {
		if !f.IsDir() {
			files = append(files, f.Name())
		}
	}

	return files, nil
}

func (r *Local) Get(file string) (string, error) {
	data, err := r.GetBytes(file)

	return string(data), err
}

func (r *Local) GetBytes(file string) ([]byte, error) {
	data, err := os.ReadFile(r.fullPath(file))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Local) LastModified(file string) (time.Time, error) {
	return supportfile.LastModified(r.fullPath(file), r.config.GetString("app.timezone"))
}

func (r *Local) MakeDirectory(directory string) error {
	return os.MkdirAll(filepath.Dir(r.fullPath(directory)+string(filepath.Separator)), os.ModePerm)
}

func (r *Local) MimeType(file string) (string, error) {
	return supportfile.MimeType(r.fullPath(file))
}

func (r *Local) Missing(file string) bool {
	return !r.Exists(file)
}

func (r *Local) Move(oldFile, newFile string) error {
	newFile = r.fullPath(newFile)
	if err := os.MkdirAll(filepath.Dir(newFile), os.ModePerm); err != nil {
		return err
	}

	if err := os.Rename(r.fullPath(oldFile), newFile); err != nil {
		return err
	}

	return nil
}

func (r *Local) Path(file string) string {
	return r.fullPath(file)
}

func (r *Local) Put(file, content string) error {
	file = r.fullPath(file)
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

func (r *Local) PutFile(filePath string, source filesystem.File) (string, error) {
	return r.PutFileAs(filePath, source, str.Random(40))
}

func (r *Local) PutFileAs(filePath string, source filesystem.File, name string) (string, error) {
	data, err := os.ReadFile(source.File())
	if err != nil {
		return "", err
	}

	fullPath, err := fullPathOfFile(filePath, source, name)
	if err != nil {
		return "", err
	}

	if err := r.Put(fullPath, string(data)); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (r *Local) Size(file string) (int64, error) {
	return supportfile.Size(r.fullPath(file))
}

func (r *Local) TemporaryUrl(file string, time time.Time) (string, error) {
	return r.Url(file), nil
}

func (r *Local) WithContext(ctx context.Context) filesystem.Driver {
	return r
}

func (r *Local) Url(file string) string {
	return strings.TrimSuffix(r.url, "/") + "/" + strings.TrimPrefix(filepath.ToSlash(file), "/")
}

func (r *Local) fullPath(path string) string {
	realPath := filepath.Clean(path)

	if realPath == "." {
		return r.rootPath()
	}

	return filepath.Join(r.rootPath(), realPath)
}

func (r *Local) rootPath() string {
	return strings.TrimSuffix(r.root, string(filepath.Separator)) + string(filepath.Separator)
}
