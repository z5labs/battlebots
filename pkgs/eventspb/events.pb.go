// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.30.0--dev
// source: events.proto

package eventspb

import (
	vectorpb "github.com/z5labs/battlebots/pkgs/vectorpb"
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

type Events2DRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Events2DRequest) Reset() {
	*x = Events2DRequest{}
	mi := &file_events_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Events2DRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Events2DRequest) ProtoMessage() {}

func (x *Events2DRequest) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type Events2DRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 Events2DRequest_builder) Build() *Events2DRequest {
	m0 := &Events2DRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type Event2D struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Event isEvent2D_Event        `protobuf_oneof:"event"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Event2D) Reset() {
	*x = Event2D{}
	mi := &file_events_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Event2D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event2D) ProtoMessage() {}

func (x *Event2D) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Event2D) GetPositionUpdate() *PositionUpdate2D {
	if x != nil {
		if x, ok := x.xxx_hidden_Event.(*event2D_PositionUpdate); ok {
			return x.PositionUpdate
		}
	}
	return nil
}

func (x *Event2D) SetPositionUpdate(v *PositionUpdate2D) {
	if v == nil {
		x.xxx_hidden_Event = nil
		return
	}
	x.xxx_hidden_Event = &event2D_PositionUpdate{v}
}

func (x *Event2D) HasEvent() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Event != nil
}

func (x *Event2D) HasPositionUpdate() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_Event.(*event2D_PositionUpdate)
	return ok
}

func (x *Event2D) ClearEvent() {
	x.xxx_hidden_Event = nil
}

func (x *Event2D) ClearPositionUpdate() {
	if _, ok := x.xxx_hidden_Event.(*event2D_PositionUpdate); ok {
		x.xxx_hidden_Event = nil
	}
}

const Event2D_Event_not_set_case case_Event2D_Event = 0
const Event2D_PositionUpdate_case case_Event2D_Event = 1

func (x *Event2D) WhichEvent() case_Event2D_Event {
	if x == nil {
		return Event2D_Event_not_set_case
	}
	switch x.xxx_hidden_Event.(type) {
	case *event2D_PositionUpdate:
		return Event2D_PositionUpdate_case
	default:
		return Event2D_Event_not_set_case
	}
}

type Event2D_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// Fields of oneof xxx_hidden_Event:
	PositionUpdate *PositionUpdate2D
	// -- end of xxx_hidden_Event
}

func (b0 Event2D_builder) Build() *Event2D {
	m0 := &Event2D{}
	b, x := &b0, m0
	_, _ = b, x
	if b.PositionUpdate != nil {
		x.xxx_hidden_Event = &event2D_PositionUpdate{b.PositionUpdate}
	}
	return m0
}

type case_Event2D_Event protoreflect.FieldNumber

func (x case_Event2D_Event) String() string {
	md := file_events_proto_msgTypes[1].Descriptor()
	if x == 0 {
		return "not set"
	}
	return protoimpl.X.MessageFieldStringOf(md, protoreflect.FieldNumber(x))
}

type isEvent2D_Event interface {
	isEvent2D_Event()
}

type event2D_PositionUpdate struct {
	PositionUpdate *PositionUpdate2D `protobuf:"bytes,1,opt,name=position_update,json=positionUpdate,oneof"`
}

func (*event2D_PositionUpdate) isEvent2D_Event() {}

type PositionUpdate2D struct {
	state               protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Position *vectorpb.Vector2      `protobuf:"bytes,1,opt,name=position"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *PositionUpdate2D) Reset() {
	*x = PositionUpdate2D{}
	mi := &file_events_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PositionUpdate2D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PositionUpdate2D) ProtoMessage() {}

func (x *PositionUpdate2D) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PositionUpdate2D) GetPosition() *vectorpb.Vector2 {
	if x != nil {
		return x.xxx_hidden_Position
	}
	return nil
}

func (x *PositionUpdate2D) SetPosition(v *vectorpb.Vector2) {
	x.xxx_hidden_Position = v
}

func (x *PositionUpdate2D) HasPosition() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Position != nil
}

func (x *PositionUpdate2D) ClearPosition() {
	x.xxx_hidden_Position = nil
}

type PositionUpdate2D_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Position *vectorpb.Vector2
}

func (b0 PositionUpdate2D_builder) Build() *PositionUpdate2D {
	m0 := &PositionUpdate2D{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Position = b.Position
	return m0
}

type Events3DRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Events3DRequest) Reset() {
	*x = Events3DRequest{}
	mi := &file_events_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Events3DRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Events3DRequest) ProtoMessage() {}

func (x *Events3DRequest) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type Events3DRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 Events3DRequest_builder) Build() *Events3DRequest {
	m0 := &Events3DRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type Event3D struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Event isEvent3D_Event        `protobuf_oneof:"event"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Event3D) Reset() {
	*x = Event3D{}
	mi := &file_events_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Event3D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event3D) ProtoMessage() {}

func (x *Event3D) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Event3D) GetPositionUpdate() *PositionUpdate3D {
	if x != nil {
		if x, ok := x.xxx_hidden_Event.(*event3D_PositionUpdate); ok {
			return x.PositionUpdate
		}
	}
	return nil
}

func (x *Event3D) SetPositionUpdate(v *PositionUpdate3D) {
	if v == nil {
		x.xxx_hidden_Event = nil
		return
	}
	x.xxx_hidden_Event = &event3D_PositionUpdate{v}
}

func (x *Event3D) HasEvent() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Event != nil
}

func (x *Event3D) HasPositionUpdate() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_Event.(*event3D_PositionUpdate)
	return ok
}

func (x *Event3D) ClearEvent() {
	x.xxx_hidden_Event = nil
}

func (x *Event3D) ClearPositionUpdate() {
	if _, ok := x.xxx_hidden_Event.(*event3D_PositionUpdate); ok {
		x.xxx_hidden_Event = nil
	}
}

const Event3D_Event_not_set_case case_Event3D_Event = 0
const Event3D_PositionUpdate_case case_Event3D_Event = 1

func (x *Event3D) WhichEvent() case_Event3D_Event {
	if x == nil {
		return Event3D_Event_not_set_case
	}
	switch x.xxx_hidden_Event.(type) {
	case *event3D_PositionUpdate:
		return Event3D_PositionUpdate_case
	default:
		return Event3D_Event_not_set_case
	}
}

type Event3D_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// Fields of oneof xxx_hidden_Event:
	PositionUpdate *PositionUpdate3D
	// -- end of xxx_hidden_Event
}

func (b0 Event3D_builder) Build() *Event3D {
	m0 := &Event3D{}
	b, x := &b0, m0
	_, _ = b, x
	if b.PositionUpdate != nil {
		x.xxx_hidden_Event = &event3D_PositionUpdate{b.PositionUpdate}
	}
	return m0
}

type case_Event3D_Event protoreflect.FieldNumber

func (x case_Event3D_Event) String() string {
	md := file_events_proto_msgTypes[4].Descriptor()
	if x == 0 {
		return "not set"
	}
	return protoimpl.X.MessageFieldStringOf(md, protoreflect.FieldNumber(x))
}

type isEvent3D_Event interface {
	isEvent3D_Event()
}

type event3D_PositionUpdate struct {
	PositionUpdate *PositionUpdate3D `protobuf:"bytes,1,opt,name=position_update,json=positionUpdate,oneof"`
}

func (*event3D_PositionUpdate) isEvent3D_Event() {}

type PositionUpdate3D struct {
	state               protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Position *vectorpb.Vector3      `protobuf:"bytes,1,opt,name=position"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *PositionUpdate3D) Reset() {
	*x = PositionUpdate3D{}
	mi := &file_events_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PositionUpdate3D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PositionUpdate3D) ProtoMessage() {}

func (x *PositionUpdate3D) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PositionUpdate3D) GetPosition() *vectorpb.Vector3 {
	if x != nil {
		return x.xxx_hidden_Position
	}
	return nil
}

func (x *PositionUpdate3D) SetPosition(v *vectorpb.Vector3) {
	x.xxx_hidden_Position = v
}

func (x *PositionUpdate3D) HasPosition() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Position != nil
}

func (x *PositionUpdate3D) ClearPosition() {
	x.xxx_hidden_Position = nil
}

type PositionUpdate3D_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Position *vectorpb.Vector3
}

func (b0 PositionUpdate3D_builder) Build() *PositionUpdate3D {
	m0 := &PositionUpdate3D{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Position = b.Position
	return m0
}

var File_events_proto protoreflect.FileDescriptor

var file_events_proto_rawDesc = string([]byte{
	0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x62, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x1a, 0x0c, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x44, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x64, 0x0a, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x32, 0x44, 0x12,
	0x50, 0x0a, 0x0f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c,
	0x65, 0x62, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x32, 0x44, 0x48,
	0x00, 0x52, 0x0e, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x4c, 0x0a, 0x10, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x32, 0x44, 0x12, 0x38,
	0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x62, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x32, 0x52, 0x08,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x11, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x33, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x64, 0x0a, 0x07, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x33, 0x44, 0x12, 0x50, 0x0a, 0x0f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x62, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x33, 0x44, 0x48, 0x00, 0x52, 0x0e, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x22, 0x4c, 0x0a, 0x10, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x33, 0x44, 0x12, 0x38, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x62, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x33, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x35,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x62, 0x6f, 0x74, 0x73, 0x2f,
	0x70, 0x6b, 0x67, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x70, 0x62, 0x62, 0x08, 0x65,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
})

var file_events_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_events_proto_goTypes = []any{
	(*Events2DRequest)(nil),  // 0: battlebots.protobuf.Events2DRequest
	(*Event2D)(nil),          // 1: battlebots.protobuf.Event2D
	(*PositionUpdate2D)(nil), // 2: battlebots.protobuf.PositionUpdate2D
	(*Events3DRequest)(nil),  // 3: battlebots.protobuf.Events3DRequest
	(*Event3D)(nil),          // 4: battlebots.protobuf.Event3D
	(*PositionUpdate3D)(nil), // 5: battlebots.protobuf.PositionUpdate3D
	(*vectorpb.Vector2)(nil), // 6: battlebots.protobuf.Vector2
	(*vectorpb.Vector3)(nil), // 7: battlebots.protobuf.Vector3
}
var file_events_proto_depIdxs = []int32{
	2, // 0: battlebots.protobuf.Event2D.position_update:type_name -> battlebots.protobuf.PositionUpdate2D
	6, // 1: battlebots.protobuf.PositionUpdate2D.position:type_name -> battlebots.protobuf.Vector2
	5, // 2: battlebots.protobuf.Event3D.position_update:type_name -> battlebots.protobuf.PositionUpdate3D
	7, // 3: battlebots.protobuf.PositionUpdate3D.position:type_name -> battlebots.protobuf.Vector3
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_events_proto_init() }
func file_events_proto_init() {
	if File_events_proto != nil {
		return
	}
	file_events_proto_msgTypes[1].OneofWrappers = []any{
		(*event2D_PositionUpdate)(nil),
	}
	file_events_proto_msgTypes[4].OneofWrappers = []any{
		(*event3D_PositionUpdate)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_events_proto_rawDesc), len(file_events_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_proto_goTypes,
		DependencyIndexes: file_events_proto_depIdxs,
		MessageInfos:      file_events_proto_msgTypes,
	}.Build()
	File_events_proto = out.File
	file_events_proto_goTypes = nil
	file_events_proto_depIdxs = nil
}
