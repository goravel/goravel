package telemetry

import (
	"context"

	otellog "go.opentelemetry.io/otel/log"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Telemetry interface {
	// Logger returns a log.Logger instance for emitting structured log records under the given instrumentation name.
	// Optional log.LoggerOption parameters allow customization of logger behavior.
	Logger(name string, opts ...otellog.LoggerOption) otellog.Logger

	// Meter returns a metric.Meter instance for recording metrics under the given instrumentation name.
	// The optional metric.MeterOption parameters allow further customization.
	Meter(name string, opts ...otelmetric.MeterOption) otelmetric.Meter

	// MeterProvider returns the underlying metric.MeterProvider responsible for creating meters.
	MeterProvider() otelmetric.MeterProvider

	// Propagator returns the configured TextMapPropagator used to inject and extract
	// context across service boundaries for distributed tracing.
	Propagator() propagation.TextMapPropagator

	// Shutdown flushes any pending telemetry data and releases associated resources.
	// This should typically be called during application shutdown.
	Shutdown(ctx context.Context) error

	// Tracer returns a trace.Tracer instance for the given instrumentation name.
	// Optional trace.TracerOption parameters allow customization of tracer behavior.
	Tracer(name string, opts ...oteltrace.TracerOption) oteltrace.Tracer

	// TracerProvider returns the underlying trace.TracerProvider responsible for creating tracers.
	TracerProvider() oteltrace.TracerProvider
}
