package grpc

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/grpc"
)

type GrpcRunner struct {
	config config.Config
	grpc   grpc.Grpc
}

func NewGrpcRunner(config config.Config, grpc grpc.Grpc) *GrpcRunner {
	return &GrpcRunner{
		config: config,
		grpc:   grpc,
	}
}

func (r *GrpcRunner) Signature() string {
	return "grpc"
}

func (r *GrpcRunner) ShouldRun() bool {
	return r.grpc != nil && r.config.GetString("grpc.host") != "" && r.config.GetBool("app.auto_run", true)
}

func (r *GrpcRunner) Run() error {
	return r.grpc.Run()
}

func (r *GrpcRunner) Shutdown() error {
	return r.grpc.Shutdown()
}
