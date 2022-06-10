// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: gateway/middleware/otel/v1/otel.proto

package v1

import (
	duration "github.com/golang/protobuf/ptypes/duration"
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

// tracing middleware config
type Otel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// report endpoint url
	HttpEndpoint string `protobuf:"bytes,1,opt,name=http_endpoint,json=httpEndpoint,proto3" json:"http_endpoint,omitempty"`
	// sample ratio
	SampleRatio *float32 `protobuf:"fixed32,2,opt,name=sample_ratio,json=sampleRatio,proto3,oneof" json:"sample_ratio,omitempty"`
	// report timeout
	Timeout *duration.Duration `protobuf:"bytes,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *Otel) Reset() {
	*x = Otel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_middleware_otel_v1_otel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Otel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Otel) ProtoMessage() {}

func (x *Otel) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_middleware_otel_v1_otel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Otel.ProtoReflect.Descriptor instead.
func (*Otel) Descriptor() ([]byte, []int) {
	return file_gateway_middleware_otel_v1_otel_proto_rawDescGZIP(), []int{0}
}

func (x *Otel) GetHttpEndpoint() string {
	if x != nil {
		return x.HttpEndpoint
	}
	return ""
}

func (x *Otel) GetSampleRatio() float32 {
	if x != nil && x.SampleRatio != nil {
		return *x.SampleRatio
	}
	return 0
}

func (x *Otel) GetTimeout() *duration.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

var File_gateway_middleware_otel_v1_otel_proto protoreflect.FileDescriptor

var file_gateway_middleware_otel_v1_otel_proto_rawDesc = []byte{
	0x0a, 0x25, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x2f, 0x6f, 0x74, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x74, 0x65,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x6f, 0x74, 0x65, 0x6c,
	0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x99, 0x01, 0x0a, 0x04, 0x4f, 0x74, 0x65, 0x6c, 0x12, 0x23, 0x0a, 0x0d,
	0x68, 0x74, 0x74, 0x70, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x12, 0x26, 0x0a, 0x0c, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x52, 0x61, 0x74, 0x69, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x42, 0x0f,
	0x0a, 0x0d, 0x5f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x42,
	0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f,
	0x2d, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x6d, 0x69, 0x64, 0x64,
	0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x6f, 0x74, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gateway_middleware_otel_v1_otel_proto_rawDescOnce sync.Once
	file_gateway_middleware_otel_v1_otel_proto_rawDescData = file_gateway_middleware_otel_v1_otel_proto_rawDesc
)

func file_gateway_middleware_otel_v1_otel_proto_rawDescGZIP() []byte {
	file_gateway_middleware_otel_v1_otel_proto_rawDescOnce.Do(func() {
		file_gateway_middleware_otel_v1_otel_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_middleware_otel_v1_otel_proto_rawDescData)
	})
	return file_gateway_middleware_otel_v1_otel_proto_rawDescData
}

var file_gateway_middleware_otel_v1_otel_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_gateway_middleware_otel_v1_otel_proto_goTypes = []interface{}{
	(*Otel)(nil),              // 0: gateway.middleware.otel.v1.Otel
	(*duration.Duration)(nil), // 1: google.protobuf.Duration
}
var file_gateway_middleware_otel_v1_otel_proto_depIdxs = []int32{
	1, // 0: gateway.middleware.otel.v1.Otel.timeout:type_name -> google.protobuf.Duration
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gateway_middleware_otel_v1_otel_proto_init() }
func file_gateway_middleware_otel_v1_otel_proto_init() {
	if File_gateway_middleware_otel_v1_otel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_middleware_otel_v1_otel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Otel); i {
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
	file_gateway_middleware_otel_v1_otel_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gateway_middleware_otel_v1_otel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gateway_middleware_otel_v1_otel_proto_goTypes,
		DependencyIndexes: file_gateway_middleware_otel_v1_otel_proto_depIdxs,
		MessageInfos:      file_gateway_middleware_otel_v1_otel_proto_msgTypes,
	}.Build()
	File_gateway_middleware_otel_v1_otel_proto = out.File
	file_gateway_middleware_otel_v1_otel_proto_rawDesc = nil
	file_gateway_middleware_otel_v1_otel_proto_goTypes = nil
	file_gateway_middleware_otel_v1_otel_proto_depIdxs = nil
}
