package filesystem

import (
	"fmt"
	"io"
	"mime/multipart"
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

type File struct {
	config  config.Config
	storage filesystem.Storage
	disk    string
	path    string
	name    string
}

func NewFile(file string) (*File, error) {
	if ConfigFacade == nil {
		return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleFilesystem)
	}

	if !supportfile.Exists(file) {
		return nil, errors.FilesystemFileNotExist
	}

	return &File{
		config:  ConfigFacade,
		disk:    ConfigFacade.GetString("filesystems.default"),
		path:    file,
		name:    filepath.Base(file),
		storage: StorageFacade,
	}, nil
}

func NewFileFromRequest(fileHeader *multipart.FileHeader) (*File, error) {
	if ConfigFacade == nil {
		return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleFilesystem)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer func(src multipart.File) {
		if err = src.Close(); err != nil {
			panic(err)
		}
	}(src)

	tempFileName := fmt.Sprintf("%s_*%s", ConfigFacade.GetString("app.name"), filepath.Ext(fileHeader.Filename))
	tempFile, err := os.CreateTemp(os.TempDir(), tempFileName)
	if err != nil {
		return nil, err
	}
	defer func(tempFile *os.File) {
		if err = tempFile.Close(); err != nil {
			panic(err)
		}
	}(tempFile)

	_, err = io.Copy(tempFile, src)
	if err != nil {
		return nil, err
	}

	return &File{
		config:  ConfigFacade,
		disk:    ConfigFacade.GetString("filesystems.default"),
		path:    tempFile.Name(),
		name:    fileHeader.Filename,
		storage: StorageFacade,
	}, nil
}

func (f *File) Disk(disk string) filesystem.File {
	f.disk = disk

	return f
}

func (f *File) Extension() (string, error) {
	return supportfile.Extension(f.path)
}

func (f *File) File() string {
	return f.path
}

func (f *File) GetClientOriginalName() string {
	return f.name
}

func (f *File) GetClientOriginalExtension() string {
	return supportfile.ClientOriginalExtension(f.name)
}

func (f *File) HashName(path ...string) string {
	var realPath string
	if len(path) > 0 {
		realPath = strings.TrimRight(path[0], "/") + "/"
	}

	extension, _ := supportfile.Extension(f.path, true)
	if extension == "" {
		return realPath + str.Random(40)
	}

	return realPath + str.Random(40) + "." + extension
}

func (f *File) LastModified() (time.Time, error) {
	return supportfile.LastModified(f.path, f.config.GetString("app.timezone"))
}

func (f *File) MimeType() (string, error) {
	return supportfile.MimeType(f.path)
}

func (f *File) Size() (int64, error) {
	return supportfile.Size(f.path)
}

func (f *File) Store(path string) (string, error) {
	if err := f.validateStorageFacade(); err != nil {
		return "", err
	}

	return f.storage.Disk(f.disk).PutFile(path, f)
}

func (f *File) StoreAs(path string, name string) (string, error) {
	if err := f.validateStorageFacade(); err != nil {
		return "", err
	}

	return f.storage.Disk(f.disk).PutFileAs(path, f, name)
}

func (f *File) validateStorageFacade() error {
	if f.storage == nil {
		return errors.StorageFacadeNotSet.SetModule(errors.ModuleFilesystem)
	}

	return nil
}
