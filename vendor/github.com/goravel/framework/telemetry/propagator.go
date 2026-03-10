package telemetry

import (
	"strings"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"

	"github.com/goravel/framework/errors"
)

const (
	propagatorTraceContext = "tracecontext"
	propagatorBaggage      = "baggage"
	propagatorB3           = "b3"
	propagatorB3Multi      = "b3multi"
)

func newCompositeTextMapPropagator(nameStr string) (propagation.TextMapPropagator, error) {
	nameStr = strings.TrimSpace(nameStr)
	if nameStr == "" {
		return nil, errors.TelemetryPropagatorRequired
	}

	var propagators []propagation.TextMapPropagator

	for _, name := range strings.Split(nameStr, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		switch name {
		case propagatorTraceContext:
			propagators = append(propagators, propagation.TraceContext{})
		case propagatorBaggage:
			propagators = append(propagators, propagation.Baggage{})
		case propagatorB3:
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)))
		case propagatorB3Multi:
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
		default:
			return nil, errors.TelemetryUnsupportedPropagator.Args(name)
		}
	}

	return propagation.NewCompositeTextMapPropagator(propagators...), nil
}
