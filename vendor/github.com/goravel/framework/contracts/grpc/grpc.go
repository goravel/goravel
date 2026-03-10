package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type Grpc interface {
	// DEPRECATED: Use Connect instead, will be removed in v1.18.
	Client(ctx context.Context, name string) (*grpc.ClientConn, error)
	// ClientStatsHandlerGroups sets the gRPC client stats handler groups.
	ClientStatsHandlerGroups(map[string][]stats.Handler)
	// Connect gets a gRPC client connection to the given server.
	// The server connection will be cached to improve performance.
	Connect(server string) (*grpc.ClientConn, error)
	// Listen starts the gRPC server with the given listener.
	Listen(l net.Listener) error
	// Run starts the gRPC server.
	Run(host ...string) error
	// Server gets the gRPC server instance.
	Server() *grpc.Server
	// ServerStatsHandlers sets the gRPC server stats handlers.
	ServerStatsHandlers([]stats.Handler)
	// Shutdown stops the gRPC server.
	Shutdown(force ...bool) error
	// UnaryClientInterceptorGroups sets the gRPC client interceptor groups.
	UnaryClientInterceptorGroups(map[string][]grpc.UnaryClientInterceptor)
	// UnaryServerInterceptors sets the gRPC server interceptors.
	UnaryServerInterceptors([]grpc.UnaryServerInterceptor)
}
