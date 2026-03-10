package telemetry

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

const (
	samplerAlwaysOn     = "always_on"
	samplerAlwaysOff    = "always_off"
	samplerTraceIDRatio = "traceidratio"
)

func newTraceSampler(cfg SamplerConfig) sdktrace.Sampler {
	var sampler sdktrace.Sampler

	switch cfg.Type {
	case samplerAlwaysOff:
		sampler = sdktrace.NeverSample()
	case samplerTraceIDRatio:
		sampler = sdktrace.TraceIDRatioBased(cfg.Ratio)
	default:
		sampler = sdktrace.AlwaysSample()
	}

	if cfg.Parent {
		return sdktrace.ParentBased(sampler)
	}

	return sampler
}
