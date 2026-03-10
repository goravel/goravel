package http

import (
	"bytes"
	"net/http"
)

type Json map[string]any

type ContextResponse interface {
	// Cookie adds a cookie to the response.
	Cookie(cookie Cookie) ContextResponse
	// Data write the given data to the response.
	Data(code int, contentType string, data []byte) AbortableResponse
	// Download initiates a file download by specifying the file path and the desired filename
	Download(filepath, filename string) Response
	// File serves a file located at the specified file path as the response.
	File(filepath string) Response
	// Header sets an HTTP header field with the given key and value.
	Header(key, value string) ContextResponse
	// Json sends a JSON response with the specified status code and data object.
	Json(code int, obj any) AbortableResponse
	// NoContent sends a response with no-body and the specified status code.
	NoContent(code ...int) AbortableResponse
	// Origin returns the ResponseOrigin
	Origin() ResponseOrigin
	// Redirect performs an HTTP redirect to the specified location with the given status code.
	Redirect(code int, location string) AbortableResponse
	// String writes a string response with the specified status code and format.
	// The 'values' parameter can be used to replace placeholders in the format string.
	String(code int, format string, values ...any) AbortableResponse
	// Success returns ResponseStatus with a 200 status code.
	Success() ResponseStatus
	// Status sets the HTTP response status code and returns the ResponseStatus.
	Status(code int) ResponseStatus
	// Stream sends a streaming response with the specified status code and the given reader.
	Stream(code int, step func(w StreamWriter) error) Response
	// View returns ResponseView
	View() ResponseView
	// Writer returns the underlying http.ResponseWriter associated with the response.
	Writer() http.ResponseWriter
	// WithoutCookie removes a cookie from the response.
	WithoutCookie(name string) ContextResponse
	// Flush flushes any buffered data to the client.
	Flush()
}

type Response interface {
	Render() error
}

type AbortableResponse interface {
	Response
	Abort() error
}

type StreamWriter interface {
	// Write writes the specified data to the response.
	Write(data []byte) (int, error)

	// WriteString writes the specified string to the response.
	WriteString(data string) (int, error)

	// Flush flushes any buffered data to the client.
	Flush() error
}

type ResponseStatus interface {
	// Data write the given data to the Response.
	Data(contentType string, data []byte) AbortableResponse
	// Json sends a JSON AbortResponse with the specified data object.
	Json(obj any) AbortableResponse
	// String writes a string AbortResponse with the specified format and values.
	String(format string, values ...any) AbortableResponse
	// Stream sends a streaming response with the specified status code and the given reader.
	Stream(step func(w StreamWriter) error) Response
}

type ResponseOrigin interface {
	// Body returns the response's body content as a *bytes.Buffer.
	Body() *bytes.Buffer
	// Header returns the response's HTTP header.
	Header() http.Header
	// Size returns the size, in bytes, of the response's body content.
	Size() int
	// Status returns the HTTP status code of the response.
	Status() int
}

type ResponseView interface {
	// Make generates a Response for the specified view with optional data.
	Make(view string, data ...any) Response
	// First generates a response for the first available view from the provided list.
	First(views []string, data ...any) Response
}
