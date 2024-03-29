// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: doc/service/hall.proto

package hallpb

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
	Hall_ClientStream_FullMethodName = "/hallpb.Hall/ClientStream"
	Hall_GateStream_FullMethodName   = "/hallpb.Hall/GateStream"
	Hall_Info_FullMethodName         = "/hallpb.Hall/Info"
	Hall_HealthCheck_FullMethodName  = "/hallpb.Hall/HealthCheck"
	Hall_MustPush_FullMethodName     = "/hallpb.Hall/MustPush"
	Hall_SendMessage_FullMethodName  = "/hallpb.Hall/SendMessage"
	Hall_CallMessage_FullMethodName  = "/hallpb.Hall/CallMessage"
	Hall_UserInstead_FullMethodName  = "/hallpb.Hall/UserInstead"
	Hall_Kick_FullMethodName         = "/hallpb.Hall/Kick"
)

// HallClient is the client API for Hall service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HallClient interface {
	// client消息流
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (Hall_ClientStreamClient, error)
	// 网关消息流
	GateStream(ctx context.Context, opts ...grpc.CallOption) (Hall_GateStreamClient, error)
	// 状态
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	// 心跳
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	// rpc PushStream (stream PushStreamNotify) returns (PushStreamPush) {}
	MustPush(ctx context.Context, in *MustPushRequest, opts ...grpc.CallOption) (*MustPushResponse, error)
	// 发送消息
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	// 发送消息
	CallMessage(ctx context.Context, in *CallMessageRequest, opts ...grpc.CallOption) (*CallMessageResponse, error)
	// 顶号下线
	UserInstead(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error)
	// 踢人下线
	Kick(ctx context.Context, in *KickRequest, opts ...grpc.CallOption) (*KickResponse, error)
}

type hallClient struct {
	cc grpc.ClientConnInterface
}

func NewHallClient(cc grpc.ClientConnInterface) HallClient {
	return &hallClient{cc}
}

func (c *hallClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (Hall_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hall_ServiceDesc.Streams[0], Hall_ClientStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hallClientStreamClient{stream}
	return x, nil
}

type Hall_ClientStreamClient interface {
	Send(*ClientMessageRequest) error
	Recv() (*ClientMessageResponse, error)
	grpc.ClientStream
}

type hallClientStreamClient struct {
	grpc.ClientStream
}

func (x *hallClientStreamClient) Send(m *ClientMessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hallClientStreamClient) Recv() (*ClientMessageResponse, error) {
	m := new(ClientMessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hallClient) GateStream(ctx context.Context, opts ...grpc.CallOption) (Hall_GateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hall_ServiceDesc.Streams[1], Hall_GateStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hallGateStreamClient{stream}
	return x, nil
}

type Hall_GateStreamClient interface {
	Send(*GateStreamRequest) error
	Recv() (*GateStreamResponse, error)
	grpc.ClientStream
}

type hallGateStreamClient struct {
	grpc.ClientStream
}

func (x *hallGateStreamClient) Send(m *GateStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hallGateStreamClient) Recv() (*GateStreamResponse, error) {
	m := new(GateStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hallClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, Hall_Info_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, Hall_HealthCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) MustPush(ctx context.Context, in *MustPushRequest, opts ...grpc.CallOption) (*MustPushResponse, error) {
	out := new(MustPushResponse)
	err := c.cc.Invoke(ctx, Hall_MustPush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, Hall_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) CallMessage(ctx context.Context, in *CallMessageRequest, opts ...grpc.CallOption) (*CallMessageResponse, error) {
	out := new(CallMessageResponse)
	err := c.cc.Invoke(ctx, Hall_CallMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) UserInstead(ctx context.Context, in *UserInsteadRequest, opts ...grpc.CallOption) (*UserInsteadResponse, error) {
	out := new(UserInsteadResponse)
	err := c.cc.Invoke(ctx, Hall_UserInstead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hallClient) Kick(ctx context.Context, in *KickRequest, opts ...grpc.CallOption) (*KickResponse, error) {
	out := new(KickResponse)
	err := c.cc.Invoke(ctx, Hall_Kick_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HallServer is the server API for Hall service.
// All implementations must embed UnimplementedHallServer
// for forward compatibility
type HallServer interface {
	// client消息流
	ClientStream(Hall_ClientStreamServer) error
	// 网关消息流
	GateStream(Hall_GateStreamServer) error
	// 状态
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	// 心跳
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	// rpc PushStream (stream PushStreamNotify) returns (PushStreamPush) {}
	MustPush(context.Context, *MustPushRequest) (*MustPushResponse, error)
	// 发送消息
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	// 发送消息
	CallMessage(context.Context, *CallMessageRequest) (*CallMessageResponse, error)
	// 顶号下线
	UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error)
	// 踢人下线
	Kick(context.Context, *KickRequest) (*KickResponse, error)
	mustEmbedUnimplementedHallServer()
}

// UnimplementedHallServer must be embedded to have forward compatible implementations.
type UnimplementedHallServer struct {
}

func (UnimplementedHallServer) ClientStream(Hall_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedHallServer) GateStream(Hall_GateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GateStream not implemented")
}
func (UnimplementedHallServer) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedHallServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedHallServer) MustPush(context.Context, *MustPushRequest) (*MustPushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MustPush not implemented")
}
func (UnimplementedHallServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedHallServer) CallMessage(context.Context, *CallMessageRequest) (*CallMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallMessage not implemented")
}
func (UnimplementedHallServer) UserInstead(context.Context, *UserInsteadRequest) (*UserInsteadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInstead not implemented")
}
func (UnimplementedHallServer) Kick(context.Context, *KickRequest) (*KickResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Kick not implemented")
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

func _Hall_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HallServer).ClientStream(&hallClientStreamServer{stream})
}

type Hall_ClientStreamServer interface {
	Send(*ClientMessageResponse) error
	Recv() (*ClientMessageRequest, error)
	grpc.ServerStream
}

type hallClientStreamServer struct {
	grpc.ServerStream
}

func (x *hallClientStreamServer) Send(m *ClientMessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hallClientStreamServer) Recv() (*ClientMessageRequest, error) {
	m := new(ClientMessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hall_GateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HallServer).GateStream(&hallGateStreamServer{stream})
}

type Hall_GateStreamServer interface {
	Send(*GateStreamResponse) error
	Recv() (*GateStreamRequest, error)
	grpc.ServerStream
}

type hallGateStreamServer struct {
	grpc.ServerStream
}

func (x *hallGateStreamServer) Send(m *GateStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hallGateStreamServer) Recv() (*GateStreamRequest, error) {
	m := new(GateStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hall_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_Info_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hall_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _Hall_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hall_CallMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).CallMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_CallMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).CallMessage(ctx, req.(*CallMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _Hall_Kick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KickRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HallServer).Kick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Hall_Kick_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HallServer).Kick(ctx, req.(*KickRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Hall_ServiceDesc is the grpc.ServiceDesc for Hall service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hall_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hallpb.Hall",
	HandlerType: (*HallServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _Hall_Info_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _Hall_HealthCheck_Handler,
		},
		{
			MethodName: "MustPush",
			Handler:    _Hall_MustPush_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _Hall_SendMessage_Handler,
		},
		{
			MethodName: "CallMessage",
			Handler:    _Hall_CallMessage_Handler,
		},
		{
			MethodName: "UserInstead",
			Handler:    _Hall_UserInstead_Handler,
		},
		{
			MethodName: "Kick",
			Handler:    _Hall_Kick_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
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
	Metadata: "doc/service/hall.proto",
}
