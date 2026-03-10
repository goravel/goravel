package http

import "net/http"

type Response interface {
	// Bind unmarshalls the JSON response body into the provided value.
	Bind(value any) error
	// Content returns the raw response body as a string.
	Content() (string, error)
	// Cookie retrieves a cookie from the response by name.
	Cookie(name string) *http.Cookie
	// Cookies returns all cookies from the response.
	Cookies() []*http.Cookie
	// Headers returns the response headers.
	Headers() http.Header
	// IsServerError returns true if the status code is >= 500 and < 600.
	IsServerError() bool
	// IsSuccessful returns true if the status code is >= 200 and < 300.
	IsSuccessful() bool
	// Json returns the response body as a map.
	Json() (map[string]any, error)
	// Session returns the session attributes associated with the response.
	Session() (map[string]any, error)

	// AssertStatus asserts that the response has the given status code.
	AssertStatus(status int) Response
	// AssertOk asserts that the response has a 200 OK status code.
	AssertOk() Response
	// AssertCreated asserts that the response has a 201 Created status code.
	AssertCreated() Response
	// AssertAccepted asserts that the response has a 202 Accepted status code.
	AssertAccepted() Response
	// AssertNoContent asserts that the response has the given status code (default 204) and empty content.
	AssertNoContent(status ...int) Response
	// AssertMovedPermanently asserts that the response has a 301 Moved Permanently status code.
	AssertMovedPermanently() Response
	// AssertFound asserts that the response has a 302 Found status code.
	AssertFound() Response
	// AssertNotModified asserts that the response has a 304 Not Modified status code.
	AssertNotModified() Response
	// AssertPartialContent asserts that the response has a 206 Partial Content status code.
	AssertPartialContent() Response
	// AssertTemporaryRedirect asserts that the response has a 307 Temporary Redirect status code.
	AssertTemporaryRedirect() Response
	// AssertBadRequest asserts that the response has a 400 Bad Request status code.
	AssertBadRequest() Response
	// AssertUnauthorized asserts that the response has a 401 Unauthorized status code.
	AssertUnauthorized() Response
	// AssertPaymentRequired asserts that the response has a 402 Payment Required status code.
	AssertPaymentRequired() Response
	// AssertForbidden asserts that the response has a 403 Forbidden status code.
	AssertForbidden() Response
	// AssertNotFound asserts that the response has a 404 Not Found status code.
	AssertNotFound() Response
	// AssertMethodNotAllowed asserts that the response has a 405 Method Not Allowed status code.
	AssertMethodNotAllowed() Response
	// AssertNotAcceptable asserts that the response has a 406 Not Acceptable status code.
	AssertNotAcceptable() Response
	// AssertConflict asserts that the response has a 409 Conflict status code.
	AssertConflict() Response
	// AssertRequestTimeout asserts that the response has a 408 Request Timeout status code.
	AssertRequestTimeout() Response
	// AssertGone asserts that the response has a 410 Gone status code.
	AssertGone() Response
	// AssertUnsupportedMediaType asserts that the response has a 415 Unsupported Media Type status code.
	AssertUnsupportedMediaType() Response
	// AssertUnprocessableEntity asserts that the response has a 422 Unprocessable Entity status code.
	AssertUnprocessableEntity() Response
	// AssertTooManyRequests asserts that the response has a 429 Too Many Requests status code.
	AssertTooManyRequests() Response
	// AssertInternalServerError asserts that the response has a 500 Internal Server Error status code.
	AssertInternalServerError() Response
	// AssertServiceUnavailable asserts that the response has a 503 Service Unavailable status code.
	AssertServiceUnavailable() Response
	// AssertHeader asserts that the given header exists and optionally matches the value.
	AssertHeader(headerName, value string) Response
	// AssertHeaderMissing asserts that the given header is not present.
	AssertHeaderMissing(string) Response
	// AssertCookie asserts that the given cookie exists and optionally matches the value.
	AssertCookie(name, value string) Response
	// AssertCookieExpired asserts that the given cookie has expired.
	AssertCookieExpired(string) Response
	// AssertCookieNotExpired asserts that the given cookie has not expired.
	AssertCookieNotExpired(string) Response
	// AssertCookieMissing asserts that the given cookie is not present.
	AssertCookieMissing(string) Response
	// AssertSuccessful asserts that the response status code is >= 200 and < 300.
	AssertSuccessful() Response
	// AssertServerError asserts that the response status code is >= 500 and < 600.
	AssertServerError() Response
	// AssertDontSee asserts that the given strings are not present in the response body.
	AssertDontSee(value []string, escaped ...bool) Response
	// AssertSee asserts that the given strings are present in the response body.
	AssertSee(value []string, escaped ...bool) Response
	// AssertSeeInOrder asserts that the given strings are present in the response body in order.
	AssertSeeInOrder(value []string, escaped ...bool) Response
	// AssertJson asserts that the response JSON contains the given data.
	AssertJson(map[string]any) Response
	// AssertExactJson asserts that the response JSON matches the given data exactly.
	AssertExactJson(map[string]any) Response
	// AssertJsonMissing asserts that the response JSON does not contain the given keys or values.
	AssertJsonMissing(map[string]any) Response
	// AssertFluentJson allows for complex JSON assertions using a callback.
	AssertFluentJson(func(json AssertableJSON)) Response
}
