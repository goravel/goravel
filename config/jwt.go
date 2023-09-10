package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("jwt", map[string]any{
		// JWT Authentication Secret
		//
		// Don't forget to set this in your .env file, as it will be used to sign
		// your tokens. A helper command is provided for this:
		// `go run . artisan jwt:secret`
		"secret": config.Env("JWT_SECRET", ""),

		// JWT time to live
		//
		// Specify the length of time (in minutes) that the token will be valid for.
		// Defaults to 1 hour.
		//
		// You can also set this to 0, to yield a never expiring token.
		// Some people may want this behaviour for e.g. a mobile app.
		// This is not particularly recommended, so make sure you have appropriate
		// systems in place to revoke the token if necessary.
		"ttl": config.Env("JWT_TTL", 60),

		// Refresh time to live
		//
		// Specify the length of time (in minutes) that the token can be refreshed
		// within. I.E. The user can refresh their token within a 2 week window of
		// the original token being created until they must re-authenticate.
		// Defaults to 2 weeks.
		//
		// You can also set this to 0, to yield an infinite refresh time.
		// Some may want this instead of never expiring tokens for e.g. a mobile app.
		// This is not particularly recommended, so make sure you have appropriate
		// systems in place to revoke the token if necessary.
		"refresh_ttl": config.Env("JWT_REFRESH_TTL", 20160),
	})
}
