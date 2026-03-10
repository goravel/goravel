package client

import (
	"context"
	"io"
	"net/http"
)

type Request interface {
	// Accept sets the "Accept" header to the specified content type.
	Accept(contentType string) Request

	// AcceptJSON sets the "Accept" header to "application/json".
	AcceptJSON() Request

	// AsForm sets the "Content-Type" header to "application/x-www-form-urlencoded".
	AsForm() Request

	// BaseUrl sets the base URL for the request, overriding the configuration.
	BaseUrl(url string) Request

	// Body returns the raw payload of the request as a string.
	Body() string

	// ClientName returns the name of the client configuration used for this request.
	ClientName() string

	// Clone creates a deep copy of the request builder.
	Clone() Request

	// Delete sends a DELETE request to the specified URI with the given body.
	Delete(uri string, body io.Reader) (Response, error)

	// FlushHeaders removes all currently configured headers from the request builder.
	FlushHeaders() Request

	// Get sends a GET request to the specified URI.
	Get(uri string) (Response, error)

	// Head sends a HEAD request to the specified URI.
	Head(uri string) (Response, error)

	// Header retrieves the value of a specific header key.
	Header(key string) string

	// Headers retrieves all headers associated with the request.
	Headers() http.Header

	// HttpClient returns the underlying standard library *http.Client.
	HttpClient() *http.Client

	// Input retrieves a specific value from the request body or query parameters.
	Input(key string) any

	// Method returns the HTTP verb of the request.
	Method() string

	// Options sends an OPTIONS request to the specified URI.
	Options(uri string) (Response, error)

	// Patch sends a PATCH request to the specified URI with the given body.
	Patch(uri string, body io.Reader) (Response, error)

	// Post sends a POST request to the specified URI with the given body.
	Post(uri string, body io.Reader) (Response, error)

	// Put sends a PUT request to the specified URI with the given body.
	Put(uri string, body io.Reader) (Response, error)

	// ReplaceHeaders replaces all existing headers with the provided map.
	ReplaceHeaders(headers map[string]string) Request

	// Url returns the full, resolved URL of the request.
	Url() string

	// WithBasicAuth sets the "Authorization" header using the Basic Auth standard.
	WithBasicAuth(username, password string) Request

	// WithContext sets the context for the request.
	WithContext(ctx context.Context) Request

	// WithCookie adds a single http.Cookie to the request.
	WithCookie(cookie *http.Cookie) Request

	// WithCookies adds multiple http.Cookie objects to the request.
	WithCookies(cookies []*http.Cookie) Request

	// WithHeader adds a specific header key-value pair to the request.
	WithHeader(key, value string) Request

	// WithHeaders adds multiple headers to the request from a map.
	WithHeaders(headers map[string]string) Request

	// WithQueryParameter adds a single query parameter to the URL.
	WithQueryParameter(key, value string) Request

	// WithQueryParameters adds multiple query parameters to the URL from a map.
	WithQueryParameters(params map[string]string) Request

	// WithQueryString parses a raw query string and adds it to the URL.
	WithQueryString(query string) Request

	// WithToken sets the "Authorization" header using a Bearer token.
	WithToken(token string, ttype ...string) Request

	// WithUrlParameter replaces a URL parameter placeholder with the given value.
	WithUrlParameter(key, value string) Request

	// WithUrlParameters replaces multiple URL parameter placeholders with values from a map.
	WithUrlParameters(params map[string]string) Request

	// WithoutHeader removes a specific header from the request by key.
	WithoutHeader(key string) Request

	// WithoutToken removes the "Authorization" header from the request.
	WithoutToken() Request
}
