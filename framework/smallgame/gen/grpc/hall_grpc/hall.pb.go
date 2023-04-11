// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v3.21.12
// source: doc/grpc/hall.proto

package hall_grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PacketType int32

const (
	PacketType_NONE         PacketType = 0
	PacketType_DATA         PacketType = 1
	PacketType_KICK         PacketType = 5
	PacketType_USER_INSTEAD PacketType = 8
)

// Enum value maps for PacketType.
var (
	PacketType_name = map[int32]string{
		0: "NONE",
		1: "DATA",
		5: "KICK",
		8: "USER_INSTEAD",
	}
	PacketType_value = map[string]int32{
		"NONE":         0,
		"DATA":         1,
		"KICK":         5,
		"USER_INSTEAD": 8,
	}
)

func (x PacketType) Enum() *PacketType {
	p := new(PacketType)
	*p = x
	return p
}

func (x PacketType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketType) Descriptor() protoreflect.EnumDescriptor {
	return file_doc_grpc_hall_proto_enumTypes[0].Descriptor()
}

func (PacketType) Type() protoreflect.EnumType {
	return &file_doc_grpc_hall_proto_enumTypes[0]
}

func (x PacketType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketType.Descriptor instead.
func (PacketType) EnumDescriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{0}
}

type HallDataType int32

const (
	HallDataType_SERVER_SUSPEND HallDataType = 0
	HallDataType_SERVER_REPORT  HallDataType = 1
)

// Enum value maps for HallDataType.
var (
	HallDataType_name = map[int32]string{
		0: "SERVER_SUSPEND",
		1: "SERVER_REPORT",
	}
	HallDataType_value = map[string]int32{
		"SERVER_SUSPEND": 0,
		"SERVER_REPORT":  1,
	}
)

func (x HallDataType) Enum() *HallDataType {
	p := new(HallDataType)
	*p = x
	return p
}

func (x HallDataType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HallDataType) Descriptor() protoreflect.EnumDescriptor {
	return file_doc_grpc_hall_proto_enumTypes[1].Descriptor()
}

func (HallDataType) Type() protoreflect.EnumType {
	return &file_doc_grpc_hall_proto_enumTypes[1]
}

func (x HallDataType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HallDataType.Descriptor instead.
func (HallDataType) EnumDescriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{1}
}

type GateDataPush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateDataPush) Reset() {
	*x = GateDataPush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateDataPush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateDataPush) ProtoMessage() {}

func (x *GateDataPush) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateDataPush.ProtoReflect.Descriptor instead.
func (*GateDataPush) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{0}
}

// stream返回结构
type HallDataPush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   HallDataType    `protobuf:"varint,1,opt,name=Type,proto3,enum=hall_grpc.HallDataType" json:"Type,omitempty"`
	Report *HallReportPush `protobuf:"bytes,2,opt,name=Report,proto3" json:"Report,omitempty"`
}

func (x *HallDataPush) Reset() {
	*x = HallDataPush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HallDataPush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HallDataPush) ProtoMessage() {}

func (x *HallDataPush) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HallDataPush.ProtoReflect.Descriptor instead.
func (*HallDataPush) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{1}
}

func (x *HallDataPush) GetType() HallDataType {
	if x != nil {
		return x.Type
	}
	return HallDataType_SERVER_SUSPEND
}

func (x *HallDataPush) GetReport() *HallReportPush {
	if x != nil {
		return x.Report
	}
	return nil
}

type HallReportPush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerCount int64 `protobuf:"varint,1,opt,name=PlayerCount,proto3" json:"PlayerCount,omitempty"`
}

func (x *HallReportPush) Reset() {
	*x = HallReportPush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HallReportPush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HallReportPush) ProtoMessage() {}

func (x *HallReportPush) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HallReportPush.ProtoReflect.Descriptor instead.
func (*HallReportPush) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{2}
}

func (x *HallReportPush) GetPlayerCount() int64 {
	if x != nil {
		return x.PlayerCount
	}
	return 0
}

type StreamDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MemberId  string `protobuf:"bytes,1,opt,name=MemberId,proto3" json:"MemberId,omitempty"`
	SessionId uint64 `protobuf:"varint,2,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	ReqId     uint64 `protobuf:"varint,3,opt,name=ReqId,proto3" json:"ReqId,omitempty"`
	Data      []byte `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *StreamDataRequest) Reset() {
	*x = StreamDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamDataRequest) ProtoMessage() {}

func (x *StreamDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamDataRequest.ProtoReflect.Descriptor instead.
func (*StreamDataRequest) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{3}
}

func (x *StreamDataRequest) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *StreamDataRequest) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *StreamDataRequest) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *StreamDataRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// stream返回结构
type StreamDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      PacketType `protobuf:"varint,1,opt,name=Type,proto3,enum=hall_grpc.PacketType" json:"Type,omitempty"`
	SessionId uint64     `protobuf:"varint,2,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	ReqId     uint64     `protobuf:"varint,3,opt,name=ReqId,proto3" json:"ReqId,omitempty"`
	Data      []byte     `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
	Route     string     `protobuf:"bytes,5,opt,name=Route,proto3" json:"Route,omitempty"`
}

func (x *StreamDataResponse) Reset() {
	*x = StreamDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamDataResponse) ProtoMessage() {}

func (x *StreamDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamDataResponse.ProtoReflect.Descriptor instead.
func (*StreamDataResponse) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{4}
}

func (x *StreamDataResponse) GetType() PacketType {
	if x != nil {
		return x.Type
	}
	return PacketType_NONE
}

func (x *StreamDataResponse) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *StreamDataResponse) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *StreamDataResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *StreamDataResponse) GetRoute() string {
	if x != nil {
		return x.Route
	}
	return ""
}

// 请求消息
type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{5}
}

func (x *HelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// 响应消息
type HelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Replay string `protobuf:"bytes,1,opt,name=replay,proto3" json:"replay,omitempty"`
}

func (x *HelloResponse) Reset() {
	*x = HelloResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResponse) ProtoMessage() {}

func (x *HelloResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResponse.ProtoReflect.Descriptor instead.
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{6}
}

func (x *HelloResponse) GetReplay() string {
	if x != nil {
		return x.Replay
	}
	return ""
}

// 顶号下线
type UserInsteadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *UserInsteadRequest) Reset() {
	*x = UserInsteadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInsteadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInsteadRequest) ProtoMessage() {}

func (x *UserInsteadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInsteadRequest.ProtoReflect.Descriptor instead.
func (*UserInsteadRequest) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{7}
}

func (x *UserInsteadRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserInsteadRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

// 顶号下线
type UserInsteadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
}

func (x *UserInsteadResponse) Reset() {
	*x = UserInsteadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInsteadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInsteadResponse) ProtoMessage() {}

func (x *UserInsteadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInsteadResponse.ProtoReflect.Descriptor instead.
func (*UserInsteadResponse) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{8}
}

func (x *UserInsteadResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *UserInsteadResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

// 推送消息
type PushStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Data   []byte `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *PushStreamRequest) Reset() {
	*x = PushStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushStreamRequest) ProtoMessage() {}

func (x *PushStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushStreamRequest.ProtoReflect.Descriptor instead.
func (*PushStreamRequest) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{9}
}

func (x *PushStreamRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *PushStreamRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// 推送消息
type PushStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
}

func (x *PushStreamResponse) Reset() {
	*x = PushStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushStreamResponse) ProtoMessage() {}

func (x *PushStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushStreamResponse.ProtoReflect.Descriptor instead.
func (*PushStreamResponse) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{10}
}

func (x *PushStreamResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *PushStreamResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

// 推送消息
type MustPushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Data   []byte `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *MustPushRequest) Reset() {
	*x = MustPushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MustPushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MustPushRequest) ProtoMessage() {}

func (x *MustPushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MustPushRequest.ProtoReflect.Descriptor instead.
func (*MustPushRequest) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{11}
}

func (x *MustPushRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MustPushRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// 推送消息
type MustPushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
}

func (x *MustPushResponse) Reset() {
	*x = MustPushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doc_grpc_hall_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MustPushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MustPushResponse) ProtoMessage() {}

func (x *MustPushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doc_grpc_hall_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MustPushResponse.ProtoReflect.Descriptor instead.
func (*MustPushResponse) Descriptor() ([]byte, []int) {
	return file_doc_grpc_hall_proto_rawDescGZIP(), []int{12}
}

func (x *MustPushResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *MustPushResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

var File_doc_grpc_hall_proto protoreflect.FileDescriptor

var file_doc_grpc_hall_proto_rawDesc = []byte{
	0x0a, 0x13, 0x64, 0x6f, 0x63, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x61, 0x6c, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63,
	0x22, 0x0e, 0x0a, 0x0c, 0x47, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x50, 0x75, 0x73, 0x68,
	0x22, 0x6e, 0x0a, 0x0c, 0x48, 0x61, 0x6c, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x50, 0x75, 0x73, 0x68,
	0x12, 0x2b, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x61, 0x6c, 0x6c, 0x44,
	0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x31, 0x0a,
	0x06, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x61, 0x6c, 0x6c, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x06, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x22, 0x32, 0x0a, 0x0e, 0x48, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x75,
	0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x77, 0x0a, 0x11, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x71, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x52, 0x65, 0x71, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x9d, 0x01,
	0x0a, 0x12, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x15, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x52, 0x65, 0x71, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x52, 0x65,
	0x71, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x22, 0x22, 0x0a,
	0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x27, 0x0a, 0x0d, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x22, 0x46, 0x0a, 0x12, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x22, 0x4f, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x4d, 0x73, 0x67, 0x22, 0x3f, 0x0a, 0x11, 0x50, 0x75, 0x73, 0x68, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x22, 0x4e, 0x0a, 0x12, 0x50, 0x75, 0x73, 0x68, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x73, 0x67, 0x22, 0x3d, 0x0a, 0x0f, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44,
	0x61, 0x74, 0x61, 0x22, 0x4c, 0x0a, 0x10, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73,
	0x67, 0x2a, 0x3c, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54,
	0x41, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4b, 0x49, 0x43, 0x4b, 0x10, 0x05, 0x12, 0x10, 0x0a,
	0x0c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x45, 0x41, 0x44, 0x10, 0x08, 0x2a,
	0x35, 0x0a, 0x0c, 0x48, 0x61, 0x6c, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x0e, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e,
	0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x52, 0x45,
	0x50, 0x4f, 0x52, 0x54, 0x10, 0x01, 0x32, 0x85, 0x03, 0x0a, 0x04, 0x48, 0x61, 0x6c, 0x6c, 0x12,
	0x4e, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61, 0x64, 0x12, 0x1d,
	0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x73, 0x74, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e,
	0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x73, 0x74, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4d, 0x0a, 0x0a, 0x50, 0x75, 0x73, 0x68, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e,
	0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x68, 0x61,
	0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x45,
	0x0a, 0x08, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x12, 0x1a, 0x2e, 0x68, 0x61, 0x6c,
	0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x44, 0x0a, 0x0a, 0x47, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x50, 0x75, 0x73, 0x68, 0x1a,
	0x17, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x61, 0x6c, 0x6c,
	0x44, 0x61, 0x74, 0x61, 0x50, 0x75, 0x73, 0x68, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x14,
	0x5a, 0x12, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x61, 0x6c, 0x6c, 0x5f,
	0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_doc_grpc_hall_proto_rawDescOnce sync.Once
	file_doc_grpc_hall_proto_rawDescData = file_doc_grpc_hall_proto_rawDesc
)

func file_doc_grpc_hall_proto_rawDescGZIP() []byte {
	file_doc_grpc_hall_proto_rawDescOnce.Do(func() {
		file_doc_grpc_hall_proto_rawDescData = protoimpl.X.CompressGZIP(file_doc_grpc_hall_proto_rawDescData)
	})
	return file_doc_grpc_hall_proto_rawDescData
}

var file_doc_grpc_hall_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_doc_grpc_hall_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_doc_grpc_hall_proto_goTypes = []interface{}{
	(PacketType)(0),             // 0: hall_grpc.PacketType
	(HallDataType)(0),           // 1: hall_grpc.HallDataType
	(*GateDataPush)(nil),        // 2: hall_grpc.GateDataPush
	(*HallDataPush)(nil),        // 3: hall_grpc.HallDataPush
	(*HallReportPush)(nil),      // 4: hall_grpc.HallReportPush
	(*StreamDataRequest)(nil),   // 5: hall_grpc.StreamDataRequest
	(*StreamDataResponse)(nil),  // 6: hall_grpc.StreamDataResponse
	(*HelloRequest)(nil),        // 7: hall_grpc.HelloRequest
	(*HelloResponse)(nil),       // 8: hall_grpc.HelloResponse
	(*UserInsteadRequest)(nil),  // 9: hall_grpc.UserInsteadRequest
	(*UserInsteadResponse)(nil), // 10: hall_grpc.UserInsteadResponse
	(*PushStreamRequest)(nil),   // 11: hall_grpc.PushStreamRequest
	(*PushStreamResponse)(nil),  // 12: hall_grpc.PushStreamResponse
	(*MustPushRequest)(nil),     // 13: hall_grpc.MustPushRequest
	(*MustPushResponse)(nil),    // 14: hall_grpc.MustPushResponse
}
var file_doc_grpc_hall_proto_depIdxs = []int32{
	1,  // 0: hall_grpc.HallDataPush.Type:type_name -> hall_grpc.HallDataType
	4,  // 1: hall_grpc.HallDataPush.Report:type_name -> hall_grpc.HallReportPush
	0,  // 2: hall_grpc.StreamDataResponse.Type:type_name -> hall_grpc.PacketType
	9,  // 3: hall_grpc.Hall.UserInstead:input_type -> hall_grpc.UserInsteadRequest
	11, // 4: hall_grpc.Hall.PushStream:input_type -> hall_grpc.PushStreamRequest
	13, // 5: hall_grpc.Hall.MustPush:input_type -> hall_grpc.MustPushRequest
	5,  // 6: hall_grpc.Hall.ClientStream:input_type -> hall_grpc.StreamDataRequest
	2,  // 7: hall_grpc.Hall.GateStream:input_type -> hall_grpc.GateDataPush
	10, // 8: hall_grpc.Hall.UserInstead:output_type -> hall_grpc.UserInsteadResponse
	12, // 9: hall_grpc.Hall.PushStream:output_type -> hall_grpc.PushStreamResponse
	14, // 10: hall_grpc.Hall.MustPush:output_type -> hall_grpc.MustPushResponse
	6,  // 11: hall_grpc.Hall.ClientStream:output_type -> hall_grpc.StreamDataResponse
	3,  // 12: hall_grpc.Hall.GateStream:output_type -> hall_grpc.HallDataPush
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_doc_grpc_hall_proto_init() }
func file_doc_grpc_hall_proto_init() {
	if File_doc_grpc_hall_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_doc_grpc_hall_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateDataPush); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HallDataPush); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HallReportPush); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamDataRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamDataResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInsteadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInsteadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushStreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MustPushRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_doc_grpc_hall_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MustPushResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_doc_grpc_hall_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_doc_grpc_hall_proto_goTypes,
		DependencyIndexes: file_doc_grpc_hall_proto_depIdxs,
		EnumInfos:         file_doc_grpc_hall_proto_enumTypes,
		MessageInfos:      file_doc_grpc_hall_proto_msgTypes,
	}.Build()
	File_doc_grpc_hall_proto = out.File
	file_doc_grpc_hall_proto_rawDesc = nil
	file_doc_grpc_hall_proto_goTypes = nil
	file_doc_grpc_hall_proto_depIdxs = nil
}
