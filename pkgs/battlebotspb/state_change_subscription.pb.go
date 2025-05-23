// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: state_change_subscription.proto

package battlebotspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StateChangeSubscription struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StateChangeSubscription) Reset() {
	*x = StateChangeSubscription{}
	mi := &file_state_change_subscription_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StateChangeSubscription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateChangeSubscription) ProtoMessage() {}

func (x *StateChangeSubscription) ProtoReflect() protoreflect.Message {
	mi := &file_state_change_subscription_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type StateChangeSubscription_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 StateChangeSubscription_builder) Build() *StateChangeSubscription {
	m0 := &StateChangeSubscription{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

var File_state_change_subscription_proto protoreflect.FileDescriptor

const file_state_change_subscription_proto_rawDesc = "" +
	"\n" +
	"\x1fstate_change_subscription.proto\x12\x13battlebots.protobuf\"\x19\n" +
	"\x17StateChangeSubscriptionB=Z;github.com/z5labs/battlebots/pkgs/battlebotspb;battlebotspbb\beditionsp\xe8\a"

var file_state_change_subscription_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_state_change_subscription_proto_goTypes = []any{
	(*StateChangeSubscription)(nil), // 0: battlebots.protobuf.StateChangeSubscription
}
var file_state_change_subscription_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_state_change_subscription_proto_init() }
func file_state_change_subscription_proto_init() {
	if File_state_change_subscription_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_state_change_subscription_proto_rawDesc), len(file_state_change_subscription_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_state_change_subscription_proto_goTypes,
		DependencyIndexes: file_state_change_subscription_proto_depIdxs,
		MessageInfos:      file_state_change_subscription_proto_msgTypes,
	}.Build()
	File_state_change_subscription_proto = out.File
	file_state_change_subscription_proto_goTypes = nil
	file_state_change_subscription_proto_depIdxs = nil
}
