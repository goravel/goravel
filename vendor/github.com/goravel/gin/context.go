package gin

import (
	"context"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/contracts/http"
)

const (
	contextKey        = "goravel_contextKey"
	responseOriginKey = "goravel_responseOrigin"
	sessionKey        = "goravel_session"
)

var (
	internalContextKeys = []any{
		contextKey,
		responseOriginKey,
		sessionKey,
	}
)

func Background() http.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	return NewContext(ctx)
}

var contextPool = sync.Pool{New: func() any {
	return &Context{}
}}

type Context struct {
	instance *gin.Context
	request  http.ContextRequest
	response http.ContextResponse
}

func NewContext(c *gin.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.instance = c
	return ctx
}

func (c *Context) Request() http.ContextRequest {
	if c.request == nil {
		request := NewContextRequest(c, LogFacade, ValidationFacade)
		c.request = request
	}

	return c.request
}

func (c *Context) Response() http.ContextResponse {
	if c.response == nil {
		response := NewContextResponse(c.instance, &BodyWriter{ResponseWriter: c.instance.Writer})
		c.response = response
	}

	responseOrigin := c.Value(responseOriginKey)
	if responseOrigin != nil {
		c.response.(*ContextResponse).origin = responseOrigin.(http.ResponseOrigin)
	}

	return c.response
}

func (c *Context) WithValue(key any, value any) {
	values := c.getGoravelContextValues()
	values[key] = value
	c.instance.Set(contextKey, values)
}

func (c *Context) WithContext(ctx context.Context) {
	// Changing the request context to a new context
	c.instance.Request = c.instance.Request.WithContext(ctx)
}

func (c *Context) Context() context.Context {
	ctx := c.instance.Request.Context()
	values := c.getGoravelContextValues()
	for key, value := range values {
		skip := false
		for _, internalContextKey := range internalContextKeys {
			if key == internalContextKey {
				skip = true
			}
		}

		if !skip {
			ctx = context.WithValue(ctx, key, value)
		}
	}

	return ctx
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.instance.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.instance.Done()
}

func (c *Context) Err() error {
	return c.instance.Err()
}

func (c *Context) Value(key any) any {
	if value, exist := c.getGoravelContextValues()[key]; exist {
		return value
	}

	return c.instance.Request.Context().Value(key)
}

func (c *Context) Instance() *gin.Context {
	return c.instance
}

func (c *Context) getGoravelContextValues() map[any]any {
	if val, exist := c.instance.Get(contextKey); exist {
		if goravelCtxVal, ok := val.(map[any]any); ok {
			return goravelCtxVal
		}
	}
	return make(map[any]any)
}
