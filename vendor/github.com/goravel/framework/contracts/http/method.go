package http

import (
	"net/http"
)

const (
	MethodHead       = http.MethodHead
	MethodGet        = http.MethodGet
	MethodPost       = http.MethodPost
	MethodPut        = http.MethodPut
	MethodDelete     = http.MethodDelete
	MethodPatch      = http.MethodPatch
	MethodOptions    = http.MethodOptions
	MethodConnect    = http.MethodConnect
	MethodTrace      = http.MethodTrace
	MethodAny        = "ANY"
	MethodResource   = "RESOURCE"
	MethodStatic     = "STATIC"
	MethodStaticFile = "STATIC_FILE"
	MethodStaticFS   = "STATIC_FS"
)
