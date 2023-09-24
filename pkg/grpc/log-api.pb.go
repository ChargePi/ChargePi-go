// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.15.8
// source: log-api.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LogResponse) Reset() {
	*x = LogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogResponse) ProtoMessage() {}

func (x *LogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_log_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogResponse.ProtoReflect.Descriptor instead.
func (*LogResponse) Descriptor() ([]byte, []int) {
	return file_log_api_proto_rawDescGZIP(), []int{0}
}

var File_log_api_proto protoreflect.FileDescriptor

var file_log_api_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x2d, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x3e, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x37, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4c, 0x6f,
	0x67, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01,
	0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78,
	0x42, 0x6c, 0x61, 0x7a, 0x33, 0x6b, 0x78, 0x2f, 0x43, 0x68, 0x61, 0x72, 0x67, 0x65, 0x50, 0x69,
	0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_log_api_proto_rawDescOnce sync.Once
	file_log_api_proto_rawDescData = file_log_api_proto_rawDesc
)

func file_log_api_proto_rawDescGZIP() []byte {
	file_log_api_proto_rawDescOnce.Do(func() {
		file_log_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_log_api_proto_rawDescData)
	})
	return file_log_api_proto_rawDescData
}

var file_log_api_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_log_api_proto_goTypes = []interface{}{
	(*LogResponse)(nil),   // 0: api.LogResponse
	(*emptypb.Empty)(nil), // 1: google.protobuf.Empty
}
var file_log_api_proto_depIdxs = []int32{
	1, // 0: api.Log.GetLogs:input_type -> google.protobuf.Empty
	0, // 1: api.Log.GetLogs:output_type -> api.LogResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_log_api_proto_init() }
func file_log_api_proto_init() {
	if File_log_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_log_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogResponse); i {
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
			RawDescriptor: file_log_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_log_api_proto_goTypes,
		DependencyIndexes: file_log_api_proto_depIdxs,
		MessageInfos:      file_log_api_proto_msgTypes,
	}.Build()
	File_log_api_proto = out.File
	file_log_api_proto_rawDesc = nil
	file_log_api_proto_goTypes = nil
	file_log_api_proto_depIdxs = nil
}