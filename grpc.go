package gira

import "google.golang.org/grpc"

type GrpcServer interface {
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
}

type GrpcServerComponent interface {
	GetGrpcServer() GrpcServer
	GetGrpcService(name string) (svr interface{}, ok bool)
}
