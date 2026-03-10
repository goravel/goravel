package client

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/convert"
)

var _ client.Response = (*Response)(nil)

type Response struct {
	json     foundation.Json
	decoded  map[string]any
	response *http.Response
	content  string
	mu       sync.Mutex
	// streamed tracks if the raw body has been handed over to the user via Stream().
	// If true, subsequent calls to Bind/Body/Json must fail to avoid data corruption.
	streamed bool
}

func NewResponse(response *http.Response, json foundation.Json) *Response {
	return &Response{
		json:     json,
		response: response,
	}
}

func (r *Response) Bind(value any) error {
	// Do NOT lock here. getContent() handles its own locking.
	// Locking here would cause a Deadlock because getContent() re-acquires the mutex.
	content, err := r.getContent()
	if err != nil {
		return err
	}

	if err = r.json.UnmarshalString(content, value); err != nil {
		return err
	}

	return nil
}

func (r *Response) Body() (string, error) {
	return r.getContent()
}

func (r *Response) ClientError() bool {
	return r.getStatusCode() >= http.StatusBadRequest && r.getStatusCode() < http.StatusInternalServerError
}

func (r *Response) Cookie(name string) *http.Cookie {
	return r.getCookie(name)
}

func (r *Response) Cookies() []*http.Cookie {
	if r.response != nil {
		return r.response.Cookies()
	}
	return []*http.Cookie{}
}

func (r *Response) Failed() bool {
	return r.ServerError() || r.ClientError()
}

func (r *Response) Header(name string) string {
	return r.getHeader(name)
}

func (r *Response) Headers() http.Header {
	if r.response != nil {
		return r.response.Header
	}
	return http.Header{}
}

func (r *Response) Json() (map[string]any, error) {
	if r.decoded != nil {
		return r.decoded, nil
	}

	content, err := r.getContent()
	if err != nil {
		return nil, err
	}

	if err := r.json.UnmarshalString(content, &r.decoded); err != nil {
		return nil, err
	}

	return r.decoded, nil
}

func (r *Response) Origin() *http.Response {
	return r.response
}

func (r *Response) Redirect() bool {
	status := r.getStatusCode()
	return status >= http.StatusMultipleChoices && status < http.StatusBadRequest
}

func (r *Response) ServerError() bool {
	return r.getStatusCode() >= http.StatusInternalServerError
}

func (r *Response) Status() int {
	return r.getStatusCode()
}

func (r *Response) Stream() (io.ReadCloser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If the user already called Bind(), Body(), or Json(), the content is
	// stored in memory (r.content). We return a reader for this cached string
	// so the stream works seamlessly even after parsing.
	if r.content != "" {
		return io.NopCloser(strings.NewReader(r.content)), nil
	}

	// If the stream was already given away, we cannot read it again or cache it.
	if r.streamed {
		return nil, errors.HttpClientResponseAlreadyStreamed
	}

	if r.response == nil {
		return nil, errors.HttpClientResponseIsNil
	}

	// Valid State: Body is nil (e.g., 204 No Content or Mock)
	// Return an empty reader so the user code doesn't crash on io.ReadAll(stream)
	if r.response.Body == nil {
		return io.NopCloser(strings.NewReader("")), nil
	}

	// Mark as streamed so getContent() knows the body is gone.
	r.streamed = true
	return r.response.Body, nil
}

func (r *Response) Successful() bool {
	status := r.getStatusCode()
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

func (r *Response) OK() bool               { return r.getStatusCode() == http.StatusOK }
func (r *Response) Created() bool          { return r.getStatusCode() == http.StatusCreated }
func (r *Response) Accepted() bool         { return r.getStatusCode() == http.StatusAccepted }
func (r *Response) NoContent() bool        { return r.getStatusCode() == http.StatusNoContent }
func (r *Response) MovedPermanently() bool { return r.getStatusCode() == http.StatusMovedPermanently }
func (r *Response) Found() bool            { return r.getStatusCode() == http.StatusFound }
func (r *Response) BadRequest() bool       { return r.getStatusCode() == http.StatusBadRequest }
func (r *Response) Unauthorized() bool     { return r.getStatusCode() == http.StatusUnauthorized }
func (r *Response) PaymentRequired() bool  { return r.getStatusCode() == http.StatusPaymentRequired }
func (r *Response) Forbidden() bool        { return r.getStatusCode() == http.StatusForbidden }
func (r *Response) NotFound() bool         { return r.getStatusCode() == http.StatusNotFound }
func (r *Response) RequestTimeout() bool   { return r.getStatusCode() == http.StatusRequestTimeout }
func (r *Response) Conflict() bool         { return r.getStatusCode() == http.StatusConflict }
func (r *Response) UnprocessableEntity() bool {
	return r.getStatusCode() == http.StatusUnprocessableEntity
}
func (r *Response) TooManyRequests() bool { return r.getStatusCode() == http.StatusTooManyRequests }

func (r *Response) getStatusCode() int {
	if r.response != nil {
		return r.response.StatusCode
	}
	return 0
}

func (r *Response) getContent() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.content != "" {
		return r.content, nil
	}

	if r.streamed {
		return "", errors.HttpClientResponseAlreadyStreamed
	}

	if r.response == nil {
		return "", errors.HttpClientResponseIsNil
	}

	if r.response.Body == nil {
		return "", nil
	}

	defer errors.Ignore(r.response.Body.Close)

	content, err := io.ReadAll(r.response.Body)
	if err != nil {
		return "", err
	}

	r.content = convert.UnsafeString(content)
	return r.content, nil
}

func (r *Response) getCookie(name string) *http.Cookie {
	if r.response == nil {
		return nil
	}
	for _, c := range r.response.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (r *Response) getHeader(name string) string {
	if r.response != nil {
		return r.response.Header.Get(name)
	}
	return ""
}
