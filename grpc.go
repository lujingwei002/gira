package gira

import "google.golang.org/grpc"

type GrpcServer interface {
	RegisterGrpc(f func(server *grpc.Server) error) error
}
