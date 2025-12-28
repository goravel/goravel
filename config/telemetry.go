package config

import (
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("telemetry", map[string]any{
		// Service Identification
		//
		// Defines the logical identity of the service. This information is attached
		// to every trace and metric, allowing observability platforms to group
		// data by service name, version, and environment.
		// Reference: https://opentelemetry.io/docs/specs/semconv/resource/
		"service": map[string]any{
			"name":        config.Env("APP_NAME", "goravel"),
			"version":     config.Env("APP_VERSION", ""),
			"environment": config.Env("APP_ENV", ""),
		},

		// Resource Attributes
		//
		// Additional custom attributes to attach to the global Resource object.
		// These attributes provide static metadata about the entity producing
		// telemetry (e.g., "k8s.pod.name", "region", "team", "service.instance.id").
		"resource": map[string]any{
			// "service.instance.id": config.Env("APP_INSTANCE_ID", ""),
		},

		// Context Propagation
		//
		// Defines how trace context is encoded and propagated across process
		// boundaries (e.g., HTTP headers).
		//
		// Supported:
		// - "tracecontext": W3C Trace Context (Standard)
		// - "baggage": W3C Baggage
		// - "b3": B3 Single Header (Zipkin)
		// - "b3multi": B3 Multi Header
		"propagators": config.Env("OTEL_PROPAGATORS", "tracecontext"),

		// Traces Configuration
		//
		// Configures the distributed tracing signal. Traces record the path of
		// a request as it propagates through the distributed system.
		"traces": map[string]any{
			// Exporter
			//
			// The name of the exporter definition in the "exporters" section below.
			// Set to "" to disable tracing.
			"exporter": config.Env("OTEL_TRACES_EXPORTER", "otlptrace"),

			// Sampler Configuration
			//
			// Controls which traces are recorded and exported. Sampling reduces
			// overhead and storage costs by only recording a subset of traces.
			"sampler": map[string]any{
				// If true, respects the sampling decision of the upstream service.
				"parent": config.Env("OTEL_TRACES_SAMPLER_PARENT", true),

				// Sampling Strategy:
				// - "always_on": Record every trace (Dev/Test).
				// - "always_off": Record nothing.
				// - "traceidratio": Probabilistic sampling based on the ratio.
				"type": config.Env("OTEL_TRACES_SAMPLER_TYPE", "always_on"),

				// The ratio for "traceidratio" sampling (0.0 to 1.0).
				// e.g., 0.1 records ~10% of traces.
				"ratio": config.Env("OTEL_TRACES_SAMPLER_RATIO", 0.05),
			},
		},

		// Metrics Configuration
		//
		// Configures the metric signal. Metrics are numerical data points
		// aggregated over time (e.g., request counters, memory usage).
		"metrics": map[string]any{
			// Exporter
			//
			// The name of the exporter definition in the "exporters" section below.
			// Set to "" to disable metrics.
			"exporter": config.Env("OTEL_METRICS_EXPORTER", "otlpmetric"),

			// Reader Configuration
			//
			// Configures the PeriodicReader, which collects and pushes metrics
			// to the exporter at a fixed interval.
			"reader": map[string]any{
				// Interval: How often metrics are pushed.
				// Format: Duration string (e.g., "60s", "1m", "500ms").
				"interval": config.GetString("OTEL_METRIC_EXPORT_INTERVAL", "60s"),

				// Timeout: Max time allowed for export before cancelling.
				// Format: Duration string (e.g., "30s", "10s").
				"timeout": config.GetString("OTEL_METRIC_EXPORT_TIMEOUT", "30s"),
			},
		},

		// Logs Configuration
		//
		// Configures the logging signal. Logs are textual records of events
		// linked to traces for correlation.
		"logs": map[string]any{
			// Exporter
			//
			// The name of the exporter definition in the "exporters" section below.
			// Set to "" to disable OTel logging.
			"exporter": config.Env("OTEL_LOGS_EXPORTER", "otlplog"),

			// Processor Configuration
			//
			// Configures the BatchLogProcessor, which batches logs before export.
			"processor": map[string]any{
				// Interval: How often logs are flushed.
				// Format: Duration string (e.g., "1s", "500ms").
				"interval": config.GetString("OTEL_LOG_EXPORT_INTERVAL", "1s"),

				// Timeout: Max time allowed for export before cancelling.
				// Format: Duration string (e.g., "30s").
				"timeout": config.GetString("OTEL_LOG_EXPORT_TIMEOUT", "30s"),
			},
		},

		// Exporters Configuration
		//
		// Defines the details for connecting to external telemetry backends.
		// These definitions are referenced by name in the signal sections above.
		//
		// Supported drivers: "otlp", "zipkin", "console", "custom"
		"exporters": map[string]any{
			// OTLP Trace Exporter
			// Reference: https://opentelemetry.io/docs/specs/otel/protocol/
			"otlptrace": map[string]any{
				"driver":   "otlp",
				"endpoint": config.Env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "http://localhost:4318"),

				// Protocol: "http/protobuf", "http/json" or "grpc".
				"protocol": config.Env("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "http/protobuf"),

				// Set to false to require TLS/SSL.
				"insecure": config.Env("OTEL_EXPORTER_OTLP_TRACES_INSECURE", true),

				// Timeout: Max time to wait for the backend to acknowledge.
				// Format: Duration string (e.g., "10s", "500ms").
				"timeout": config.GetString("OTEL_EXPORTER_OTLP_TRACES_TIMEOUT", "10s"),
			},

			// OTLP Metric Exporter
			"otlpmetric": map[string]any{
				"driver":   "otlp",
				"endpoint": config.Env("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "http://localhost:4318"),
				"protocol": config.Env("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL", "http/protobuf"),
				"insecure": config.Env("OTEL_EXPORTER_OTLP_METRICS_INSECURE", true),

				// Timeout: Max time to wait for the backend to acknowledge.
				// Format: Duration string (e.g., "10s", "500ms").
				"timeout": config.GetString("OTEL_EXPORTER_OTLP_METRICS_TIMEOUT", "10s"),

				// Metric Temporality: "cumulative" or "delta".
				// - "cumulative": Standard for Prometheus (counts never reset).
				// - "delta": Standard for Datadog/StatsD (counts per interval).
				"metric_temporality": config.Env("OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY", "cumulative"),
			},

			// OTLP Log Exporter
			"otlplog": map[string]any{
				"driver":   "otlp",
				"endpoint": config.Env("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT", "http://localhost:4318"),
				"protocol": config.Env("OTEL_EXPORTER_OTLP_LOGS_PROTOCOL", "http/protobuf"),
				"insecure": config.Env("OTEL_EXPORTER_OTLP_LOGS_INSECURE", true),

				// Timeout: Max time to wait for the backend to acknowledge.
				// Format: Duration string (e.g., "10s", "500ms").
				"timeout": config.GetString("OTEL_EXPORTER_OTLP_LOGS_TIMEOUT", "10s"),
			},

			// Zipkin Trace Exporter (Tracing only)
			"zipkin": map[string]any{
				"driver":   "zipkin",
				"endpoint": config.Env("OTEL_EXPORTER_ZIPKIN_ENDPOINT", "http://localhost:9411/api/v2/spans"),
			},

			// Console Exporter (Debugging)
			// Prints telemetry data to stdout.
			"console": map[string]any{
				"driver": "console",
			},
		},
	})
}
