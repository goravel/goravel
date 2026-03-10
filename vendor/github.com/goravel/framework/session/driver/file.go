package driver

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/file"
)

type File struct {
	path    string
	minutes int
	mu      sync.RWMutex
}

func NewFile(path string, minutes int) *File {
	return &File{
		path:    path,
		minutes: minutes,
	}
}

func (f *File) Close() error {
	return nil
}

func (f *File) Destroy(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return file.Remove(f.getFilePath(id))
}

func (f *File) Gc(maxLifetime int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	cutoffTime := carbon.Now(carbon.UTC).SubSeconds(maxLifetime)

	if !file.Exists(f.path) {
		return nil
	}

	return filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.ModTime().UTC().Before(cutoffTime.StdTime()) {
			return os.Remove(path)
		}

		return nil
	})
}

func (f *File) Open(string, string) error {
	return nil
}

func (f *File) Read(id string) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	path := f.getFilePath(id)
	if file.Exists(path) {
		modified, err := file.LastModified(path, carbon.UTC)
		if err != nil {
			return "", err
		}
		if modified.After(carbon.Now(carbon.UTC).SubMinutes(f.minutes).StdTime()) {
			data, err := os.ReadFile(path)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}
	}

	return "", nil
}

func (f *File) Write(id string, data string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return file.PutContent(f.getFilePath(id), data)
}

func (f *File) getFilePath(id string) string {
	return filepath.Join(f.path, id)
}
