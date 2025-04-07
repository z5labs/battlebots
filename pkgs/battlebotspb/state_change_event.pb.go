// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: state_change_event.proto

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

type StateChangeEvent struct {
	state              protoimpl.MessageState     `protogen:"opaque.v1"`
	xxx_hidden_Subject isStateChangeEvent_Subject `protobuf_oneof:"subject"`
	xxx_hidden_Event   isStateChangeEvent_Event   `protobuf_oneof:"event"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *StateChangeEvent) Reset() {
	*x = StateChangeEvent{}
	mi := &file_state_change_event_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StateChangeEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateChangeEvent) ProtoMessage() {}

func (x *StateChangeEvent) ProtoReflect() protoreflect.Message {
	mi := &file_state_change_event_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *StateChangeEvent) GetBot() *Bot {
	if x != nil {
		if x, ok := x.xxx_hidden_Subject.(*stateChangeEvent_Bot); ok {
			return x.Bot
		}
	}
	return nil
}

func (x *StateChangeEvent) GetPosition() *Vector {
	if x != nil {
		if x, ok := x.xxx_hidden_Event.(*stateChangeEvent_Position); ok {
			return x.Position
		}
	}
	return nil
}

func (x *StateChangeEvent) SetBot(v *Bot) {
	if v == nil {
		x.xxx_hidden_Subject = nil
		return
	}
	x.xxx_hidden_Subject = &stateChangeEvent_Bot{v}
}

func (x *StateChangeEvent) SetPosition(v *Vector) {
	if v == nil {
		x.xxx_hidden_Event = nil
		return
	}
	x.xxx_hidden_Event = &stateChangeEvent_Position{v}
}

func (x *StateChangeEvent) HasSubject() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Subject != nil
}

func (x *StateChangeEvent) HasBot() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_Subject.(*stateChangeEvent_Bot)
	return ok
}

func (x *StateChangeEvent) HasEvent() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Event != nil
}

func (x *StateChangeEvent) HasPosition() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_Event.(*stateChangeEvent_Position)
	return ok
}

func (x *StateChangeEvent) ClearSubject() {
	x.xxx_hidden_Subject = nil
}

func (x *StateChangeEvent) ClearBot() {
	if _, ok := x.xxx_hidden_Subject.(*stateChangeEvent_Bot); ok {
		x.xxx_hidden_Subject = nil
	}
}

func (x *StateChangeEvent) ClearEvent() {
	x.xxx_hidden_Event = nil
}

func (x *StateChangeEvent) ClearPosition() {
	if _, ok := x.xxx_hidden_Event.(*stateChangeEvent_Position); ok {
		x.xxx_hidden_Event = nil
	}
}

const StateChangeEvent_Subject_not_set_case case_StateChangeEvent_Subject = 0
const StateChangeEvent_Bot_case case_StateChangeEvent_Subject = 1

func (x *StateChangeEvent) WhichSubject() case_StateChangeEvent_Subject {
	if x == nil {
		return StateChangeEvent_Subject_not_set_case
	}
	switch x.xxx_hidden_Subject.(type) {
	case *stateChangeEvent_Bot:
		return StateChangeEvent_Bot_case
	default:
		return StateChangeEvent_Subject_not_set_case
	}
}

const StateChangeEvent_Event_not_set_case case_StateChangeEvent_Event = 0
const StateChangeEvent_Position_case case_StateChangeEvent_Event = 2

func (x *StateChangeEvent) WhichEvent() case_StateChangeEvent_Event {
	if x == nil {
		return StateChangeEvent_Event_not_set_case
	}
	switch x.xxx_hidden_Event.(type) {
	case *stateChangeEvent_Position:
		return StateChangeEvent_Position_case
	default:
		return StateChangeEvent_Event_not_set_case
	}
}

type StateChangeEvent_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// Fields of oneof xxx_hidden_Subject:
	Bot *Bot
	// -- end of xxx_hidden_Subject
	// Fields of oneof xxx_hidden_Event:
	Position *Vector
	// -- end of xxx_hidden_Event
}

func (b0 StateChangeEvent_builder) Build() *StateChangeEvent {
	m0 := &StateChangeEvent{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Bot != nil {
		x.xxx_hidden_Subject = &stateChangeEvent_Bot{b.Bot}
	}
	if b.Position != nil {
		x.xxx_hidden_Event = &stateChangeEvent_Position{b.Position}
	}
	return m0
}

type case_StateChangeEvent_Subject protoreflect.FieldNumber

func (x case_StateChangeEvent_Subject) String() string {
	md := file_state_change_event_proto_msgTypes[0].Descriptor()
	if x == 0 {
		return "not set"
	}
	return protoimpl.X.MessageFieldStringOf(md, protoreflect.FieldNumber(x))
}

type case_StateChangeEvent_Event protoreflect.FieldNumber

func (x case_StateChangeEvent_Event) String() string {
	md := file_state_change_event_proto_msgTypes[0].Descriptor()
	if x == 0 {
		return "not set"
	}
	return protoimpl.X.MessageFieldStringOf(md, protoreflect.FieldNumber(x))
}

type isStateChangeEvent_Subject interface {
	isStateChangeEvent_Subject()
}

type stateChangeEvent_Bot struct {
	Bot *Bot `protobuf:"bytes,1,opt,name=bot,oneof"`
}

func (*stateChangeEvent_Bot) isStateChangeEvent_Subject() {}

type isStateChangeEvent_Event interface {
	isStateChangeEvent_Event()
}

type stateChangeEvent_Position struct {
	Position *Vector `protobuf:"bytes,2,opt,name=position,oneof"`
}

func (*stateChangeEvent_Position) isStateChangeEvent_Event() {}

var File_state_change_event_proto protoreflect.FileDescriptor

const file_state_change_event_proto_rawDesc = "" +
	"\n" +
	"\x18state_change_event.proto\x12\x13battlebots.protobuf\x1a\tbot.proto\x1a\fvector.proto\"\x8f\x01\n" +
	"\x10StateChangeEvent\x12,\n" +
	"\x03bot\x18\x01 \x01(\v2\x18.battlebots.protobuf.BotH\x00R\x03bot\x129\n" +
	"\bposition\x18\x02 \x01(\v2\x1b.battlebots.protobuf.VectorH\x01R\bpositionB\t\n" +
	"\asubjectB\a\n" +
	"\x05eventB=Z;github.com/z5labs/battlebots/pkgs/battlebotspb;battlebotspbb\beditionsp\xe8\a"

var file_state_change_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_state_change_event_proto_goTypes = []any{
	(*StateChangeEvent)(nil), // 0: battlebots.protobuf.StateChangeEvent
	(*Bot)(nil),              // 1: battlebots.protobuf.Bot
	(*Vector)(nil),           // 2: battlebots.protobuf.Vector
}
var file_state_change_event_proto_depIdxs = []int32{
	1, // 0: battlebots.protobuf.StateChangeEvent.bot:type_name -> battlebots.protobuf.Bot
	2, // 1: battlebots.protobuf.StateChangeEvent.position:type_name -> battlebots.protobuf.Vector
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_state_change_event_proto_init() }
func file_state_change_event_proto_init() {
	if File_state_change_event_proto != nil {
		return
	}
	file_bot_proto_init()
	file_vector_proto_init()
	file_state_change_event_proto_msgTypes[0].OneofWrappers = []any{
		(*stateChangeEvent_Bot)(nil),
		(*stateChangeEvent_Position)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_state_change_event_proto_rawDesc), len(file_state_change_event_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_state_change_event_proto_goTypes,
		DependencyIndexes: file_state_change_event_proto_depIdxs,
		MessageInfos:      file_state_change_event_proto_msgTypes,
	}.Build()
	File_state_change_event_proto = out.File
	file_state_change_event_proto_goTypes = nil
	file_state_change_event_proto_depIdxs = nil
}
