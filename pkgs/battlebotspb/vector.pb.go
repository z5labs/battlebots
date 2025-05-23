// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: vector.proto

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

type Vector struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_X0          float64                `protobuf:"fixed64,1,opt,name=x0"`
	xxx_hidden_X1          float64                `protobuf:"fixed64,2,opt,name=x1"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *Vector) Reset() {
	*x = Vector{}
	mi := &file_vector_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Vector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vector) ProtoMessage() {}

func (x *Vector) ProtoReflect() protoreflect.Message {
	mi := &file_vector_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Vector) GetX0() float64 {
	if x != nil {
		return x.xxx_hidden_X0
	}
	return 0
}

func (x *Vector) GetX1() float64 {
	if x != nil {
		return x.xxx_hidden_X1
	}
	return 0
}

func (x *Vector) SetX0(v float64) {
	x.xxx_hidden_X0 = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 2)
}

func (x *Vector) SetX1(v float64) {
	x.xxx_hidden_X1 = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 1, 2)
}

func (x *Vector) HasX0() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *Vector) HasX1() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 1)
}

func (x *Vector) ClearX0() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_X0 = 0
}

func (x *Vector) ClearX1() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 1)
	x.xxx_hidden_X1 = 0
}

type Vector_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	X0 *float64
	X1 *float64
}

func (b0 Vector_builder) Build() *Vector {
	m0 := &Vector{}
	b, x := &b0, m0
	_, _ = b, x
	if b.X0 != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 2)
		x.xxx_hidden_X0 = *b.X0
	}
	if b.X1 != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 1, 2)
		x.xxx_hidden_X1 = *b.X1
	}
	return m0
}

var File_vector_proto protoreflect.FileDescriptor

const file_vector_proto_rawDesc = "" +
	"\n" +
	"\fvector.proto\x12\x13battlebots.protobuf\"(\n" +
	"\x06Vector\x12\x0e\n" +
	"\x02x0\x18\x01 \x01(\x01R\x02x0\x12\x0e\n" +
	"\x02x1\x18\x02 \x01(\x01R\x02x1B=Z;github.com/z5labs/battlebots/pkgs/battlebotspb;battlebotspbb\beditionsp\xe8\a"

var file_vector_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_vector_proto_goTypes = []any{
	(*Vector)(nil), // 0: battlebots.protobuf.Vector
}
var file_vector_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_vector_proto_init() }
func file_vector_proto_init() {
	if File_vector_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_vector_proto_rawDesc), len(file_vector_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vector_proto_goTypes,
		DependencyIndexes: file_vector_proto_depIdxs,
		MessageInfos:      file_vector_proto_msgTypes,
	}.Build()
	File_vector_proto = out.File
	file_vector_proto_goTypes = nil
	file_vector_proto_depIdxs = nil
}
