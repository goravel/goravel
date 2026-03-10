package telemetry

import (
	"strings"
	"time"
)

type Protocol string

const (
	ProtocolGRPC         Protocol = "grpc"
	ProtocolHTTPProtobuf Protocol = "http/protobuf"
	ProtocolHTTPJSON     Protocol = "http/json"
)

const defaultOTLPTimeout = 10 * time.Second

func buildOTLPOptions[T any](
	cfg ExporterEntry,
	withEndpoint func(string) T,
	withInsecure func() T,
	withTimeout func(time.Duration) T,
	withHeaders func(map[string]string) T,
) []T {
	var opts []T

	if cfg.Endpoint != "" {
		endpoint := strings.TrimPrefix(cfg.Endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		opts = append(opts, withEndpoint(endpoint))
	}

	if cfg.Insecure {
		opts = append(opts, withInsecure())
	}

	timeout := defaultOTLPTimeout
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}
	opts = append(opts, withTimeout(timeout))

	if headers := cfg.Headers; len(headers) > 0 {
		opts = append(opts, withHeaders(headers))
	}

	return opts
}
