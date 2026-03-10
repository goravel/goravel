package console

type Stubs struct {
}

func (r Stubs) Request() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DummyRequest struct {
	DummyField
}

func (r *DummyRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *DummyRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
`
}

func (r Stubs) Controller() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/http"
)

type DummyController struct {
	// Dependent services
}

func NewDummyController() *DummyController {
	return &DummyController{
		// Inject services
	}
}

func (r *DummyController) Index(ctx http.Context) http.Response {
	return nil
}	
`
}

func (r Stubs) ResourceController() string {
	return r.Controller() + `
func (r *DummyController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *DummyController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *DummyController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *DummyController) Destroy(ctx http.Context) http.Response {
	return nil
}
`
}

func (r Stubs) Middleware() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/http"
)

func DummyMiddleware() http.Middleware {
	return func(ctx http.Context) {
		ctx.Request().Next()
	}
}
`
}
