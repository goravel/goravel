package http

type RateLimiter interface {
	// For register a new rate limiter.
	For(name string, callback func(ctx Context) Limit)
	// ForWithLimits register a new rate limiter with limits.
	ForWithLimits(name string, callback func(ctx Context) []Limit)
	// Limiter get a rate limiter instance by name.
	Limiter(name string) func(ctx Context) []Limit
}

type Limit interface {
	// By set the signature key name for the rate limiter.
	By(key string) Limit
	// GetKey get the signature key name for the rate limiter.
	GetKey() string
	// GetResponse get the response callback that should be used.
	GetResponse() func(ctx Context)
	// GetStore get the store instance.
	GetStore() Store
	// Response set the response callback that should be used.
	Response(func(ctx Context)) Limit
}
