// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: admin.proto

package admins

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

type ValidateAdminRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *ValidateAdminRequest) Reset() {
	*x = ValidateAdminRequest{}
	mi := &file_admin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateAdminRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateAdminRequest) ProtoMessage() {}

func (x *ValidateAdminRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateAdminRequest.ProtoReflect.Descriptor instead.
func (*ValidateAdminRequest) Descriptor() ([]byte, []int) {
	return file_admin_proto_rawDescGZIP(), []int{0}
}

func (x *ValidateAdminRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ValidateAdminRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type ValidateAdminResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsValid bool   `protobuf:"varint,1,opt,name=is_valid,json=isValid,proto3" json:"is_valid,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	AdminId int64  `protobuf:"varint,3,opt,name=admin_id,json=adminId,proto3" json:"admin_id,omitempty"`
}

func (x *ValidateAdminResponse) Reset() {
	*x = ValidateAdminResponse{}
	mi := &file_admin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateAdminResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateAdminResponse) ProtoMessage() {}

func (x *ValidateAdminResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateAdminResponse.ProtoReflect.Descriptor instead.
func (*ValidateAdminResponse) Descriptor() ([]byte, []int) {
	return file_admin_proto_rawDescGZIP(), []int{1}
}

func (x *ValidateAdminResponse) GetIsValid() bool {
	if x != nil {
		return x.IsValid
	}
	return false
}

func (x *ValidateAdminResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ValidateAdminResponse) GetAdminId() int64 {
	if x != nil {
		return x.AdminId
	}
	return 0
}

var File_admin_proto protoreflect.FileDescriptor

var file_admin_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x73, 0x22, 0x48, 0x0a, 0x14, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x67, 0x0a, 0x15, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x64, 0x32, 0x5c, 0x0a, 0x0c, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4c, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x73,
	0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_admin_proto_rawDescOnce sync.Once
	file_admin_proto_rawDescData = file_admin_proto_rawDesc
)

func file_admin_proto_rawDescGZIP() []byte {
	file_admin_proto_rawDescOnce.Do(func() {
		file_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_proto_rawDescData)
	})
	return file_admin_proto_rawDescData
}

var file_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_admin_proto_goTypes = []any{
	(*ValidateAdminRequest)(nil),  // 0: admins.ValidateAdminRequest
	(*ValidateAdminResponse)(nil), // 1: admins.ValidateAdminResponse
}
var file_admin_proto_depIdxs = []int32{
	0, // 0: admins.AdminService.ValidateAdmin:input_type -> admins.ValidateAdminRequest
	1, // 1: admins.AdminService.ValidateAdmin:output_type -> admins.ValidateAdminResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_admin_proto_init() }
func file_admin_proto_init() {
	if File_admin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_proto_goTypes,
		DependencyIndexes: file_admin_proto_depIdxs,
		MessageInfos:      file_admin_proto_msgTypes,
	}.Build()
	File_admin_proto = out.File
	file_admin_proto_rawDesc = nil
	file_admin_proto_goTypes = nil
	file_admin_proto_depIdxs = nil
}