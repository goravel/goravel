package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type KeyValue = attribute.KeyValue

// Attribute functions for creating span attributes.
var (
	Bool         = attribute.Bool
	BoolSlice    = attribute.BoolSlice
	Float64      = attribute.Float64
	Float64Slice = attribute.Float64Slice
	Int          = attribute.Int
	Int64        = attribute.Int64
	Int64Slice   = attribute.Int64Slice
	IntSlice     = attribute.IntSlice
	String       = attribute.String
	StringSlice  = attribute.StringSlice
	Stringer     = attribute.Stringer
)

// Span status codes.
const (
	CodeUnset = codes.Unset
	CodeError = codes.Error
	CodeOk    = codes.Ok
)

// SpanKind constants.
const (
	SpanKindUnspecified = trace.SpanKindUnspecified
	SpanKindInternal    = trace.SpanKindInternal
	SpanKindServer      = trace.SpanKindServer
	SpanKindClient      = trace.SpanKindClient
	SpanKindProducer    = trace.SpanKindProducer
	SpanKindConsumer    = trace.SpanKindConsumer
)

// Span start options.
var (
	WithAttributes = trace.WithAttributes
	WithLinks      = trace.WithLinks
	WithSpanKind   = trace.WithSpanKind
	WithTimestamp  = trace.WithTimestamp
)

// Propagation carriers for context propagation.
var (
	PropagationHeaderCarrier = func(h http.Header) propagation.HeaderCarrier {
		return propagation.HeaderCarrier(h)
	}
	PropagationMapCarrier = func(m map[string]string) propagation.MapCarrier {
		return propagation.MapCarrier(m)
	}
)
