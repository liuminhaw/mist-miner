// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/miner.proto

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

type NoParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoParam) Reset() {
	*x = NoParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoParam) ProtoMessage() {}

func (x *NoParam) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoParam.ProtoReflect.Descriptor instead.
func (*NoParam) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{0}
}

type MinerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *MinerConfig) Reset() {
	*x = MinerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerConfig) ProtoMessage() {}

func (x *MinerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerConfig.ProtoReflect.Descriptor instead.
func (*MinerConfig) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{1}
}

func (x *MinerConfig) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type MinerPropertyLabel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Unique bool   `protobuf:"varint,2,opt,name=unique,proto3" json:"unique,omitempty"`
}

func (x *MinerPropertyLabel) Reset() {
	*x = MinerPropertyLabel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerPropertyLabel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerPropertyLabel) ProtoMessage() {}

func (x *MinerPropertyLabel) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerPropertyLabel.ProtoReflect.Descriptor instead.
func (*MinerPropertyLabel) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{2}
}

func (x *MinerPropertyLabel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MinerPropertyLabel) GetUnique() bool {
	if x != nil {
		return x.Unique
	}
	return false
}

type MinerPropertyContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Format string `protobuf:"bytes,1,opt,name=format,proto3" json:"format,omitempty"`
	Value  string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MinerPropertyContent) Reset() {
	*x = MinerPropertyContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerPropertyContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerPropertyContent) ProtoMessage() {}

func (x *MinerPropertyContent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerPropertyContent.ProtoReflect.Descriptor instead.
func (*MinerPropertyContent) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{3}
}

func (x *MinerPropertyContent) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *MinerPropertyContent) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type MinerProperty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    string                `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Label   *MinerPropertyLabel   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Content *MinerPropertyContent `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *MinerProperty) Reset() {
	*x = MinerProperty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerProperty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerProperty) ProtoMessage() {}

func (x *MinerProperty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerProperty.ProtoReflect.Descriptor instead.
func (*MinerProperty) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{4}
}

func (x *MinerProperty) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MinerProperty) GetLabel() *MinerPropertyLabel {
	if x != nil {
		return x.Label
	}
	return nil
}

func (x *MinerProperty) GetContent() *MinerPropertyContent {
	if x != nil {
		return x.Content
	}
	return nil
}

type MinerResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Identifier string           `protobuf:"bytes,1,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Properties []*MinerProperty `protobuf:"bytes,2,rep,name=properties,proto3" json:"properties,omitempty"`
}

func (x *MinerResource) Reset() {
	*x = MinerResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerResource) ProtoMessage() {}

func (x *MinerResource) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerResource.ProtoReflect.Descriptor instead.
func (*MinerResource) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{5}
}

func (x *MinerResource) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *MinerResource) GetProperties() []*MinerProperty {
	if x != nil {
		return x.Properties
	}
	return nil
}

type MinerResources struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resources []*MinerResource `protobuf:"bytes,1,rep,name=resources,proto3" json:"resources,omitempty"`
}

func (x *MinerResources) Reset() {
	*x = MinerResources{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinerResources) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinerResources) ProtoMessage() {}

func (x *MinerResources) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinerResources.ProtoReflect.Descriptor instead.
func (*MinerResources) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{6}
}

func (x *MinerResources) GetResources() []*MinerResource {
	if x != nil {
		return x.Resources
	}
	return nil
}

type TestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_miner_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_miner_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_proto_miner_proto_rawDescGZIP(), []int{7}
}

func (x *TestResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_miner_proto protoreflect.FileDescriptor

var file_proto_miner_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x09, 0x0a, 0x07, 0x4e, 0x6f,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0x21, 0x0a, 0x0b, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x40, 0x0a, 0x12, 0x4d, 0x69, 0x6e, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x22, 0x44, 0x0a, 0x14, 0x4d, 0x69,
	0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x8b, 0x01, 0x0a, 0x0d, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72,
	0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69,
	0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x35, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x65,
	0x0a, 0x0d, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12,
	0x34, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x69, 0x65, 0x73, 0x22, 0x44, 0x0a, 0x0e, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22, 0x28, 0x0a, 0x0c, 0x54,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x41, 0x0a, 0x0c, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x4d, 0x69, 0x6e, 0x65, 0x12, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_miner_proto_rawDescOnce sync.Once
	file_proto_miner_proto_rawDescData = file_proto_miner_proto_rawDesc
)

func file_proto_miner_proto_rawDescGZIP() []byte {
	file_proto_miner_proto_rawDescOnce.Do(func() {
		file_proto_miner_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_miner_proto_rawDescData)
	})
	return file_proto_miner_proto_rawDescData
}

var file_proto_miner_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_miner_proto_goTypes = []interface{}{
	(*NoParam)(nil),              // 0: proto.NoParam
	(*MinerConfig)(nil),          // 1: proto.MinerConfig
	(*MinerPropertyLabel)(nil),   // 2: proto.MinerPropertyLabel
	(*MinerPropertyContent)(nil), // 3: proto.MinerPropertyContent
	(*MinerProperty)(nil),        // 4: proto.MinerProperty
	(*MinerResource)(nil),        // 5: proto.MinerResource
	(*MinerResources)(nil),       // 6: proto.MinerResources
	(*TestResponse)(nil),         // 7: proto.TestResponse
}
var file_proto_miner_proto_depIdxs = []int32{
	2, // 0: proto.MinerProperty.label:type_name -> proto.MinerPropertyLabel
	3, // 1: proto.MinerProperty.content:type_name -> proto.MinerPropertyContent
	4, // 2: proto.MinerResource.properties:type_name -> proto.MinerProperty
	5, // 3: proto.MinerResources.resources:type_name -> proto.MinerResource
	1, // 4: proto.MinerService.Mine:input_type -> proto.MinerConfig
	6, // 5: proto.MinerService.Mine:output_type -> proto.MinerResources
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_miner_proto_init() }
func file_proto_miner_proto_init() {
	if File_proto_miner_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_miner_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoParam); i {
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
		file_proto_miner_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerConfig); i {
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
		file_proto_miner_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerPropertyLabel); i {
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
		file_proto_miner_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerPropertyContent); i {
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
		file_proto_miner_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerProperty); i {
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
		file_proto_miner_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerResource); i {
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
		file_proto_miner_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinerResources); i {
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
		file_proto_miner_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResponse); i {
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
			RawDescriptor: file_proto_miner_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_miner_proto_goTypes,
		DependencyIndexes: file_proto_miner_proto_depIdxs,
		MessageInfos:      file_proto_miner_proto_msgTypes,
	}.Build()
	File_proto_miner_proto = out.File
	file_proto_miner_proto_rawDesc = nil
	file_proto_miner_proto_goTypes = nil
	file_proto_miner_proto_depIdxs = nil
}
