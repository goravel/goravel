package http

import (
	"net/http"
)

const (
	DefaultAbortStatus = http.StatusBadRequest

	StatusContinue           = http.StatusContinue
	StatusSwitchingProtocols = http.StatusSwitchingProtocols
	StatusProcessing         = http.StatusProcessing
	StatusEarlyHints         = http.StatusEarlyHints

	StatusOK                   = http.StatusOK
	StatusCreated              = http.StatusCreated
	StatusAccepted             = http.StatusAccepted
	StatusNonAuthoritativeInfo = http.StatusNonAuthoritativeInfo
	StatusNoContent            = http.StatusNoContent
	StatusResetContent         = http.StatusResetContent
	StatusPartialContent       = http.StatusPartialContent
	StatusMultiStatus          = http.StatusMultiStatus
	StatusAlreadyReported      = http.StatusAlreadyReported
	StatusIMUsed               = http.StatusIMUsed

	StatusMultipleChoices   = http.StatusMultipleChoices
	StatusMovedPermanently  = http.StatusMovedPermanently
	StatusFound             = http.StatusFound
	StatusSeeOther          = http.StatusSeeOther
	StatusNotModified       = http.StatusNotModified
	StatusUseProxy          = http.StatusUseProxy
	StatusTemporaryRedirect = http.StatusTemporaryRedirect
	StatusPermanentRedirect = http.StatusPermanentRedirect

	StatusBadRequest                   = http.StatusBadRequest
	StatusUnauthorized                 = http.StatusUnauthorized
	StatusPaymentRequired              = http.StatusPaymentRequired
	StatusForbidden                    = http.StatusForbidden
	StatusNotFound                     = http.StatusNotFound
	StatusMethodNotAllowed             = http.StatusMethodNotAllowed
	StatusNotAcceptable                = http.StatusNotAcceptable
	StatusProxyAuthRequired            = http.StatusProxyAuthRequired
	StatusRequestTimeout               = http.StatusRequestTimeout
	StatusConflict                     = http.StatusConflict
	StatusGone                         = http.StatusGone
	StatusLengthRequired               = http.StatusLengthRequired
	StatusPreconditionFailed           = http.StatusPreconditionFailed
	StatusRequestEntityTooLarge        = http.StatusRequestEntityTooLarge
	StatusRequestURITooLong            = http.StatusRequestURITooLong
	StatusUnsupportedMediaType         = http.StatusUnsupportedMediaType
	StatusRequestedRangeNotSatisfiable = http.StatusRequestedRangeNotSatisfiable
	StatusExpectationFailed            = http.StatusExpectationFailed
	StatusTeapot                       = http.StatusTeapot
	StatusMisdirectedRequest           = http.StatusMisdirectedRequest
	StatusUnprocessableEntity          = http.StatusUnprocessableEntity
	StatusLocked                       = http.StatusLocked
	StatusFailedDependency             = http.StatusFailedDependency
	StatusTooEarly                     = http.StatusTooEarly
	StatusUpgradeRequired              = http.StatusUpgradeRequired
	StatusPreconditionRequired         = http.StatusPreconditionRequired
	StatusTooManyRequests              = http.StatusTooManyRequests
	StatusRequestHeaderFieldsTooLarge  = http.StatusRequestHeaderFieldsTooLarge
	StatusUnavailableForLegalReasons   = http.StatusUnavailableForLegalReasons

	StatusInternalServerError           = http.StatusInternalServerError
	StatusNotImplemented                = http.StatusNotImplemented
	StatusBadGateway                    = http.StatusBadGateway
	StatusServiceUnavailable            = http.StatusServiceUnavailable
	StatusGatewayTimeout                = http.StatusGatewayTimeout
	StatusHTTPVersionNotSupported       = http.StatusHTTPVersionNotSupported
	StatusVariantAlsoNegotiates         = http.StatusVariantAlsoNegotiates
	StatusInsufficientStorage           = http.StatusInsufficientStorage
	StatusLoopDetected                  = http.StatusLoopDetected
	StatusNotExtended                   = http.StatusNotExtended
	StatusNetworkAuthenticationRequired = http.StatusNetworkAuthenticationRequired
	StatusTokenMismatch                 = 419
)

var statusText = map[int]string{
	StatusContinue:           http.StatusText(StatusContinue),
	StatusSwitchingProtocols: http.StatusText(StatusSwitchingProtocols),
	StatusProcessing:         http.StatusText(StatusProcessing),
	StatusEarlyHints:         http.StatusText(StatusEarlyHints),

	StatusOK:                   http.StatusText(StatusOK),
	StatusCreated:              http.StatusText(StatusCreated),
	StatusAccepted:             http.StatusText(StatusAccepted),
	StatusNonAuthoritativeInfo: http.StatusText(StatusNonAuthoritativeInfo),
	StatusNoContent:            http.StatusText(StatusNoContent),
	StatusResetContent:         http.StatusText(StatusResetContent),
	StatusPartialContent:       http.StatusText(StatusPartialContent),
	StatusMultiStatus:          http.StatusText(StatusMultiStatus),
	StatusAlreadyReported:      http.StatusText(StatusAlreadyReported),
	StatusIMUsed:               http.StatusText(StatusIMUsed),

	StatusMultipleChoices:   http.StatusText(StatusMultipleChoices),
	StatusMovedPermanently:  http.StatusText(StatusMovedPermanently),
	StatusFound:             http.StatusText(StatusFound),
	StatusSeeOther:          http.StatusText(StatusSeeOther),
	StatusNotModified:       http.StatusText(StatusNotModified),
	StatusUseProxy:          http.StatusText(StatusUseProxy),
	StatusTemporaryRedirect: http.StatusText(StatusTemporaryRedirect),
	StatusPermanentRedirect: http.StatusText(StatusPermanentRedirect),

	StatusBadRequest:                   http.StatusText(StatusBadRequest),
	StatusUnauthorized:                 http.StatusText(StatusUnauthorized),
	StatusPaymentRequired:              http.StatusText(StatusPaymentRequired),
	StatusForbidden:                    http.StatusText(StatusForbidden),
	StatusNotFound:                     http.StatusText(StatusNotFound),
	StatusMethodNotAllowed:             http.StatusText(StatusMethodNotAllowed),
	StatusNotAcceptable:                http.StatusText(StatusNotAcceptable),
	StatusProxyAuthRequired:            http.StatusText(StatusProxyAuthRequired),
	StatusRequestTimeout:               http.StatusText(StatusRequestTimeout),
	StatusConflict:                     http.StatusText(StatusConflict),
	StatusGone:                         http.StatusText(StatusGone),
	StatusLengthRequired:               http.StatusText(StatusLengthRequired),
	StatusPreconditionFailed:           http.StatusText(StatusPreconditionFailed),
	StatusRequestEntityTooLarge:        http.StatusText(StatusRequestEntityTooLarge),
	StatusRequestURITooLong:            http.StatusText(StatusRequestURITooLong),
	StatusUnsupportedMediaType:         http.StatusText(StatusUnsupportedMediaType),
	StatusRequestedRangeNotSatisfiable: http.StatusText(StatusRequestedRangeNotSatisfiable),
	StatusExpectationFailed:            http.StatusText(StatusExpectationFailed),
	StatusTeapot:                       http.StatusText(StatusTeapot),
	StatusMisdirectedRequest:           http.StatusText(StatusMisdirectedRequest),
	StatusUnprocessableEntity:          http.StatusText(StatusUnprocessableEntity),
	StatusLocked:                       http.StatusText(StatusLocked),
	StatusFailedDependency:             http.StatusText(StatusFailedDependency),
	StatusTooEarly:                     http.StatusText(StatusTooEarly),
	StatusUpgradeRequired:              http.StatusText(StatusUpgradeRequired),
	StatusPreconditionRequired:         http.StatusText(StatusPreconditionRequired),
	StatusTooManyRequests:              http.StatusText(StatusTooManyRequests),
	StatusRequestHeaderFieldsTooLarge:  http.StatusText(StatusRequestHeaderFieldsTooLarge),
	StatusUnavailableForLegalReasons:   http.StatusText(StatusUnavailableForLegalReasons),

	StatusInternalServerError:           http.StatusText(StatusInternalServerError),
	StatusNotImplemented:                http.StatusText(StatusNotImplemented),
	StatusBadGateway:                    http.StatusText(StatusBadGateway),
	StatusServiceUnavailable:            http.StatusText(StatusServiceUnavailable),
	StatusGatewayTimeout:                http.StatusText(StatusGatewayTimeout),
	StatusHTTPVersionNotSupported:       http.StatusText(StatusHTTPVersionNotSupported),
	StatusVariantAlsoNegotiates:         http.StatusText(StatusVariantAlsoNegotiates),
	StatusInsufficientStorage:           http.StatusText(StatusInsufficientStorage),
	StatusLoopDetected:                  http.StatusText(StatusLoopDetected),
	StatusNotExtended:                   http.StatusText(StatusNotExtended),
	StatusNetworkAuthenticationRequired: http.StatusText(StatusNetworkAuthenticationRequired),
	StatusTokenMismatch:                 "CSRF token mismatch",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
