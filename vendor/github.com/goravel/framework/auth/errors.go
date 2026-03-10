package auth

import "github.com/goravel/framework/errors"

// These errors may be used by user's project, so we can't remove them.
var (
	ErrorRefreshTimeExceeded     = errors.AuthRefreshTimeExceeded
	ErrorTokenExpired            = errors.AuthTokenExpired
	ErrorNoPrimaryKeyField       = errors.AuthNoPrimaryKeyField
	ErrorEmptySecret             = errors.AuthEmptySecret
	ErrorTokenDisabled           = errors.AuthTokenDisabled
	ErrorParseTokenFirst         = errors.AuthParseTokenFirst
	ErrorInvalidClaims           = errors.AuthInvalidClaims
	ErrorInvalidToken            = errors.AuthInvalidToken
	ErrorInvalidKey              = errors.AuthInvalidKey
	ErrorUnsupportedDriverMethod = errors.AuthUnsupportedDriverMethod
)
