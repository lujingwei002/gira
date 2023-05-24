package gira

import "google.golang.org/grpc"

type GrpcServer interface {
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
	GetService(name string) (svr interface{}, ok bool)
}

type GrpcServerComponent interface {
	GetGrpcServer() GrpcServer
}
