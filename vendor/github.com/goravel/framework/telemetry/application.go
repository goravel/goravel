package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	otellog "go.opentelemetry.io/otel/log"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/goravel/framework/contracts/telemetry"
	"github.com/goravel/framework/errors"
)

var _ telemetry.Telemetry = (*Application)(nil)

type Application struct {
	loggerProvider otellog.LoggerProvider
	meterProvider  otelmetric.MeterProvider
	tracerProvider oteltrace.TracerProvider
	propagator     propagation.TextMapPropagator
	shutdownFuncs  []ShutdownFunc
}

func NewApplication(cfg Config) (*Application, error) {
	propagator, err := newCompositeTextMapPropagator(cfg.Propagators)
	if err != nil {
		return nil, err
	}
	otel.SetTextMapPropagator(propagator)

	ctx := context.Background()
	resource, err := newResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	traceProvider, traceShutdown, err := NewTracerProvider(ctx, cfg, sdktrace.WithResource(resource))
	if err != nil {
		return nil, err
	}

	meterProvider, metricShutdown, err := NewMeterProvider(ctx, cfg, sdkmetric.WithResource(resource))
	if err != nil {
		_ = traceShutdown(ctx)
		return nil, err
	}

	loggerProvider, loggerShutdown, err := NewLoggerProvider(ctx, cfg, sdklog.WithResource(resource))
	if err != nil {
		_ = traceShutdown(ctx)
		_ = metricShutdown(ctx)
		return nil, err
	}

	return &Application{
		loggerProvider: loggerProvider,
		meterProvider:  meterProvider,
		tracerProvider: traceProvider,
		propagator:     propagator,
		shutdownFuncs: []ShutdownFunc{
			traceShutdown,
			metricShutdown,
			loggerShutdown,
		},
	}, nil
}

func (r *Application) Logger(name string, opts ...otellog.LoggerOption) otellog.Logger {
	return r.loggerProvider.Logger(name, opts...)
}

func (r *Application) Meter(name string, opts ...otelmetric.MeterOption) otelmetric.Meter {
	return r.meterProvider.Meter(name, opts...)
}

func (r *Application) MeterProvider() otelmetric.MeterProvider {
	return r.meterProvider
}

func (r *Application) Propagator() propagation.TextMapPropagator {
	return r.propagator
}

func (r *Application) Shutdown(ctx context.Context) error {
	var err error

	for _, fn := range r.shutdownFuncs {
		if fn == nil {
			continue
		}
		if e := fn(ctx); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func (r *Application) Tracer(name string, opts ...oteltrace.TracerOption) oteltrace.Tracer {
	return r.tracerProvider.Tracer(name, opts...)
}

func (r *Application) TracerProvider() oteltrace.TracerProvider {
	return r.tracerProvider
}
