// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: inventory/v1/message.proto

package inventoryv1

import (
	v1 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
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

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_inventory_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_inventory_v1_message_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetGamesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version           string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	IncludeProductIds []string `protobuf:"bytes,2,rep,name=include_product_ids,json=includeProductIds,proto3" json:"include_product_ids,omitempty"`
}

func (x *GetGamesRequest) Reset() {
	*x = GetGamesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_v1_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGamesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGamesRequest) ProtoMessage() {}

func (x *GetGamesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_v1_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGamesRequest.ProtoReflect.Descriptor instead.
func (*GetGamesRequest) Descriptor() ([]byte, []int) {
	return file_inventory_v1_message_proto_rawDescGZIP(), []int{2}
}

func (x *GetGamesRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *GetGamesRequest) GetIncludeProductIds() []string {
	if x != nil {
		return x.IncludeProductIds
	}
	return nil
}

type GetGamesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Games []*v1.Game `protobuf:"bytes,1,rep,name=games,proto3" json:"games,omitempty"`
}

func (x *GetGamesResponse) Reset() {
	*x = GetGamesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_v1_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGamesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGamesResponse) ProtoMessage() {}

func (x *GetGamesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_v1_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGamesResponse.ProtoReflect.Descriptor instead.
func (*GetGamesResponse) Descriptor() ([]byte, []int) {
	return file_inventory_v1_message_proto_rawDescGZIP(), []int{3}
}

func (x *GetGamesResponse) GetGames() []*v1.Game {
	if x != nil {
		return x.Games
	}
	return nil
}

var File_inventory_v1_message_proto protoreflect.FileDescriptor

var file_inventory_v1_message_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x69, 0x6e,
	0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x25, 0x74, 0x65, 0x6d, 0x70,
	0x6f, 0x72, 0x61, 0x6c, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x21, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x22, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x47,
	0x61, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x13, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65,
	0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x11, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x49, 0x64, 0x73, 0x22, 0x47, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x47, 0x61, 0x6d, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x05, 0x67, 0x61, 0x6d,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f,
	0x72, 0x61, 0x6c, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x05, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x42, 0xcb,
	0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x74, 0x65, 0x6d,
	0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2d, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x64, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x76,
	0x31, 0x3b, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x49, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x0c, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x18, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inventory_v1_message_proto_rawDescOnce sync.Once
	file_inventory_v1_message_proto_rawDescData = file_inventory_v1_message_proto_rawDesc
)

func file_inventory_v1_message_proto_rawDescGZIP() []byte {
	file_inventory_v1_message_proto_rawDescOnce.Do(func() {
		file_inventory_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_inventory_v1_message_proto_rawDescData)
	})
	return file_inventory_v1_message_proto_rawDescData
}

var file_inventory_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_inventory_v1_message_proto_goTypes = []interface{}{
	(*PingRequest)(nil),      // 0: inventory.v1.PingRequest
	(*PingResponse)(nil),     // 1: inventory.v1.PingResponse
	(*GetGamesRequest)(nil),  // 2: inventory.v1.GetGamesRequest
	(*GetGamesResponse)(nil), // 3: inventory.v1.GetGamesResponse
	(*v1.Game)(nil),          // 4: temporal_shop.values.v1.Game
}
var file_inventory_v1_message_proto_depIdxs = []int32{
	4, // 0: inventory.v1.GetGamesResponse.games:type_name -> temporal_shop.values.v1.Game
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_inventory_v1_message_proto_init() }
func file_inventory_v1_message_proto_init() {
	if File_inventory_v1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inventory_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
		file_inventory_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
		file_inventory_v1_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGamesRequest); i {
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
		file_inventory_v1_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGamesResponse); i {
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
			RawDescriptor: file_inventory_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_inventory_v1_message_proto_goTypes,
		DependencyIndexes: file_inventory_v1_message_proto_depIdxs,
		MessageInfos:      file_inventory_v1_message_proto_msgTypes,
	}.Build()
	File_inventory_v1_message_proto = out.File
	file_inventory_v1_message_proto_rawDesc = nil
	file_inventory_v1_message_proto_goTypes = nil
	file_inventory_v1_message_proto_depIdxs = nil
}
