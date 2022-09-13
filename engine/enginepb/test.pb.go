// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: engine/proto/test.proto

package enginepb

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

type Record_RecordType int32

const (
	Record_Data Record_RecordType = 0
	Record_DDL  Record_RecordType = 1
)

// Enum value maps for Record_RecordType.
var (
	Record_RecordType_name = map[int32]string{
		0: "Data",
		1: "DDL",
	}
	Record_RecordType_value = map[string]int32{
		"Data": 0,
		"DDL":  1,
	}
)

func (x Record_RecordType) Enum() *Record_RecordType {
	p := new(Record_RecordType)
	*p = x
	return p
}

func (x Record_RecordType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Record_RecordType) Descriptor() protoreflect.EnumDescriptor {
	return file_engine_proto_test_proto_enumTypes[0].Descriptor()
}

func (Record_RecordType) Type() protoreflect.EnumType {
	return &file_engine_proto_test_proto_enumTypes[0]
}

func (x Record_RecordType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Record_RecordType.Descriptor instead.
func (Record_RecordType) EnumDescriptor() ([]byte, []int) {
	return file_engine_proto_test_proto_rawDescGZIP(), []int{0, 0}
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tp        Record_RecordType `protobuf:"varint,1,opt,name=tp,proto3,enum=enginepb.Record_RecordType" json:"tp,omitempty"`
	SchemaVer int32             `protobuf:"varint,2,opt,name=schema_ver,json=schemaVer,proto3" json:"schema_ver,omitempty"`
	Tid       int32             `protobuf:"varint,3,opt,name=tid,proto3" json:"tid,omitempty"`
	Gtid      int32             `protobuf:"varint,4,opt,name=gtid,proto3" json:"gtid,omitempty"`
	Pk        int32             `protobuf:"varint,5,opt,name=pk,proto3" json:"pk,omitempty"`
	// for record time
	TimeTracer []int64 `protobuf:"varint,6,rep,packed,name=time_tracer,json=timeTracer,proto3" json:"time_tracer,omitempty"`
	// error
	Err *Error `protobuf:"bytes,7,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_engine_proto_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_engine_proto_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_engine_proto_test_proto_rawDescGZIP(), []int{0}
}

func (x *Record) GetTp() Record_RecordType {
	if x != nil {
		return x.Tp
	}
	return Record_Data
}

func (x *Record) GetSchemaVer() int32 {
	if x != nil {
		return x.SchemaVer
	}
	return 0
}

func (x *Record) GetTid() int32 {
	if x != nil {
		return x.Tid
	}
	return 0
}

func (x *Record) GetGtid() int32 {
	if x != nil {
		return x.Gtid
	}
	return 0
}

func (x *Record) GetPk() int32 {
	if x != nil {
		return x.Pk
	}
	return 0
}

func (x *Record) GetTimeTracer() []int64 {
	if x != nil {
		return x.TimeTracer
	}
	return nil
}

func (x *Record) GetErr() *Error {
	if x != nil {
		return x.Err
	}
	return nil
}

type TestBinlogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gtid int32 `protobuf:"varint,1,opt,name=gtid,proto3" json:"gtid,omitempty"`
}

func (x *TestBinlogRequest) Reset() {
	*x = TestBinlogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_engine_proto_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestBinlogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestBinlogRequest) ProtoMessage() {}

func (x *TestBinlogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_engine_proto_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestBinlogRequest.ProtoReflect.Descriptor instead.
func (*TestBinlogRequest) Descriptor() ([]byte, []int) {
	return file_engine_proto_test_proto_rawDescGZIP(), []int{1}
}

func (x *TestBinlogRequest) GetGtid() int32 {
	if x != nil {
		return x.Gtid
	}
	return 0
}

var File_engine_proto_test_proto protoreflect.FileDescriptor

var file_engine_proto_test_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x70, 0x62, 0x1a, 0x18, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x01,
	0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x2b, 0x0a, 0x02, 0x74, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x02, 0x74, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f,
	0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x56, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x74, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x74, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x67, 0x74, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6b,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x70, 0x6b, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x72, 0x18, 0x06, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x0a, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x03, 0x65,
	0x72, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x70, 0x62, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x1f,
	0x0a, 0x0a, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x44, 0x4c, 0x10, 0x01, 0x22,
	0x27, 0x0a, 0x11, 0x54, 0x65, 0x73, 0x74, 0x42, 0x69, 0x6e, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x67, 0x74, 0x69, 0x64, 0x32, 0x4c, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x46, 0x65, 0x65, 0x64, 0x42,
	0x69, 0x6e, 0x6c, 0x6f, 0x67, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x70, 0x62,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x42, 0x69, 0x6e, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x30, 0x01, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x63, 0x61, 0x70, 0x2f, 0x74, 0x69, 0x66,
	0x6c, 0x6f, 0x77, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_engine_proto_test_proto_rawDescOnce sync.Once
	file_engine_proto_test_proto_rawDescData = file_engine_proto_test_proto_rawDesc
)

func file_engine_proto_test_proto_rawDescGZIP() []byte {
	file_engine_proto_test_proto_rawDescOnce.Do(func() {
		file_engine_proto_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_engine_proto_test_proto_rawDescData)
	})
	return file_engine_proto_test_proto_rawDescData
}

var file_engine_proto_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_engine_proto_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_engine_proto_test_proto_goTypes = []interface{}{
	(Record_RecordType)(0),    // 0: enginepb.Record.RecordType
	(*Record)(nil),            // 1: enginepb.Record
	(*TestBinlogRequest)(nil), // 2: enginepb.TestBinlogRequest
	(*Error)(nil),             // 3: enginepb.Error
}
var file_engine_proto_test_proto_depIdxs = []int32{
	0, // 0: enginepb.Record.tp:type_name -> enginepb.Record.RecordType
	3, // 1: enginepb.Record.err:type_name -> enginepb.Error
	2, // 2: enginepb.TestService.FeedBinlog:input_type -> enginepb.TestBinlogRequest
	1, // 3: enginepb.TestService.FeedBinlog:output_type -> enginepb.Record
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_engine_proto_test_proto_init() }
func file_engine_proto_test_proto_init() {
	if File_engine_proto_test_proto != nil {
		return
	}
	file_engine_proto_error_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_engine_proto_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_engine_proto_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestBinlogRequest); i {
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
			RawDescriptor: file_engine_proto_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_engine_proto_test_proto_goTypes,
		DependencyIndexes: file_engine_proto_test_proto_depIdxs,
		EnumInfos:         file_engine_proto_test_proto_enumTypes,
		MessageInfos:      file_engine_proto_test_proto_msgTypes,
	}.Build()
	File_engine_proto_test_proto = out.File
	file_engine_proto_test_proto_rawDesc = nil
	file_engine_proto_test_proto_goTypes = nil
	file_engine_proto_test_proto_depIdxs = nil
}
