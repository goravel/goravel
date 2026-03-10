package http

import (
	"context"
	"io"
	"net/http"
)

type Request interface {
	// Get sends a GET request to the given URI.
	Get(uri string) (Response, error)
	// Post sends a POST request to the given URI with the provided body.
	Post(uri string, body io.Reader) (Response, error)
	// Put sends a PUT request to the given URI with the provided body.
	Put(uri string, body io.Reader) (Response, error)
	// Delete sends a DELETE request to the given URI with the provided body.
	Delete(uri string, body io.Reader) (Response, error)
	// Patch sends a PATCH request to the given URI with the provided body.
	Patch(uri string, body io.Reader) (Response, error)
	// Head sends a HEAD request to the given URI.
	Head(uri string) (Response, error)
	// Options sends an OPTIONS request to the given URI.
	Options(uri string) (Response, error)

	// FlushHeaders removes all currently configured headers.
	FlushHeaders() Request
	// WithBasicAuth sets the Authorization header using Basic Auth credentials.
	WithBasicAuth(username, password string) Request
	// WithContext sets the context for the request.
	WithContext(ctx context.Context) Request
	// WithCookie adds a single cookie to the request.
	WithCookie(cookie *http.Cookie) Request
	// WithCookies adds multiple cookies to the request.
	WithCookies(cookies []*http.Cookie) Request
	// WithHeader sets a specific header key and value.
	WithHeader(key, value string) Request
	// WithHeaders adds multiple headers to the request.
	WithHeaders(headers map[string]string) Request
	// WithoutHeader removes a specific header from the request.
	WithoutHeader(key string) Request
	// WithToken sets the Authorization header using a Bearer token.
	// An optional second argument can specify the token type (default is "Bearer").
	WithToken(token string, ttype ...string) Request
	// WithoutToken removes the Authorization header.
	WithoutToken() Request
	// WithSession sets the session attributes for the request.
	WithSession(attributes map[string]any) Request
}
