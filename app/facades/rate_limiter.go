package facades

import (
	"github.com/goravel/framework/contracts/http"
)

func RateLimiter() http.RateLimiter {
	return App().MakeRateLimiter()
}
