// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: pb/sum.proto

package sum

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// SumRequest is a request for Summator service.
type SumRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A is the number we're adding to. Can't be zero for the sake of example.
	A int64 `protobuf:"varint,1,opt,name=a,proto3" json:"a,omitempty"`
	// B is the number we're adding.
	B *NestedB `protobuf:"bytes,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *SumRequest) Reset() {
	*x = SumRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_sum_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumRequest) ProtoMessage() {}

func (x *SumRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_sum_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumRequest.ProtoReflect.Descriptor instead.
func (*SumRequest) Descriptor() ([]byte, []int) {
	return file_pb_sum_proto_rawDescGZIP(), []int{0}
}

func (x *SumRequest) GetA() int64 {
	if x != nil {
		return x.A
	}
	return 0
}

func (x *SumRequest) GetB() *NestedB {
	if x != nil {
		return x.B
	}
	return nil
}

type SumResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sum   int64  `protobuf:"varint,1,opt,name=sum,proto3" json:"sum,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SumResponse) Reset() {
	*x = SumResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_sum_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumResponse) ProtoMessage() {}

func (x *SumResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_sum_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumResponse.ProtoReflect.Descriptor instead.
func (*SumResponse) Descriptor() ([]byte, []int) {
	return file_pb_sum_proto_rawDescGZIP(), []int{1}
}

func (x *SumResponse) GetSum() int64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

func (x *SumResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type NestedB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	B int64 `protobuf:"varint,1,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *NestedB) Reset() {
	*x = NestedB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_sum_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NestedB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedB) ProtoMessage() {}

func (x *NestedB) ProtoReflect() protoreflect.Message {
	mi := &file_pb_sum_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedB.ProtoReflect.Descriptor instead.
func (*NestedB) Descriptor() ([]byte, []int) {
	return file_pb_sum_proto_rawDescGZIP(), []int{2}
}

func (x *NestedB) GetB() int64 {
	if x != nil {
		return x.B
	}
	return 0
}

var File_pb_sum_proto protoreflect.FileDescriptor

var file_pb_sum_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x62, 0x2f, 0x73, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x73, 0x75, 0x6d, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x38, 0x0a, 0x0a, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x61, 0x12,
	0x1c, 0x0a, 0x01, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x75, 0x6d,
	0x70, 0x62, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x42, 0x52, 0x01, 0x62, 0x22, 0x35, 0x0a,
	0x0b, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x17, 0x0a, 0x07, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x42, 0x12,
	0x0c, 0x0a, 0x01, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x62, 0x32, 0x58, 0x0a,
	0x08, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4c, 0x0a, 0x03, 0x53, 0x75, 0x6d,
	0x12, 0x11, 0x2e, 0x73, 0x75, 0x6d, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x73, 0x75, 0x6d, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22,
	0x13, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x73, 0x75, 0x6d,
	0x2f, 0x7b, 0x61, 0x7d, 0x3a, 0x01, 0x62, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x62, 0x3b,
	0x73, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_sum_proto_rawDescOnce sync.Once
	file_pb_sum_proto_rawDescData = file_pb_sum_proto_rawDesc
)

func file_pb_sum_proto_rawDescGZIP() []byte {
	file_pb_sum_proto_rawDescOnce.Do(func() {
		file_pb_sum_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_sum_proto_rawDescData)
	})
	return file_pb_sum_proto_rawDescData
}

var file_pb_sum_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pb_sum_proto_goTypes = []interface{}{
	(*SumRequest)(nil),  // 0: sumpb.SumRequest
	(*SumResponse)(nil), // 1: sumpb.SumResponse
	(*NestedB)(nil),     // 2: sumpb.NestedB
}
var file_pb_sum_proto_depIdxs = []int32{
	2, // 0: sumpb.SumRequest.b:type_name -> sumpb.NestedB
	0, // 1: sumpb.Summator.Sum:input_type -> sumpb.SumRequest
	1, // 2: sumpb.Summator.Sum:output_type -> sumpb.SumResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pb_sum_proto_init() }
func file_pb_sum_proto_init() {
	if File_pb_sum_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_sum_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumRequest); i {
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
		file_pb_sum_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumResponse); i {
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
		file_pb_sum_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NestedB); i {
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
			RawDescriptor: file_pb_sum_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_sum_proto_goTypes,
		DependencyIndexes: file_pb_sum_proto_depIdxs,
		MessageInfos:      file_pb_sum_proto_msgTypes,
	}.Build()
	File_pb_sum_proto = out.File
	file_pb_sum_proto_rawDesc = nil
	file_pb_sum_proto_goTypes = nil
	file_pb_sum_proto_depIdxs = nil
}
