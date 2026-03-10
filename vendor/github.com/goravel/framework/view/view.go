package view

import (
	"path/filepath"
	"sync"

	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/file"
)

type View struct {
	shared sync.Map
}

func NewView() *View {
	return &View{}
}

func (r *View) Exists(view string) bool {
	return file.Exists(filepath.Join(support.Config.Paths.Resources, "views", view))
}

func (r *View) Share(key string, value any) {
	r.shared.Store(key, value)
}

func (r *View) Shared(key string, def ...any) any {
	value, ok := r.shared.Load(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}

		return nil
	}

	return value
}

func (r *View) GetShared() map[string]any {
	shared := make(map[string]any)
	r.shared.Range(func(key, value any) bool {
		shared[key.(string)] = value
		return true
	})

	return shared
}
