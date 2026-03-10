package client

import (
	"io"
	"net/http"
)

type Response interface {
	// Accepted determines if the response status code is 202 Accepted.
	Accepted() bool

	// BadRequest determines if the response status code is 400 Bad Request.
	BadRequest() bool

	// Bind unmarshalls the JSON response body into the provided value.
	Bind(value any) error

	// Body returns the response body as a string.
	Body() (string, error)

	// ClientError determines if the response status code is in the 400-499 range.
	ClientError() bool

	// Conflict determines if the response status code is 409 Conflict.
	Conflict() bool

	// Cookie retrieves a cookie by name from the response.
	Cookie(name string) *http.Cookie

	// Cookies returns all cookies provided by the response.
	Cookies() []*http.Cookie

	// Created determines if the response status code is 201 Created.
	Created() bool

	// Failed determines if the response status code is >= 400.
	Failed() bool

	// Forbidden determines if the response status code is 403 Forbidden.
	Forbidden() bool

	// Found determines if the response status code is 302 Found.
	Found() bool

	// Header retrieves the first value of a specific header from the response.
	Header(name string) string

	// Headers returns all headers from the response.
	Headers() http.Header

	// Json returns the response body parsed as a map.
	Json() (map[string]any, error)

	// MovedPermanently determines if the response status code is 301 Moved Permanently.
	MovedPermanently() bool

	// NoContent determines if the response status code is 204 No Content.
	NoContent() bool

	// NotFound determines if the response status code is 404 Not Found.
	NotFound() bool

	// OK determines if the response status code is 200 OK.
	OK() bool

	// Origin returns the underlying standard library *http.Response.
	Origin() *http.Response

	// PaymentRequired determines if the response status code is 402 Payment Required.
	PaymentRequired() bool

	// Redirect determines if the response status code is in the 300-399 range.
	Redirect() bool

	// RequestTimeout determines if the response status code is 408 Request Timeout.
	RequestTimeout() bool

	// ServerError determines if the response status code is >= 500.
	ServerError() bool

	// Status returns the integer HTTP status code of the response.
	Status() int

	// Stream returns the underlying reader to stream the response body.
	Stream() (io.ReadCloser, error)

	// Successful determines if the response status code is in the 200-299 range.
	Successful() bool

	// TooManyRequests determines if the response status code is 429 Too Many Requests.
	TooManyRequests() bool

	// Unauthorized determines if the response status code is 401 Unauthorized.
	Unauthorized() bool

	// UnprocessableEntity determines if the response status code is 422 Unprocessable Entity.
	UnprocessableEntity() bool
}
