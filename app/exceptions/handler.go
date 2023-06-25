package exceptions

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation/exception"
)

type Handler struct {
	exception.Handler
}

func (h *Handler) Register(app foundation.Application) {
	app.Singleton(exception.Binding, func(app foundation.Application) (any, error) {
		return NewHandler(app.MakeConfig()), nil
	})
}

func (h *Handler) Report(throwable error) {
	// check your throwable here

	h.Handler.Report(throwable)
}

func NewHandler(config config.Config) *Handler {
	return &Handler{
		Handler: exception.Handler{
			Config:     config,
			DontReport: []error{},
		},
	}
}
