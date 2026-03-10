package http

import (
	"net/http"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/contracts/validation"
)

type ContextRequest interface {
	// Cookie retrieves the value of the specified cookie by its key.
	Cookie(key string, defaultValue ...string) string
	// Header retrieves the value of the specified HTTP header by its key.
	// If the header is not found, it returns the optional default value (if provided).
	Header(key string, defaultValue ...string) string
	// Headers return all the HTTP headers of the request.
	Headers() http.Header
	// Method retrieves the HTTP request method (e.g., GET, POST, PUT).
	Method() string
	// Name retrieves the name of the route for the request.
	Name() string
	// OriginPath retrieves the original path of the request: /users/{id}
	OriginPath() string
	// Path retrieves the current path information for the request: /users/1
	Path() string
	// Url retrieves the URL (excluding the query string) for the request.
	Url() string
	// FullUrl retrieves the full URL, including the query string, for the request.
	FullUrl() string
	// Info retrieves the route information for the request.
	Info() Info
	// Ip retrieves the client's IP address.
	Ip() string
	// Host retrieves the host name.
	Host() string
	// All retrieves data from JSON, form, and query parameters.
	All() map[string]any
	// Bind retrieve json and bind to obj
	Bind(obj any) error
	// BindQuery bind query parameters to obj
	BindQuery(obj any) error
	// Route retrieves a route parameter from the request path (e.g., /users/{id}).
	Route(key string) string
	// RouteInt retrieves a route parameter from the request path and attempts to parse it as an integer.
	RouteInt(key string) int
	// RouteInt64 retrieves a route parameter from the request path and attempts to parse it as a 64-bit integer.
	RouteInt64(key string) int64
	// Query retrieves a query string parameter from the request (e.g., /users?id=1).
	Query(key string, defaultValue ...string) string
	// QueryInt retrieves a query string parameter from the request and attempts to parse it as an integer.
	QueryInt(key string, defaultValue ...int) int
	// QueryInt64 retrieves a query string parameter from the request and attempts to parse it as a 64-bit integer.
	QueryInt64(key string, defaultValue ...int64) int64
	// QueryBool retrieves a query string parameter from the request and attempts to parse it as a boolean.
	QueryBool(key string, defaultValue ...bool) bool
	// QueryArray retrieves a query string parameter from the request and returns it as a slice of strings.
	QueryArray(key string) []string
	// QueryMap retrieves a query string parameter from the request and returns it as a map of key-value pairs.
	QueryMap(key string) map[string]string
	// Queries returns all the query string parameters from the request as a map of key-value pairs.
	Queries() map[string]string

	// HasSession checks if the request has a session.
	HasSession() bool
	// Session retrieves the session associated with the request.
	Session() session.Session
	// SetSession sets the session associated with the request.
	SetSession(session session.Session) ContextRequest

	// Input retrieves data from the request in the following order: JSON, form, query, and route parameters.
	Input(key string, defaultValue ...string) string
	InputArray(key string, defaultValue ...[]string) []string
	InputMap(key string, defaultValue ...map[string]any) map[string]any
	InputMapArray(key string, defaultValue ...[]map[string]any) []map[string]any
	InputInt(key string, defaultValue ...int) int
	InputInt64(key string, defaultValue ...int64) int64
	InputBool(key string, defaultValue ...bool) bool
	// File retrieves a file by its key from the request.
	File(name string) (filesystem.File, error)
	Files(name string) ([]filesystem.File, error)

	// Abort aborts the request with the specified HTTP status code, default is 400.
	Abort(code ...int)
	// AbortWithStatus aborts the request with the specified HTTP status code.
	// DEPRECATED: Use Abort instead.
	AbortWithStatus(code int)
	// AbortWithStatusJson aborts the request with the specified HTTP status code
	// and returns a JSON response object.
	// DEPRECATED: Use Response().Json().Abort() instead.
	AbortWithStatusJson(code int, jsonObj any)
	// Next skips the current request handler, allowing the next middleware or handler to be executed.
	Next()
	// Origin retrieves the underlying *http.Request object for advanced request handling.
	Origin() *http.Request

	// Validate performs request data validation using specified rules and options.
	Validate(rules map[string]string, options ...validation.Option) (validation.Validator, error)
	// ValidateRequest validates the request data against a pre-defined FormRequest structure
	// and returns validation errors, if any.
	ValidateRequest(request FormRequest) (validation.Errors, error)
}

type FormRequest interface {
	// Authorize determine if the user is authorized to make this request.
	Authorize(ctx Context) error
	// Rules get the validation rules that apply to the request.
	Rules(ctx Context) map[string]string
}

type FormRequestWithFilters interface {
	// Filters get the custom filters that apply to the request.
	Filters(ctx Context) map[string]string
}

type FormRequestWithMessages interface {
	// Messages get the validation messages that apply to the request.
	Messages(ctx Context) map[string]string
}

type FormRequestWithAttributes interface {
	// Attributes get custom attributes for validator errors.
	Attributes(ctx Context) map[string]string
}

type FormRequestWithPrepareForValidation interface {
	// PrepareForValidation prepare the data for validation.
	PrepareForValidation(ctx Context, data validation.Data) error
}

type Info struct {
	Handler string `json:"handler"`
	Method  string `json:"method"`
	Name    string `json:"name"`
	Path    string `json:"path"`
}
