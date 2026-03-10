package client

import "time"

type FactoryConfig struct {
	Default string            `json:"default_client"`
	Clients map[string]Config `json:"clients"`
}

type Config struct {
	// BaseUrl is the prefix for all requests made by this client.
	// Example: "https://goravel.dev"
	BaseUrl string `json:"base_url"`

	// Timeout specifies the time limit for requests made by this Client.
	// The timeout includes connection time, any redirects, and reading the response body.
	// A Timeout of zero means no timeout (not recommended).
	Timeout time.Duration `json:"timeout"`

	// MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts.
	// Increasing this helps performance when making many requests to distinct hosts.
	MaxIdleConns int `json:"max_idle_conns"`

	// MaxIdleConnsPerHost controls the maximum number of idle (keep-alive) connections
	// to keep per-host. This is the most critical setting for high-throughput clients
	// talking to a single backend service.
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host"`

	// MaxConnsPerHost limits the total number of connections (active + idle) per host.
	// Zero means no limit.
	MaxConnsPerHost int `json:"max_conns_per_host"`

	// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection
	// will remain idle before closing itself.
	IdleConnTimeout time.Duration `json:"idle_conn_timeout"`
}
