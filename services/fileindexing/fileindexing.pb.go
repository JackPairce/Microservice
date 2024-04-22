// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: fileindexing.proto

package fileindexing

import (
	types "github.com/JackPairce/MicroService/services/types"
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

type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string        `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values []*types.File `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileindexing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_fileindexing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_fileindexing_proto_rawDescGZIP(), []int{0}
}

func (x *KeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KeyValue) GetValues() []*types.File {
	if x != nil {
		return x.Values
	}
	return nil
}

type JSONData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*KeyValue `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *JSONData) Reset() {
	*x = JSONData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileindexing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JSONData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JSONData) ProtoMessage() {}

func (x *JSONData) ProtoReflect() protoreflect.Message {
	mi := &file_fileindexing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JSONData.ProtoReflect.Descriptor instead.
func (*JSONData) Descriptor() ([]byte, []int) {
	return file_fileindexing_proto_rawDescGZIP(), []int{1}
}

func (x *JSONData) GetData() []*KeyValue {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_fileindexing_proto protoreflect.FileDescriptor

var file_fileindexing_proto_rawDesc = []byte{
	0x0a, 0x12, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x69, 0x6e, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x1a, 0x12, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x3b, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1d,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x30, 0x0a,
	0x08, 0x4a, 0x53, 0x4f, 0x4e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x42,
	0x0f, 0x5a, 0x0d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x69, 0x6e, 0x67,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fileindexing_proto_rawDescOnce sync.Once
	file_fileindexing_proto_rawDescData = file_fileindexing_proto_rawDesc
)

func file_fileindexing_proto_rawDescGZIP() []byte {
	file_fileindexing_proto_rawDescOnce.Do(func() {
		file_fileindexing_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileindexing_proto_rawDescData)
	})
	return file_fileindexing_proto_rawDescData
}

var file_fileindexing_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fileindexing_proto_goTypes = []interface{}{
	(*KeyValue)(nil),   // 0: protos.KeyValue
	(*JSONData)(nil),   // 1: protos.JSONData
	(*types.File)(nil), // 2: File
}
var file_fileindexing_proto_depIdxs = []int32{
	2, // 0: protos.KeyValue.values:type_name -> File
	0, // 1: protos.JSONData.data:type_name -> protos.KeyValue
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_fileindexing_proto_init() }
func file_fileindexing_proto_init() {
	if File_fileindexing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fileindexing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
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
		file_fileindexing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JSONData); i {
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
			RawDescriptor: file_fileindexing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fileindexing_proto_goTypes,
		DependencyIndexes: file_fileindexing_proto_depIdxs,
		MessageInfos:      file_fileindexing_proto_msgTypes,
	}.Build()
	File_fileindexing_proto = out.File
	file_fileindexing_proto_rawDesc = nil
	file_fileindexing_proto_goTypes = nil
	file_fileindexing_proto_depIdxs = nil
}
