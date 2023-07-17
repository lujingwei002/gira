// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: service/peer/peer.proto

package peerpb

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

type HealthCheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthCheckRequest) Reset() {
	*x = HealthCheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_peer_peer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckRequest) ProtoMessage() {}

func (x *HealthCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_peer_peer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckRequest.ProtoReflect.Descriptor instead.
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return file_service_peer_peer_proto_rawDescGZIP(), []int{0}
}

type HealthCheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BuildTime     int64  `protobuf:"varint,1,opt,name=BuildTime,proto3" json:"BuildTime,omitempty"`
	AppVersion    string `protobuf:"bytes,2,opt,name=AppVersion,proto3" json:"AppVersion,omitempty"`
	UpTime        int64  `protobuf:"varint,3,opt,name=UpTime,proto3" json:"UpTime,omitempty"`
	ResVersion    string `protobuf:"bytes,5,opt,name=ResVersion,proto3" json:"ResVersion,omitempty"`
	LoaderVersion string `protobuf:"bytes,4,opt,name=LoaderVersion,proto3" json:"LoaderVersion,omitempty"`
}

func (x *HealthCheckResponse) Reset() {
	*x = HealthCheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_peer_peer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponse) ProtoMessage() {}

func (x *HealthCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_peer_peer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckResponse.ProtoReflect.Descriptor instead.
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return file_service_peer_peer_proto_rawDescGZIP(), []int{1}
}

func (x *HealthCheckResponse) GetBuildTime() int64 {
	if x != nil {
		return x.BuildTime
	}
	return 0
}

func (x *HealthCheckResponse) GetAppVersion() string {
	if x != nil {
		return x.AppVersion
	}
	return ""
}

func (x *HealthCheckResponse) GetUpTime() int64 {
	if x != nil {
		return x.UpTime
	}
	return 0
}

func (x *HealthCheckResponse) GetResVersion() string {
	if x != nil {
		return x.ResVersion
	}
	return ""
}

func (x *HealthCheckResponse) GetLoaderVersion() string {
	if x != nil {
		return x.LoaderVersion
	}
	return ""
}

type MemStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MemStatsRequest) Reset() {
	*x = MemStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_peer_peer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemStatsRequest) ProtoMessage() {}

func (x *MemStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_peer_peer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemStatsRequest.ProtoReflect.Descriptor instead.
func (*MemStatsRequest) Descriptor() ([]byte, []int) {
	return file_service_peer_peer_proto_rawDescGZIP(), []int{2}
}

type MemStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alloc        uint64 `protobuf:"varint,1,opt,name=Alloc,proto3" json:"Alloc,omitempty"`
	TotalAlloc   uint64 `protobuf:"varint,2,opt,name=TotalAlloc,proto3" json:"TotalAlloc,omitempty"`
	Sys          uint64 `protobuf:"varint,3,opt,name=Sys,proto3" json:"Sys,omitempty"`
	Mallocs      uint64 `protobuf:"varint,4,opt,name=Mallocs,proto3" json:"Mallocs,omitempty"`
	Frees        uint64 `protobuf:"varint,5,opt,name=Frees,proto3" json:"Frees,omitempty"`
	HeapAlloc    uint64 `protobuf:"varint,6,opt,name=HeapAlloc,proto3" json:"HeapAlloc,omitempty"`
	HeapSys      uint64 `protobuf:"varint,7,opt,name=HeapSys,proto3" json:"HeapSys,omitempty"`
	HeapIdle     uint64 `protobuf:"varint,8,opt,name=HeapIdle,proto3" json:"HeapIdle,omitempty"`
	HeapInuse    uint64 `protobuf:"varint,9,opt,name=HeapInuse,proto3" json:"HeapInuse,omitempty"`
	HeapReleased uint64 `protobuf:"varint,10,opt,name=HeapReleased,proto3" json:"HeapReleased,omitempty"`
	HeapObjects  uint64 `protobuf:"varint,11,opt,name=HeapObjects,proto3" json:"HeapObjects,omitempty"`
}

func (x *MemStatsResponse) Reset() {
	*x = MemStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_peer_peer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemStatsResponse) ProtoMessage() {}

func (x *MemStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_peer_peer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemStatsResponse.ProtoReflect.Descriptor instead.
func (*MemStatsResponse) Descriptor() ([]byte, []int) {
	return file_service_peer_peer_proto_rawDescGZIP(), []int{3}
}

func (x *MemStatsResponse) GetAlloc() uint64 {
	if x != nil {
		return x.Alloc
	}
	return 0
}

func (x *MemStatsResponse) GetTotalAlloc() uint64 {
	if x != nil {
		return x.TotalAlloc
	}
	return 0
}

func (x *MemStatsResponse) GetSys() uint64 {
	if x != nil {
		return x.Sys
	}
	return 0
}

func (x *MemStatsResponse) GetMallocs() uint64 {
	if x != nil {
		return x.Mallocs
	}
	return 0
}

func (x *MemStatsResponse) GetFrees() uint64 {
	if x != nil {
		return x.Frees
	}
	return 0
}

func (x *MemStatsResponse) GetHeapAlloc() uint64 {
	if x != nil {
		return x.HeapAlloc
	}
	return 0
}

func (x *MemStatsResponse) GetHeapSys() uint64 {
	if x != nil {
		return x.HeapSys
	}
	return 0
}

func (x *MemStatsResponse) GetHeapIdle() uint64 {
	if x != nil {
		return x.HeapIdle
	}
	return 0
}

func (x *MemStatsResponse) GetHeapInuse() uint64 {
	if x != nil {
		return x.HeapInuse
	}
	return 0
}

func (x *MemStatsResponse) GetHeapReleased() uint64 {
	if x != nil {
		return x.HeapReleased
	}
	return 0
}

func (x *MemStatsResponse) GetHeapObjects() uint64 {
	if x != nil {
		return x.HeapObjects
	}
	return 0
}

var File_service_peer_peer_proto protoreflect.FileDescriptor

var file_service_peer_peer_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x2f, 0x70,
	0x65, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x65, 0x65, 0x72, 0x70,
	0x62, 0x22, 0x14, 0x0a, 0x12, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xb1, 0x01, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x55, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55,
	0x70, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x65, 0x73, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x4c, 0x6f,
	0x61, 0x64, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x11, 0x0a, 0x0f, 0x4d,
	0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xc2,
	0x02, 0x0a, 0x10, 0x4d, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x79, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x53, 0x79, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4d,
	0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x4d, 0x61,
	0x6c, 0x6c, 0x6f, 0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x72, 0x65, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x46, 0x72, 0x65, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x48,
	0x65, 0x61, 0x70, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x48, 0x65, 0x61, 0x70, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x48, 0x65, 0x61,
	0x70, 0x53, 0x79, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x48, 0x65, 0x61, 0x70,
	0x53, 0x79, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x48, 0x65, 0x61, 0x70, 0x49, 0x64, 0x6c, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x48, 0x65, 0x61, 0x70, 0x49, 0x64, 0x6c, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x70, 0x49, 0x6e, 0x75, 0x73, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x48, 0x65, 0x61, 0x70, 0x49, 0x6e, 0x75, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x0c, 0x48, 0x65, 0x61, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0c, 0x48, 0x65, 0x61, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x70, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x48, 0x65, 0x61, 0x70, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x32, 0x91, 0x01, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x48, 0x0a, 0x0b,
	0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1a, 0x2e, 0x70, 0x65,
	0x65, 0x72, 0x70, 0x62, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x70, 0x62,
	0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x08, 0x4d, 0x65, 0x6d, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x12, 0x17, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x6d, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x65,
	0x65, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x65, 0x65,
	0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_peer_peer_proto_rawDescOnce sync.Once
	file_service_peer_peer_proto_rawDescData = file_service_peer_peer_proto_rawDesc
)

func file_service_peer_peer_proto_rawDescGZIP() []byte {
	file_service_peer_peer_proto_rawDescOnce.Do(func() {
		file_service_peer_peer_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_peer_peer_proto_rawDescData)
	})
	return file_service_peer_peer_proto_rawDescData
}

var file_service_peer_peer_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_service_peer_peer_proto_goTypes = []interface{}{
	(*HealthCheckRequest)(nil),  // 0: peerpb.HealthCheckRequest
	(*HealthCheckResponse)(nil), // 1: peerpb.HealthCheckResponse
	(*MemStatsRequest)(nil),     // 2: peerpb.MemStatsRequest
	(*MemStatsResponse)(nil),    // 3: peerpb.MemStatsResponse
}
var file_service_peer_peer_proto_depIdxs = []int32{
	0, // 0: peerpb.Peer.HealthCheck:input_type -> peerpb.HealthCheckRequest
	2, // 1: peerpb.Peer.MemStats:input_type -> peerpb.MemStatsRequest
	1, // 2: peerpb.Peer.HealthCheck:output_type -> peerpb.HealthCheckResponse
	3, // 3: peerpb.Peer.MemStats:output_type -> peerpb.MemStatsResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_peer_peer_proto_init() }
func file_service_peer_peer_proto_init() {
	if File_service_peer_peer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_peer_peer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckRequest); i {
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
		file_service_peer_peer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckResponse); i {
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
		file_service_peer_peer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemStatsRequest); i {
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
		file_service_peer_peer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemStatsResponse); i {
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
			RawDescriptor: file_service_peer_peer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_peer_peer_proto_goTypes,
		DependencyIndexes: file_service_peer_peer_proto_depIdxs,
		MessageInfos:      file_service_peer_peer_proto_msgTypes,
	}.Build()
	File_service_peer_peer_proto = out.File
	file_service_peer_peer_proto_rawDesc = nil
	file_service_peer_peer_proto_goTypes = nil
	file_service_peer_peer_proto_depIdxs = nil
}