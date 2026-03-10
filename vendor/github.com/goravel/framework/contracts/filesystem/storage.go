package filesystem

import (
	"context"
	"time"
)

type Storage interface {
	Driver
	// Disk gets the instance of the given disk.
	Disk(disk string) Driver
}

type Driver interface {
	// AllDirectories gets all the directories within a given directory(recursive).
	AllDirectories(path string) ([]string, error)
	// AllFiles gets all the files from the given directory(recursive).
	AllFiles(path string) ([]string, error)
	// Copy the given file to a new location.
	Copy(oldFile, newFile string) error
	// Delete deletes the given file(s).
	Delete(file ...string) error
	// DeleteDirectory deletes the given directory(recursive).
	DeleteDirectory(directory string) error
	// Directories get all the directories within a given directory.
	Directories(path string) ([]string, error)
	// Exists determines if a file exists.
	Exists(file string) bool
	// Files gets all the files from the given directory.
	Files(path string) ([]string, error)
	// Get gets the contents of a file.
	Get(file string) (string, error)
	// GetBytes gets the contents of a file as a byte array.
	GetBytes(file string) ([]byte, error)
	// LastModified gets the file's last modified time.
	LastModified(file string) (time.Time, error)
	// MakeDirectory creates a directory.
	MakeDirectory(directory string) error
	// MimeType gets the file's mime type.
	MimeType(file string) (string, error)
	// Missing determines if a file is missing.
	Missing(file string) bool
	// Move a file to a new location.
	Move(oldFile, newFile string) error
	// Path gets the full path for the file.
	Path(file string) string
	// Put writes the contents of a file.
	Put(file, content string) error
	// PutFile upload the given file.
	PutFile(path string, source File) (string, error)
	// PutFileAs upload the given file with a new name.
	PutFileAs(path string, source File, name string) (string, error)
	// Size gets the file size of a given file.
	Size(file string) (int64, error)
	// TemporaryUrl get a temporary URL for the file.
	TemporaryUrl(file string, time time.Time) (string, error)
	// WithContext sets the context to be used by the driver.
	WithContext(ctx context.Context) Driver
	// Url get the URL for the file at the given path.
	Url(file string) string
}

type File interface {
	// Disk gets the instance of the given disk.
	Disk(disk string) File
	// Extension gets the file extension.
	Extension() (string, error)
	// File gets the file path.
	File() string
	// GetClientOriginalName gets the client original name.
	GetClientOriginalName() string
	// GetClientOriginalExtension gets the client original extension.
	GetClientOriginalExtension() string
	// HashName gets the file's hash name.
	HashName(path ...string) string
	// LastModified gets the file's last modified time.
	LastModified() (time.Time, error)
	// MimeType gets the file's mime type.
	MimeType() (string, error)
	// Size gets the file size.
	Size() (int64, error)
	// Store the file at the given path.
	Store(path string) (string, error)
	// StoreAs store the file at the given path with a new name.
	StoreAs(path string, name string) (string, error)
}
