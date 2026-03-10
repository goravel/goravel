package gin

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

var contextResponsePool = sync.Pool{New: func() any {
	return &ContextResponse{}
}}

type ContextResponse struct {
	instance *gin.Context
	origin   contractshttp.ResponseOrigin
}

func NewContextResponse(instance *gin.Context, origin contractshttp.ResponseOrigin) contractshttp.ContextResponse {
	response := contextResponsePool.Get().(*ContextResponse)
	response.instance = instance
	response.origin = origin
	return response
}

func (r *ContextResponse) Cookie(cookie contractshttp.Cookie) contractshttp.ContextResponse {
	if cookie.MaxAge == 0 {
		if !cookie.Expires.IsZero() {
			cookie.MaxAge = int(cookie.Expires.Sub(carbon.Now().StdTime()).Seconds())
		}
	}

	sameSiteOptions := map[string]http.SameSite{
		"strict": http.SameSiteStrictMode,
		"lax":    http.SameSiteLaxMode,
		"none":   http.SameSiteNoneMode,
	}
	if sameSite, ok := sameSiteOptions[cookie.SameSite]; ok {
		r.instance.SetSameSite(sameSite)
	}

	r.instance.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	return r
}

func (r *ContextResponse) Data(code int, contentType string, data []byte) contractshttp.AbortableResponse {
	return &DataResponse{code, contentType, data, r.instance}
}

func (r *ContextResponse) Download(filepath, filename string) contractshttp.Response {
	return &DownloadResponse{filename, filepath, r.instance}
}

func (r *ContextResponse) File(filepath string) contractshttp.Response {
	return &FileResponse{filepath, r.instance}
}

func (r *ContextResponse) Header(key, value string) contractshttp.ContextResponse {
	r.instance.Header(key, value)

	return r
}

func (r *ContextResponse) Json(code int, obj any) contractshttp.AbortableResponse {
	return &JsonResponse{code, obj, r.instance}
}

func (r *ContextResponse) NoContent(code ...int) contractshttp.AbortableResponse {
	if len(code) > 0 {
		return &NoContentResponse{code[0], r.instance}
	}

	return &NoContentResponse{http.StatusNoContent, r.instance}
}

func (r *ContextResponse) Origin() contractshttp.ResponseOrigin {
	return r.origin
}

func (r *ContextResponse) Redirect(code int, location string) contractshttp.AbortableResponse {
	return &RedirectResponse{code, location, r.instance}
}

func (r *ContextResponse) String(code int, format string, values ...any) contractshttp.AbortableResponse {
	return &StringResponse{code, format, r.instance, values}
}

func (r *ContextResponse) Success() contractshttp.ResponseStatus {
	return NewStatus(r.instance, http.StatusOK)
}

func (r *ContextResponse) Status(code int) contractshttp.ResponseStatus {
	return NewStatus(r.instance, code)
}

func (r *ContextResponse) Stream(code int, step func(w contractshttp.StreamWriter) error) contractshttp.Response {
	return &StreamResponse{code, r.instance, step}
}

func (r *ContextResponse) View() contractshttp.ResponseView {
	return NewView(r.instance)
}

func (r *ContextResponse) WithoutCookie(name string) contractshttp.ContextResponse {
	r.instance.SetCookie(name, "", -1, "", "", false, false)

	return r
}

func (r *ContextResponse) Writer() http.ResponseWriter {
	return r.instance.Writer
}

func (r *ContextResponse) Flush() {
	r.instance.Writer.Flush()
}

type Status struct {
	instance *gin.Context
	status   int
}

func NewStatus(instance *gin.Context, code int) *Status {
	return &Status{instance, code}
}

func (r *Status) Data(contentType string, data []byte) contractshttp.AbortableResponse {
	return &DataResponse{r.status, contentType, data, r.instance}
}

func (r *Status) Json(obj any) contractshttp.AbortableResponse {
	return &JsonResponse{r.status, obj, r.instance}
}

func (r *Status) String(format string, values ...any) contractshttp.AbortableResponse {
	return &StringResponse{r.status, format, r.instance, values}
}

func (r *Status) Stream(step func(w contractshttp.StreamWriter) error) contractshttp.Response {
	return &StreamResponse{r.status, r.instance, step}
}

func ResponseMiddleware() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		blw := &BodyWriter{body: bytes.NewBufferString("")}
		switch ctx := ctx.(type) {
		case *Context:
			blw.ResponseWriter = ctx.Instance().Writer
			ctx.Instance().Writer = blw
		}

		ctx.WithValue(responseOriginKey, blw)
		ctx.Request().Next()
	}
}

type BodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func (w *BodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)

	return w.ResponseWriter.WriteString(s)
}

func (w *BodyWriter) Body() *bytes.Buffer {
	return w.body
}

func (w *BodyWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
