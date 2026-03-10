package client

import "net/http"

type FakeResponse interface {
	// File creates a mock response using the contents of a file at the specified path.
	File(status int, path string) Response

	// Json creates a mock response with a JSON body and "application/json" content type.
	Json(status int, obj any) Response

	// Make constructs a custom mock response with the specified body, status, and headers.
	Make(status int, body string, header http.Header) Response

	// OK creates a generic 200 OK mock response with an empty body.
	OK() Response

	// Status creates a mock response with the specified status code and an empty body.
	Status(status int) Response

	// String creates a mock response with a raw string body.
	String(status int, body string) Response
}
