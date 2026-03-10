package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
)

type Application struct {
	config config.Config
	server *grpc.Server

	// Server Options
	unaryServerInterceptors []grpc.UnaryServerInterceptor
	serverStatsHandlers     []stats.Handler

	// Client Options
	unaryClientInterceptorGroups map[string][]grpc.UnaryClientInterceptor
	clientStatsHandlerGroups     map[string][]stats.Handler

	// Mutex protects the servers map
	mu      sync.RWMutex
	servers map[string]*grpc.ClientConn
}

func NewApplication(config config.Config) *Application {
	return &Application{
		config:                       config,
		servers:                      make(map[string]*grpc.ClientConn),
		unaryServerInterceptors:      make([]grpc.UnaryServerInterceptor, 0),
		serverStatsHandlers:          make([]stats.Handler, 0),
		unaryClientInterceptorGroups: make(map[string][]grpc.UnaryClientInterceptor),
		clientStatsHandlerGroups:     make(map[string][]stats.Handler),
	}
}

// DEPRECATED: Use Connect instead, will be removed in v1.18.
func (r *Application) Client(ctx context.Context, name string) (*grpc.ClientConn, error) {
	return r.Connect(name)
}

func (r *Application) Connect(server string) (*grpc.ClientConn, error) {
	r.mu.RLock()
	conn, ok := r.servers[server]
	r.mu.RUnlock()

	if ok {
		// If connection exists and is healthy, return it immediately
		if conn.GetState() != connectivity.Shutdown {
			return conn, nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-Check: Someone else might have created it while we waited for the lock
	if conn, ok = r.servers[server]; ok {
		if conn.GetState() != connectivity.Shutdown {
			return conn, nil
		}
		// Found a Shutdown connection. Close and remove it immediately.
		// This prevents stale connections from lingering if the subsequent creation fails.
		_ = conn.Close()
		delete(r.servers, server)
	}

	host := r.config.GetString(fmt.Sprintf("grpc.servers.%s.host", server))
	if host == "" {
		return nil, errors.GrpcEmptyClientHost
	}
	if !strings.Contains(host, ":") {
		port := r.config.GetString(fmt.Sprintf("grpc.servers.%s.port", server))
		if port == "" {
			return nil, errors.GrpcEmptyClientPort
		}

		host += ":" + port
	}

	interceptorKeys, ok := r.config.Get(fmt.Sprintf("grpc.servers.%s.interceptors", server)).([]string)
	if !ok {
		return nil, errors.GrpcInvalidInterceptorsType.Args(server)
	}

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if interceptors := r.getClientInterceptors(interceptorKeys); len(interceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	statsHandlerKeys, ok := r.config.Get(fmt.Sprintf("grpc.servers.%s.stats_handlers", server)).([]string)
	if ok {
		if handlers := r.getClientStatsHandlers(statsHandlerKeys); len(handlers) > 0 {
			for _, h := range handlers {
				if h == nil {
					continue
				}
				dialOpts = append(dialOpts, grpc.WithStatsHandler(h))
			}
		}
	}

	newConn, err := grpc.NewClient(host, dialOpts...)
	if err != nil {
		return nil, err
	}

	r.servers[server] = newConn

	return newConn, nil
}

func (r *Application) Listen(l net.Listener) error {
	color.Green().Println("[GRPC] Listening on: " + l.Addr().String())
	return r.Server().Serve(l)
}

func (r *Application) Run(host ...string) error {
	if len(host) == 0 {
		defaultHost := r.config.GetString("grpc.host")
		if defaultHost == "" {
			return errors.GrpcEmptyServerHost
		}

		if !strings.Contains(defaultHost, ":") {
			defaultPort := r.config.GetString("grpc.port")
			if defaultPort == "" {
				return errors.GrpcEmptyServerPort
			}
			defaultHost += ":" + defaultPort
		}

		host = append(host, defaultHost)
	}

	listen, err := net.Listen("tcp", host[0])
	if err != nil {
		return err
	}

	color.Green().Println("[GRPC] Listening on: " + host[0])
	return r.Server().Serve(listen)
}

func (r *Application) Server() *grpc.Server {
	if r.server != nil {
		return r.server
	}

	var opts []grpc.ServerOption

	if len(r.unaryServerInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(r.unaryServerInterceptors...))
	}

	for _, h := range r.serverStatsHandlers {
		if h == nil {
			continue
		}
		opts = append(opts, grpc.StatsHandler(h))
	}

	r.server = grpc.NewServer(opts...)
	return r.server
}

func (r *Application) Shutdown(force ...bool) error {
	if r.server != nil {
		if len(force) > 0 && force[0] {
			r.server.Stop()
		} else {
			r.server.GracefulStop()
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, conn := range r.servers {
		_ = conn.Close()
	}

	// Clear the map to allow Garbage Collection
	r.servers = make(map[string]*grpc.ClientConn)

	return nil
}

func (r *Application) UnaryServerInterceptors(unaryServerInterceptors []grpc.UnaryServerInterceptor) {
	if r.server != nil {
		color.Warningln("[GRPC] Server already initialized; unary server interceptor registration ignored.")
		return
	}
	r.unaryServerInterceptors = append(r.unaryServerInterceptors, unaryServerInterceptors...)
}

func (r *Application) ServerStatsHandlers(handlers []stats.Handler) {
	if r.server != nil {
		color.Warningln("[GRPC] Server already initialized; server stats handler registration ignored.")
		return
	}
	r.serverStatsHandlers = append(r.serverStatsHandlers, handlers...)
}

func (r *Application) UnaryClientInterceptorGroups(groups map[string][]grpc.UnaryClientInterceptor) {
	for key, interceptors := range groups {
		r.unaryClientInterceptorGroups[key] = append(r.unaryClientInterceptorGroups[key], interceptors...)
	}
}

func (r *Application) ClientStatsHandlerGroups(groups map[string][]stats.Handler) {
	for key, handlers := range groups {
		r.clientStatsHandlerGroups[key] = append(r.clientStatsHandlerGroups[key], handlers...)
	}
}

func (r *Application) getClientInterceptors(keys []string) []grpc.UnaryClientInterceptor {
	var result []grpc.UnaryClientInterceptor
	for _, key := range keys {
		if group, ok := r.unaryClientInterceptorGroups[key]; ok {
			result = append(result, group...)
		}
	}
	return result
}

func (r *Application) getClientStatsHandlers(keys []string) []stats.Handler {
	var result []stats.Handler
	for _, key := range keys {
		if group, ok := r.clientStatsHandlerGroups[key]; ok {
			result = append(result, group...)
		}
	}
	return result
}
