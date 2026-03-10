package telemetry

import "context"

type ShutdownFunc = func(context.Context) error

// NoopShutdown returns a ShutdownFunc that does nothing and returns nil.
func NoopShutdown() ShutdownFunc {
	return func(context.Context) error { return nil }
}
