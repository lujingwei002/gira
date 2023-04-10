// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
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
	Hall_UserInstead_FullMethodName  = "/hall_grpc.Hall/UserInstead"
	Hall_PushStream_FullMethodName   = "/hall_grpc.Hall/PushStream"
	Hall_MustPush_FullMethodName     = "/hall_grpc.Hall/MustPush"
	Hall_ClientStream_FullMethodName = "/hall_grpc.Hall/ClientStream"
	Hall_GateStream_FullMethodName   = "/hall_grpc.Hall/GateStream"
)

// HallClient is the client API for Hall service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HallClient interface {
	// SayHello 方法
	UserInstead(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error)
	PushStream(ctx context.Context, opts ...grpc.CallOption) (Hall_PushStreamClient, error)
	MustPush(ctx context.Context, in *MustPushRequest, opts ...grpc.CallOption) (*MustPushResponse, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (Hall_ClientStreamClient, error)
	GateStream(ctx context.Context, opts ...grpc.CallOption) (Hall_GateStreamClient, error)
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

func (c *hallClient) PushStream(ctx context.Context, opts ...grpc.CallOption) (Hall_PushStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hall_ServiceDesc.Streams[0], Hall_PushStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hallPushStreamClient{stream}
	return x, nil
}

type Hall_PushStreamClient interface {
	Send(*PushStreamRequest) error
	CloseAndRecv() (*PushStreamResponse, error)
	grpc.ClientStream
}

type hallPushStreamClient struct {
	grpc.ClientStream
}

func (x *hallPushStreamClient) Send(m *PushStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hallPushStreamClient) CloseAndRecv() (*PushStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PushStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hallClient) MustPush(ctx context.Context, in *MustPushRequest, opts ...grpc.CallOption) (*MustPushResponse, error) {
	out := new(MustPushResponse)
	err := c.cc.Invoke(ctx, Hall_MustPush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (Hall_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hall_ServiceDesc.Streams[1], Hall_ClientStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hallClientStreamClient{stream}
	return x, nil
}

type Hall_ClientStreamClient interface {
	Send(*StreamDataRequest) error
	Recv() (*StreamDataResponse, error)
	grpc.ClientStream
}

type hallClientStreamClient struct {
	grpc.ClientStream
}

func (x *hallClientStreamClient) Send(m *StreamDataRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hallClientStreamClient) Recv() (*StreamDataResponse, error) {
	m := new(StreamDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hallClient) GateStream(ctx context.Context, opts ...grpc.CallOption) (Hall_GateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hall_ServiceDesc.Streams[2], Hall_GateStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hallGateStreamClient{stream}
	return x, nil
}

type Hall_GateStreamClient interface {
	Send(*GateDataRequest) error
	Recv() (*GateDataResponse, error)
	grpc.ClientStream
}

type hallGateStreamClient struct {
	grpc.ClientStream
}

func (x *hallGateStreamClient) Send(m *GateDataRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hallGateStreamClient) Recv() (*GateDataResponse, error) {
	m := new(GateDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HallServer is the server API for Hall service.
// All implementations must embed UnimplementedHallServer
// for forward compatibility
type HallServer interface {
	// SayHello 方法
	UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error)
	PushStream(Hall_PushStreamServer) error
	MustPush(context.Context, *MustPushRequest) (*MustPushResponse, error)
	ClientStream(Hall_ClientStreamServer) error
	GateStream(Hall_GateStreamServer) error
	mustEmbedUnimplementedHallServer()
}

// UnimplementedHallServer must be embedded to have forward compatible implementations.
type UnimplementedHallServer struct {
}

func (UnimplementedHallServer) UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInstead not implemented")
}
func (UnimplementedHallServer) PushStream(Hall_PushStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PushStream not implemented")
}
func (UnimplementedHallServer) MustPush(context.Context, *MustPushRequest) (*MustPushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MustPush not implemented")
}
func (UnimplementedHallServer) ClientStream(Hall_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedHallServer) GateStream(Hall_GateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GateStream not implemented")
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

func _Hall_PushStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HallServer).PushStream(&hallPushStreamServer{stream})
}

type Hall_PushStreamServer interface {
	SendAndClose(*PushStreamResponse) error
	Recv() (*PushStreamRequest, error)
	grpc.ServerStream
}

type hallPushStreamServer struct {
	grpc.ServerStream
}

func (x *hallPushStreamServer) SendAndClose(m *PushStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hallPushStreamServer) Recv() (*PushStreamRequest, error) {
	m := new(PushStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hall_MustPush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MustPushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).MustPush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_MustPush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).MustPush(ctx, req.(*MustPushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hall_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HallServer).ClientStream(&hallClientStreamServer{stream})
}

type Hall_ClientStreamServer interface {
	Send(*StreamDataResponse) error
	Recv() (*StreamDataRequest, error)
	grpc.ServerStream
}

type hallClientStreamServer struct {
	grpc.ServerStream
}

func (x *hallClientStreamServer) Send(m *StreamDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hallClientStreamServer) Recv() (*StreamDataRequest, error) {
	m := new(StreamDataRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hall_GateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HallServer).GateStream(&hallGateStreamServer{stream})
}

type Hall_GateStreamServer interface {
	Send(*GateDataResponse) error
	Recv() (*GateDataRequest, error)
	grpc.ServerStream
}

type hallGateStreamServer struct {
	grpc.ServerStream
}

func (x *hallGateStreamServer) Send(m *GateDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hallGateStreamServer) Recv() (*GateDataRequest, error) {
	m := new(GateDataRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
		{
			MethodName: "MustPush",
			Handler:    _Hall_MustPush_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushStream",
			Handler:       _Hall_PushStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _Hall_ClientStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "GateStream",
			Handler:       _Hall_GateStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "doc/grpc/hall.proto",
}

const (
	Gate_HallSuspend_FullMethodName = "/hall_grpc.Gate/HallSuspend"
	Gate_HallRestart_FullMethodName = "/hall_grpc.Gate/HallRestart"
)

// GateClient is the client API for Gate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GateClient interface {
	HallSuspend(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error)
	HallRestart(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error)
}

type gateClient struct {
	cc grpc.ClientConnInterface
}

func NewGateClient(cc grpc.ClientConnInterface) GateClient {
	return &gateClient{cc}
}

func (c *gateClient) HallSuspend(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error) {
	out := new(UserInsteadResponse)
	err := c.cc.Invoke(ctx, Gate_HallSuspend_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateClient) HallRestart(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error) {
	out := new(UserInsteadResponse)
	err := c.cc.Invoke(ctx, Gate_HallRestart_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GateServer is the server API for Gate service.
// All implementations must embed UnimplementedGateServer
// for forward compatibility
type GateServer interface {
	HallSuspend(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error)
	HallRestart(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error)
	mustEmbedUnimplementedGateServer()
}

// UnimplementedGateServer must be embedded to have forward compatible implementations.
type UnimplementedGateServer struct {
}

func (UnimplementedGateServer) HallSuspend(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HallSuspend not implemented")
}
func (UnimplementedGateServer) HallRestart(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HallRestart not implemented")
}
func (UnimplementedGateServer) mustEmbedUnimplementedGateServer() {}

// UnsafeGateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GateServer will
// result in compilation errors.
type UnsafeGateServer interface {
	mustEmbedUnimplementedGateServer()
}

func RegisterGateServer(s grpc.ServiceRegistrar, srv GateServer) {
	s.RegisterService(&Gate_ServiceDesc, srv)
}

func _Gate_HallSuspend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInsteadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).HallSuspend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gate_HallSuspend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).HallSuspend(ctx, req.(*UserInsteadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gate_HallRestart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInsteadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).HallRestart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gate_HallRestart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).HallRestart(ctx, req.(*UserInsteadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gate_ServiceDesc is the grpc.ServiceDesc for Gate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hall_grpc.Gate",
	HandlerType: (*GateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HallSuspend",
			Handler:    _Gate_HallSuspend_Handler,
		},
		{
			MethodName: "HallRestart",
			Handler:    _Gate_HallRestart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "doc/grpc/hall.proto",
}
