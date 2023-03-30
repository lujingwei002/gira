// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: doc/grpc/hall.proto

package hall_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Upstream_SayHello_FullMethodName   = "/hall_grpc.Upstream/SayHello"
	Upstream_DataStream_FullMethodName = "/hall_grpc.Upstream/DataStream"
)

// UpstreamClient is the client API for Upstream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpstreamClient interface {
	// SayHello 方法
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	DataStream(ctx context.Context, opts ...grpc.CallOption) (Upstream_DataStreamClient, error)
}

type upstreamClient struct {
	cc grpc.ClientConnInterface
}

func NewUpstreamClient(cc grpc.ClientConnInterface) UpstreamClient {
	return &upstreamClient{cc}
}

func (c *upstreamClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, Upstream_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *upstreamClient) DataStream(ctx context.Context, opts ...grpc.CallOption) (Upstream_DataStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Upstream_ServiceDesc.Streams[0], Upstream_DataStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &upstreamDataStreamClient{stream}
	return x, nil
}

type Upstream_DataStreamClient interface {
	Send(*StreamDataRequest) error
	Recv() (*StreamDataResponse, error)
	grpc.ClientStream
}

type upstreamDataStreamClient struct {
	grpc.ClientStream
}

func (x *upstreamDataStreamClient) Send(m *StreamDataRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *upstreamDataStreamClient) Recv() (*StreamDataResponse, error) {
	m := new(StreamDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UpstreamServer is the server API for Upstream service.
// All implementations must embed UnimplementedUpstreamServer
// for forward compatibility
type UpstreamServer interface {
	// SayHello 方法
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	DataStream(Upstream_DataStreamServer) error
	mustEmbedUnimplementedUpstreamServer()
}

// UnimplementedUpstreamServer must be embedded to have forward compatible implementations.
type UnimplementedUpstreamServer struct {
}

func (UnimplementedUpstreamServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedUpstreamServer) DataStream(Upstream_DataStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method DataStream not implemented")
}
func (UnimplementedUpstreamServer) mustEmbedUnimplementedUpstreamServer() {}

// UnsafeUpstreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpstreamServer will
// result in compilation errors.
type UnsafeUpstreamServer interface {
	mustEmbedUnimplementedUpstreamServer()
}

func RegisterUpstreamServer(s grpc.ServiceRegistrar, srv UpstreamServer) {
	s.RegisterService(&Upstream_ServiceDesc, srv)
}

func _Upstream_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpstreamServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Upstream_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpstreamServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Upstream_DataStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UpstreamServer).DataStream(&upstreamDataStreamServer{stream})
}

type Upstream_DataStreamServer interface {
	Send(*StreamDataResponse) error
	Recv() (*StreamDataRequest, error)
	grpc.ServerStream
}

type upstreamDataStreamServer struct {
	grpc.ServerStream
}

func (x *upstreamDataStreamServer) Send(m *StreamDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *upstreamDataStreamServer) Recv() (*StreamDataRequest, error) {
	m := new(StreamDataRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Upstream_ServiceDesc is the grpc.ServiceDesc for Upstream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Upstream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hall_grpc.Upstream",
	HandlerType: (*UpstreamServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Upstream_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DataStream",
			Handler:       _Upstream_DataStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "doc/grpc/hall.proto",
}

const (
	Hall_UserInstead_FullMethodName = "/hall_grpc.Hall/UserInstead"
)

// HallClient is the client API for Hall service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HallClient interface {
	// SayHello 方法
	UserInstead(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error)
}

type hallClient struct {
	cc grpc.ClientConnInterface
}

func NewHallClient(cc grpc.ClientConnInterface) HallClient {
	return &hallClient{cc}
}

func (c *hallClient) UserInstead(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error) {
	out := new(UserInsteadResponse)
	err := c.cc.Invoke(ctx, Hall_UserInstead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HallServer is the server API for Hall service.
// All implementations must embed UnimplementedHallServer
// for forward compatibility
type HallServer interface {
	// SayHello 方法
	UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error)
	mustEmbedUnimplementedHallServer()
}

// UnimplementedHallServer must be embedded to have forward compatible implementations.
type UnimplementedHallServer struct {
}

func (UnimplementedHallServer) UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInstead not implemented")
}
func (UnimplementedHallServer) mustEmbedUnimplementedHallServer() {}

// UnsafeHallServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HallServer will
// result in compilation errors.
type UnsafeHallServer interface {
	mustEmbedUnimplementedHallServer()
}

func RegisterHallServer(s grpc.ServiceRegistrar, srv HallServer) {
	s.RegisterService(&Hall_ServiceDesc, srv)
}

func _Hall_UserInstead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInsteadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).UserInstead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_UserInstead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).UserInstead(ctx, req.(*UserInsteadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Hall_ServiceDesc is the grpc.ServiceDesc for Hall service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hall_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hall_grpc.Hall",
	HandlerType: (*HallServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserInstead",
			Handler:    _Hall_UserInstead_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "doc/grpc/hall.proto",
}
