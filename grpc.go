package gira

import "google.golang.org/grpc"

const GRPC_CATALOG_KEY = "catalog-key"

type GrpcServer interface {
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
	GetServer(name string) (svr interface{}, ok bool)
}

type GrpcServerComponent interface {
	GetGrpcServer() GrpcServer
}
