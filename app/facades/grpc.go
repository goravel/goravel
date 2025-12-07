package facades

import (
	"github.com/goravel/framework/contracts/grpc"
)

func Grpc() grpc.Grpc {
	return App().MakeGrpc()
}
