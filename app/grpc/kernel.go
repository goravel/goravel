package grpc

import (
	"context"
	"runtime/debug"

	"github.com/goravel/framework/facades"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Kernel struct {
}

// The application's global GRPC interceptor stack.
// These middleware are run during every request to your application.
func (kernel Kernel) UnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryInterceptor(),
	}
}

// The application's client interceptor groups.
func (kernel Kernel) UnaryClientInterceptorGroups() map[string][]grpc.UnaryClientInterceptor {
	return map[string][]grpc.UnaryClientInterceptor{}
}

// recoveryInterceptor recovers from panics and returns an internal error
func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				facades.Log().Errorf("gRPC panic recovered: %v\n%s", r, debug.Stack())
				er	r = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
