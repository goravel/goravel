package http

import "time"

// Cookie represents an HTTP cookie as defined by RFC 6265.
type Cookie struct {

	// Expires specifies the maximum age of the cookie.It is considered
	// expired if the current time is after the Expires value.
	Expires time.Time

	// Name is the name of the cookie.
	Name string

	// Value is the value associated with the cookie's name.
	Value string

	// Path specifies the subset of URLs to which this cookie applies.
	Path string

	// Domain specifies the domain for which the cookie is valid.
	Domain string

	// Raw is the unparsed value of the "Set-Cookie" header received from
	// the server.
	Raw string

	// SameSite allows a server to define a cookie attribute, making it
	// impossible for the browser to send this cookie along with cross-site
	// requests.It helps mitigate the risk of cross-origin information leaks.
	SameSite string

	// MaxAge specifies the maximum age of the cookie in seconds.A zero or
	// negative MaxAge means that the cookie is not persistent and will be
	// deleted when the browser is closed.
	MaxAge int

	// Secure indicates whether the cookie should only be sent over secure
	// (HTTPS) connections.
	Secure bool

	// HttpOnly indicates whether the cookie is accessible only through
	// HTTP requests, and not through JavaScript.
	HttpOnly bool
}
