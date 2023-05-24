// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: proto/client.proto

package proto

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

type ClientBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   int32                          `protobuf:"varint,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver int32                          `protobuf:"varint,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	UniqueId string                         `protobuf:"bytes,3,opt,name=unique_id,json=uniqueId,proto3" json:"unique_id,omitempty"`
	Type     int32                          `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"` // 1 for request and 2 for response
	Note     string                         `protobuf:"bytes,5,opt,name=note,proto3" json:"note,omitempty"`
	Requests []*ClientBatch_SingleOperation `protobuf:"bytes,6,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (x *ClientBatch) Reset() {
	*x = ClientBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientBatch) ProtoMessage() {}

func (x *ClientBatch) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientBatch.ProtoReflect.Descriptor instead.
func (*ClientBatch) Descriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{0}
}

func (x *ClientBatch) GetSender() int32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

func (x *ClientBatch) GetReceiver() int32 {
	if x != nil {
		return x.Receiver
	}
	return 0
}

func (x *ClientBatch) GetUniqueId() string {
	if x != nil {
		return x.UniqueId
	}
	return ""
}

func (x *ClientBatch) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ClientBatch) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *ClientBatch) GetRequests() []*ClientBatch_SingleOperation {
	if x != nil {
		return x.Requests
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   int32  `protobuf:"varint,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver int32  `protobuf:"varint,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	UniqueId string `protobuf:"bytes,3,opt,name=unique_id,json=uniqueId,proto3" json:"unique_id,omitempty"`
	Type     int32  `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"` // 1 for bootstrap, 2 for log print, 3 consensus start
	Note     string `protobuf:"bytes,5,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_client_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetSender() int32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

func (x *Status) GetReceiver() int32 {
	if x != nil {
		return x.Receiver
	}
	return 0
}

func (x *Status) GetUniqueId() string {
	if x != nil {
		return x.UniqueId
	}
	return ""
}

func (x *Status) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Status) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

type ClientBatch_SingleOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Command string `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *ClientBatch_SingleOperation) Reset() {
	*x = ClientBatch_SingleOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_client_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientBatch_SingleOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientBatch_SingleOperation) ProtoMessage() {}

func (x *ClientBatch_SingleOperation) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientBatch_SingleOperation.ProtoReflect.Descriptor instead.
func (*ClientBatch_SingleOperation) Descriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ClientBatch_SingleOperation) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClientBatch_SingleOperation) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

var File_proto_client_proto protoreflect.FileDescriptor

var file_proto_client_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfd, 0x01, 0x0a, 0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x6e, 0x69, 0x71,
	0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x6e, 0x69,
	0x71, 0x75, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x12, 0x38, 0x0a,
	0x08, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x53, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x1a, 0x3b, 0x0a, 0x0f, 0x53, 0x69, 0x6e, 0x67, 0x6c,
	0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x22, 0x81, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_client_proto_rawDescOnce sync.Once
	file_proto_client_proto_rawDescData = file_proto_client_proto_rawDesc
)

func file_proto_client_proto_rawDescGZIP() []byte {
	file_proto_client_proto_rawDescOnce.Do(func() {
		file_proto_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_client_proto_rawDescData)
	})
	return file_proto_client_proto_rawDescData
}

var file_proto_client_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_client_proto_goTypes = []interface{}{
	(*ClientBatch)(nil),                 // 0: ClientBatch
	(*Status)(nil),                      // 1: Status
	(*ClientBatch_SingleOperation)(nil), // 2: ClientBatch.SingleOperation
}
var file_proto_client_proto_depIdxs = []int32{
	2, // 0: ClientBatch.requests:type_name -> ClientBatch.SingleOperation
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_client_proto_init() }
func file_proto_client_proto_init() {
	if File_proto_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientBatch); i {
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
		file_proto_client_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_proto_client_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientBatch_SingleOperation); i {
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
			RawDescriptor: file_proto_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_client_proto_goTypes,
		DependencyIndexes: file_proto_client_proto_depIdxs,
		MessageInfos:      file_proto_client_proto_msgTypes,
	}.Build()
	File_proto_client_proto = out.File
	file_proto_client_proto_rawDesc = nil
	file_proto_client_proto_goTypes = nil
	file_proto_client_proto_depIdxs = nil
}
