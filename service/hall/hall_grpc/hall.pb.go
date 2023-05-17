// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: service/hall/hall.proto

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
	return file_service_hall_hall_proto_enumTypes[0].Descriptor()
}

func (PacketType) Type() protoreflect.EnumType {
	return &file_service_hall_hall_proto_enumTypes[0]
}

func (x PacketType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketType.Descriptor instead.
func (PacketType) EnumDescriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{0}
}

type GateStreamNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateStreamNotify) Reset() {
	*x = GateStreamNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateStreamNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateStreamNotify) ProtoMessage() {}

func (x *GateStreamNotify) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateStreamNotify.ProtoReflect.Descriptor instead.
func (*GateStreamNotify) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{0}
}

type GateStreamPush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateStreamPush) Reset() {
	*x = GateStreamPush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateStreamPush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateStreamPush) ProtoMessage() {}

func (x *GateStreamPush) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateStreamPush.ProtoReflect.Descriptor instead.
func (*GateStreamPush) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{1}
}

type HeartbeatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HeartbeatRequest) Reset() {
	*x = HeartbeatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatRequest) ProtoMessage() {}

func (x *HeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatRequest.ProtoReflect.Descriptor instead.
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{2}
}

type HeartbeatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerCount int64 `protobuf:"varint,1,opt,name=PlayerCount,proto3" json:"PlayerCount,omitempty"`
}

func (x *HeartbeatResponse) Reset() {
	*x = HeartbeatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatResponse) ProtoMessage() {}

func (x *HeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatResponse.ProtoReflect.Descriptor instead.
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{3}
}

func (x *HeartbeatResponse) GetPlayerCount() int64 {
	if x != nil {
		return x.PlayerCount
	}
	return 0
}

type InfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InfoRequest) Reset() {
	*x = InfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRequest) ProtoMessage() {}

func (x *InfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRequest.ProtoReflect.Descriptor instead.
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{4}
}

type InfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BuildTime    int64  `protobuf:"varint,1,opt,name=BuildTime,proto3" json:"BuildTime,omitempty"`
	BuildVersion string `protobuf:"bytes,2,opt,name=BuildVersion,proto3" json:"BuildVersion,omitempty"`
	// 在线人数
	SessionCount int64 `protobuf:"varint,3,opt,name=SessionCount,proto3" json:"SessionCount,omitempty"`
}

func (x *InfoResponse) Reset() {
	*x = InfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoResponse) ProtoMessage() {}

func (x *InfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoResponse.ProtoReflect.Descriptor instead.
func (*InfoResponse) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{5}
}

func (x *InfoResponse) GetBuildTime() int64 {
	if x != nil {
		return x.BuildTime
	}
	return 0
}

func (x *InfoResponse) GetBuildVersion() string {
	if x != nil {
		return x.BuildVersion
	}
	return ""
}

func (x *InfoResponse) GetSessionCount() int64 {
	if x != nil {
		return x.SessionCount
	}
	return 0
}

type ClientMessageNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MemberId  string `protobuf:"bytes,1,opt,name=MemberId,proto3" json:"MemberId,omitempty"`
	SessionId uint64 `protobuf:"varint,2,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	ReqId     uint64 `protobuf:"varint,3,opt,name=ReqId,proto3" json:"ReqId,omitempty"`
	Data      []byte `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *ClientMessageNotify) Reset() {
	*x = ClientMessageNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessageNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessageNotify) ProtoMessage() {}

func (x *ClientMessageNotify) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessageNotify.ProtoReflect.Descriptor instead.
func (*ClientMessageNotify) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{6}
}

func (x *ClientMessageNotify) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *ClientMessageNotify) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *ClientMessageNotify) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *ClientMessageNotify) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ClientMessagePush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      PacketType `protobuf:"varint,1,opt,name=Type,proto3,enum=hall_grpc.PacketType" json:"Type,omitempty"`
	SessionId uint64     `protobuf:"varint,2,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	ReqId     uint64     `protobuf:"varint,3,opt,name=ReqId,proto3" json:"ReqId,omitempty"`
	Data      []byte     `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
	Route     string     `protobuf:"bytes,5,opt,name=Route,proto3" json:"Route,omitempty"`
}

func (x *ClientMessagePush) Reset() {
	*x = ClientMessagePush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessagePush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessagePush) ProtoMessage() {}

func (x *ClientMessagePush) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessagePush.ProtoReflect.Descriptor instead.
func (*ClientMessagePush) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{7}
}

func (x *ClientMessagePush) GetType() PacketType {
	if x != nil {
		return x.Type
	}
	return PacketType_NONE
}

func (x *ClientMessagePush) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *ClientMessagePush) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *ClientMessagePush) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ClientMessagePush) GetRoute() string {
	if x != nil {
		return x.Route
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
		mi := &file_service_hall_hall_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInsteadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInsteadRequest) ProtoMessage() {}

func (x *UserInsteadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[8]
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
	return file_service_hall_hall_proto_rawDescGZIP(), []int{8}
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
		mi := &file_service_hall_hall_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInsteadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInsteadResponse) ProtoMessage() {}

func (x *UserInsteadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[9]
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
	return file_service_hall_hall_proto_rawDescGZIP(), []int{9}
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
		mi := &file_service_hall_hall_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MustPushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MustPushRequest) ProtoMessage() {}

func (x *MustPushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[10]
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
	return file_service_hall_hall_proto_rawDescGZIP(), []int{10}
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
		mi := &file_service_hall_hall_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MustPushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MustPushResponse) ProtoMessage() {}

func (x *MustPushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[11]
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
	return file_service_hall_hall_proto_rawDescGZIP(), []int{11}
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

type KickRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Reason string `protobuf:"bytes,2,opt,name=Reason,proto3" json:"Reason,omitempty"`
}

func (x *KickRequest) Reset() {
	*x = KickRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KickRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KickRequest) ProtoMessage() {}

func (x *KickRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KickRequest.ProtoReflect.Descriptor instead.
func (*KickRequest) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{12}
}

func (x *KickRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *KickRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type KickResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode int32  `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	ErrorMsg  string `protobuf:"bytes,2,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
}

func (x *KickResponse) Reset() {
	*x = KickResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_hall_hall_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KickResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KickResponse) ProtoMessage() {}

func (x *KickResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_hall_hall_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KickResponse.ProtoReflect.Descriptor instead.
func (*KickResponse) Descriptor() ([]byte, []int) {
	return file_service_hall_hall_proto_rawDescGZIP(), []int{13}
}

func (x *KickResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *KickResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

var File_service_hall_hall_proto protoreflect.FileDescriptor

var file_service_hall_hall_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x68, 0x61, 0x6c, 0x6c, 0x2f, 0x68,
	0x61, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x68, 0x61, 0x6c, 0x6c, 0x5f,
	0x67, 0x72, 0x70, 0x63, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x75, 0x73, 0x68, 0x22, 0x12, 0x0a, 0x10, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x35,
	0x0a, 0x11, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x74, 0x0a, 0x0c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x79, 0x0a, 0x13, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x52,
	0x65, 0x71, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x52, 0x65, 0x71, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x9c, 0x01, 0x0a, 0x11, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x75, 0x73, 0x68, 0x12, 0x29, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x68, 0x61, 0x6c, 0x6c,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x71, 0x49, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x52, 0x65, 0x71, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61,
	0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14,
	0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x22, 0x46, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x4f, 0x0a, 0x13,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x22, 0x3d, 0x0a,
	0x0f, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x4c, 0x0a, 0x10,
	0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x22, 0x3d, 0x0a, 0x0b, 0x4b, 0x69,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x48, 0x0a, 0x0c, 0x4b, 0x69, 0x63,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x4d, 0x73, 0x67, 0x2a, 0x3c, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x44,
	0x41, 0x54, 0x41, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4b, 0x49, 0x43, 0x4b, 0x10, 0x05, 0x12,
	0x10, 0x0a, 0x0c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x45, 0x41, 0x44, 0x10,
	0x08, 0x32, 0xfd, 0x03, 0x0a, 0x04, 0x48, 0x61, 0x6c, 0x6c, 0x12, 0x4e, 0x0a, 0x0b, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61, 0x64, 0x12, 0x1d, 0x2e, 0x68, 0x61, 0x6c, 0x6c,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x08, 0x4d, 0x75,
	0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x12, 0x1a, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x4d, 0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d,
	0x75, 0x73, 0x74, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x52, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x1e, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x1a, 0x1c, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x75, 0x73, 0x68, 0x22,
	0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x4a, 0x0a, 0x0a, 0x47, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x1b, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x47, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x1a, 0x19, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x75, 0x73, 0x68, 0x22, 0x00, 0x28, 0x01, 0x30,
	0x01, 0x12, 0x39, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x68, 0x61, 0x6c, 0x6c,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x09,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1b, 0x2e, 0x68, 0x61, 0x6c, 0x6c,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x04, 0x4b, 0x69, 0x63, 0x6b, 0x12, 0x16,
	0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4b, 0x69, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x4b, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x68, 0x61, 0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_hall_hall_proto_rawDescOnce sync.Once
	file_service_hall_hall_proto_rawDescData = file_service_hall_hall_proto_rawDesc
)

func file_service_hall_hall_proto_rawDescGZIP() []byte {
	file_service_hall_hall_proto_rawDescOnce.Do(func() {
		file_service_hall_hall_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_hall_hall_proto_rawDescData)
	})
	return file_service_hall_hall_proto_rawDescData
}

var file_service_hall_hall_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_service_hall_hall_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_service_hall_hall_proto_goTypes = []interface{}{
	(PacketType)(0),             // 0: hall_grpc.PacketType
	(*GateStreamNotify)(nil),    // 1: hall_grpc.GateStreamNotify
	(*GateStreamPush)(nil),      // 2: hall_grpc.GateStreamPush
	(*HeartbeatRequest)(nil),    // 3: hall_grpc.HeartbeatRequest
	(*HeartbeatResponse)(nil),   // 4: hall_grpc.HeartbeatResponse
	(*InfoRequest)(nil),         // 5: hall_grpc.InfoRequest
	(*InfoResponse)(nil),        // 6: hall_grpc.InfoResponse
	(*ClientMessageNotify)(nil), // 7: hall_grpc.ClientMessageNotify
	(*ClientMessagePush)(nil),   // 8: hall_grpc.ClientMessagePush
	(*UserInsteadRequest)(nil),  // 9: hall_grpc.UserInsteadRequest
	(*UserInsteadResponse)(nil), // 10: hall_grpc.UserInsteadResponse
	(*MustPushRequest)(nil),     // 11: hall_grpc.MustPushRequest
	(*MustPushResponse)(nil),    // 12: hall_grpc.MustPushResponse
	(*KickRequest)(nil),         // 13: hall_grpc.KickRequest
	(*KickResponse)(nil),        // 14: hall_grpc.KickResponse
}
var file_service_hall_hall_proto_depIdxs = []int32{
	0,  // 0: hall_grpc.ClientMessagePush.Type:type_name -> hall_grpc.PacketType
	9,  // 1: hall_grpc.Hall.UserInstead:input_type -> hall_grpc.UserInsteadRequest
	11, // 2: hall_grpc.Hall.MustPush:input_type -> hall_grpc.MustPushRequest
	7,  // 3: hall_grpc.Hall.ClientStream:input_type -> hall_grpc.ClientMessageNotify
	1,  // 4: hall_grpc.Hall.GateStream:input_type -> hall_grpc.GateStreamNotify
	5,  // 5: hall_grpc.Hall.Info:input_type -> hall_grpc.InfoRequest
	3,  // 6: hall_grpc.Hall.Heartbeat:input_type -> hall_grpc.HeartbeatRequest
	13, // 7: hall_grpc.Hall.Kick:input_type -> hall_grpc.KickRequest
	10, // 8: hall_grpc.Hall.UserInstead:output_type -> hall_grpc.UserInsteadResponse
	12, // 9: hall_grpc.Hall.MustPush:output_type -> hall_grpc.MustPushResponse
	8,  // 10: hall_grpc.Hall.ClientStream:output_type -> hall_grpc.ClientMessagePush
	2,  // 11: hall_grpc.Hall.GateStream:output_type -> hall_grpc.GateStreamPush
	6,  // 12: hall_grpc.Hall.Info:output_type -> hall_grpc.InfoResponse
	4,  // 13: hall_grpc.Hall.Heartbeat:output_type -> hall_grpc.HeartbeatResponse
	14, // 14: hall_grpc.Hall.Kick:output_type -> hall_grpc.KickResponse
	8,  // [8:15] is the sub-list for method output_type
	1,  // [1:8] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_service_hall_hall_proto_init() }
func file_service_hall_hall_proto_init() {
	if File_service_hall_hall_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_hall_hall_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateStreamNotify); i {
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
		file_service_hall_hall_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateStreamPush); i {
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
		file_service_hall_hall_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatRequest); i {
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
		file_service_hall_hall_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatResponse); i {
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
		file_service_hall_hall_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRequest); i {
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
		file_service_hall_hall_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoResponse); i {
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
		file_service_hall_hall_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessageNotify); i {
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
		file_service_hall_hall_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessagePush); i {
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
		file_service_hall_hall_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_hall_hall_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_hall_hall_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_hall_hall_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_hall_hall_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KickRequest); i {
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
		file_service_hall_hall_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KickResponse); i {
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
			RawDescriptor: file_service_hall_hall_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_hall_hall_proto_goTypes,
		DependencyIndexes: file_service_hall_hall_proto_depIdxs,
		EnumInfos:         file_service_hall_hall_proto_enumTypes,
		MessageInfos:      file_service_hall_hall_proto_msgTypes,
	}.Build()
	File_service_hall_hall_proto = out.File
	file_service_hall_hall_proto_rawDesc = nil
	file_service_hall_hall_proto_goTypes = nil
	file_service_hall_hall_proto_depIdxs = nil
}
