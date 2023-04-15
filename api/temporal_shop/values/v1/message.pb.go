// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: temporal_shop/values/v1/message.proto

package values

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

type SessionID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *SessionID) Reset() {
	*x = SessionID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_shop_values_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionID) ProtoMessage() {}

func (x *SessionID) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_shop_values_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionID.ProtoReflect.Descriptor instead.
func (*SessionID) Descriptor() ([]byte, []int) {
	return file_temporal_shop_values_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *SessionID) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type Game struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title      string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Category   string `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	ImageUrl   string `protobuf:"bytes,4,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	PriceCents int64  `protobuf:"varint,5,opt,name=price_cents,json=priceCents,proto3" json:"price_cents,omitempty"`
}

func (x *Game) Reset() {
	*x = Game{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_shop_values_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Game) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Game) ProtoMessage() {}

func (x *Game) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_shop_values_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Game.ProtoReflect.Descriptor instead.
func (*Game) Descriptor() ([]byte, []int) {
	return file_temporal_shop_values_v1_message_proto_rawDescGZIP(), []int{1}
}

func (x *Game) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Game) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Game) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Game) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *Game) GetPriceCents() int64 {
	if x != nil {
		return x.PriceCents
	}
	return 0
}

var File_temporal_shop_values_v1_message_proto protoreflect.FileDescriptor

var file_temporal_shop_values_v1_message_proto_rawDesc = []byte{
	0x0a, 0x25, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61,
	0x6c, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x22, 0x21, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x22, 0x86, 0x01, 0x0a, 0x04, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x1b,
	0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0a, 0x70, 0x72, 0x69, 0x63, 0x65, 0x43, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x48, 0x5a, 0x46,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6f,
	0x72, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2d, 0x73,
	0x68, 0x6f, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c,
	0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_temporal_shop_values_v1_message_proto_rawDescOnce sync.Once
	file_temporal_shop_values_v1_message_proto_rawDescData = file_temporal_shop_values_v1_message_proto_rawDesc
)

func file_temporal_shop_values_v1_message_proto_rawDescGZIP() []byte {
	file_temporal_shop_values_v1_message_proto_rawDescOnce.Do(func() {
		file_temporal_shop_values_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_temporal_shop_values_v1_message_proto_rawDescData)
	})
	return file_temporal_shop_values_v1_message_proto_rawDescData
}

var file_temporal_shop_values_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_temporal_shop_values_v1_message_proto_goTypes = []interface{}{
	(*SessionID)(nil), // 0: temporal_shop.values.v1.SessionID
	(*Game)(nil),      // 1: temporal_shop.values.v1.Game
}
var file_temporal_shop_values_v1_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_shop_values_v1_message_proto_init() }
func file_temporal_shop_values_v1_message_proto_init() {
	if File_temporal_shop_values_v1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_temporal_shop_values_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionID); i {
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
		file_temporal_shop_values_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Game); i {
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
			RawDescriptor: file_temporal_shop_values_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_shop_values_v1_message_proto_goTypes,
		DependencyIndexes: file_temporal_shop_values_v1_message_proto_depIdxs,
		MessageInfos:      file_temporal_shop_values_v1_message_proto_msgTypes,
	}.Build()
	File_temporal_shop_values_v1_message_proto = out.File
	file_temporal_shop_values_v1_message_proto_rawDesc = nil
	file_temporal_shop_values_v1_message_proto_goTypes = nil
	file_temporal_shop_values_v1_message_proto_depIdxs = nil
}
