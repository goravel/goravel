package translation

import (
	"io/fs"
	"path"

	"github.com/goravel/framework/contracts/foundation"
	contractstranslation "github.com/goravel/framework/contracts/translation"
	"github.com/goravel/framework/errors"
)

type FSLoader struct {
	fs   fs.FS
	json foundation.Json
}

func NewFSLoader(fs fs.FS, json foundation.Json) contractstranslation.Loader {
	return &FSLoader{
		fs:   fs,
		json: json,
	}
}

func (f *FSLoader) Load(locale string, group string) (map[string]any, error) {
	var val map[string]any
	fullPath := path.Join(locale, group+".json")
	if group == "*" {
		fullPath = locale + ".json"
	}

	data, err := fs.ReadFile(f.fs, fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.LangFileNotExist
		}
		return nil, err
	}
	if err = f.json.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return val, nil
}
