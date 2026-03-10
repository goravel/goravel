package http

import (
	"context"

	"github.com/goravel/framework/contracts/http"
)

func Background() http.Context {
	return NewContext()
}

type Ctx context.Context

type Context struct {
	Ctx
}

func NewContext() *Context {
	return &Context{
		Ctx: context.Background(),
	}
}

func (c *Context) Context() context.Context {
	return c.Ctx
}

func (c *Context) WithContext(ctx context.Context) {
	// Changing the request context to a new context
	c.Ctx = ctx
}

func (c *Context) WithValue(key any, value any) {
	// nolint:all
	c.Ctx = context.WithValue(c.Ctx, key, value)
}

func (c *Context) Request() http.ContextRequest {
	return nil
}

func (c *Context) Response() http.ContextResponse {
	return nil
}
