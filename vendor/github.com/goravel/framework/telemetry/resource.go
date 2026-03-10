package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/goravel/framework/errors"
)

func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	serviceCfg := cfg.Service
	serviceName := serviceCfg.Name
	if serviceName == "" {
		return nil, errors.TelemetryServiceNameRequired
	}

	var attrsList []attribute.KeyValue
	attrsList = append(attrsList, semconv.ServiceName(serviceName))

	if serviceCfg.Version != "" {
		attrsList = append(attrsList, semconv.ServiceVersion(serviceCfg.Version))
	}

	if serviceCfg.Environment != "" {
		attrsList = append(attrsList, semconv.DeploymentEnvironmentName(serviceCfg.Environment))
	}

	// TODO: revisit the resource type, Need to consider any [In Next PR]
	for k, v := range cfg.Resource {
		if k != "" {
			attrsList = append(attrsList, attribute.String(k, v))
		}
	}

	resourceOptions := []resource.Option{
		resource.WithAttributes(attrsList...),

		// Add automatic detection options
		resource.WithFromEnv(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
	}

	detected, err := resource.New(ctx, resourceOptions...)
	if err != nil {
		return nil, err
	}

	return resource.Merge(resource.Default(), detected)
}
