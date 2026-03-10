package translation

import (
	"context"
	"io/fs"

	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
	contractstranslation "github.com/goravel/framework/contracts/translation"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Lang,
		},
		Dependencies: binding.Bindings[binding.Lang].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.BindWith(binding.Lang, func(app foundation.Application, parameters map[string]any) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleLang)
		}

		logger := app.MakeLog()
		if logger == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleLang)
		}

		locale := config.GetString("app.locale")
		fallback := config.GetString("app.fallback_locale")
		path := config.GetString("app.lang_path")
		if path == "" {
			path = support.Config.Paths.Lang
		}

		var fileLoader contractstranslation.Loader
		if path != "" {
			fileLoader = NewFileLoader([]string{cast.ToString(path)}, app.GetJson())
		}

		var fsLoader contractstranslation.Loader
		if f, ok := config.Get("app.lang_fs").(fs.FS); ok {
			fsLoader = NewFSLoader(f, app.GetJson())
		}

		return NewTranslator(parameters["ctx"].(context.Context), fsLoader, fileLoader, locale, fallback, logger), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {

}
