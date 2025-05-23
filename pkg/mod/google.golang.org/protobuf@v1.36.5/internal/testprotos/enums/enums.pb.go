// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/testprotos/enums/enums.proto

package enums

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

type Enum int32

const (
	Enum_DEFAULT     Enum = 1337
	Enum_ZERO        Enum = 0
	Enum_ONE         Enum = 1
	Enum_ELEVENT     Enum = 11
	Enum_SEVENTEEN   Enum = 17
	Enum_THIRTYSEVEN Enum = 37
	Enum_SIXTYSEVEN  Enum = 67
	Enum_NEGATIVE    Enum = -1
)

// Enum value maps for Enum.
var (
	Enum_name = map[int32]string{
		1337: "DEFAULT",
		0:    "ZERO",
		1:    "ONE",
		11:   "ELEVENT",
		17:   "SEVENTEEN",
		37:   "THIRTYSEVEN",
		67:   "SIXTYSEVEN",
		-1:   "NEGATIVE",
	}
	Enum_value = map[string]int32{
		"DEFAULT":     1337,
		"ZERO":        0,
		"ONE":         1,
		"ELEVENT":     11,
		"SEVENTEEN":   17,
		"THIRTYSEVEN": 37,
		"SIXTYSEVEN":  67,
		"NEGATIVE":    -1,
	}
)

func (x Enum) Enum() *Enum {
	p := new(Enum)
	*p = x
	return p
}

func (x Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Enum) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_testprotos_enums_enums_proto_enumTypes[0].Descriptor()
}

func (Enum) Type() protoreflect.EnumType {
	return &file_internal_testprotos_enums_enums_proto_enumTypes[0]
}

func (x Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Enum.Descriptor instead.
func (Enum) EnumDescriptor() ([]byte, []int) {
	return file_internal_testprotos_enums_enums_proto_rawDescGZIP(), []int{0}
}

var File_internal_testprotos_enums_enums_proto protoreflect.FileDescriptor

var file_internal_testprotos_enums_enums_proto_rawDesc = string([]byte{
	0x0a, 0x25, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2a, 0x7b, 0x0a, 0x04,
	0x45, 0x6e, 0x75, 0x6d, 0x12, 0x0c, 0x0a, 0x07, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10,
	0xb9, 0x0a, 0x12, 0x08, 0x0a, 0x04, 0x5a, 0x45, 0x52, 0x4f, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03,
	0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x4c, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x10, 0x0b, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x45, 0x45, 0x4e, 0x10,
	0x11, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x48, 0x49, 0x52, 0x54, 0x59, 0x53, 0x45, 0x56, 0x45, 0x4e,
	0x10, 0x25, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x49, 0x58, 0x54, 0x59, 0x53, 0x45, 0x56, 0x45, 0x4e,
	0x10, 0x43, 0x12, 0x15, 0x0a, 0x08, 0x4e, 0x45, 0x47, 0x41, 0x54, 0x49, 0x56, 0x45, 0x10, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x42, 0x3b, 0x5a, 0x34, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x73, 0x92, 0x03, 0x02, 0x10, 0x02, 0x62, 0x08, 0x65, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x70, 0xe8, 0x07,
})

var (
	file_internal_testprotos_enums_enums_proto_rawDescOnce sync.Once
	file_internal_testprotos_enums_enums_proto_rawDescData []byte
)

func file_internal_testprotos_enums_enums_proto_rawDescGZIP() []byte {
	file_internal_testprotos_enums_enums_proto_rawDescOnce.Do(func() {
		file_internal_testprotos_enums_enums_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_testprotos_enums_enums_proto_rawDesc), len(file_internal_testprotos_enums_enums_proto_rawDesc)))
	})
	return file_internal_testprotos_enums_enums_proto_rawDescData
}

var file_internal_testprotos_enums_enums_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_testprotos_enums_enums_proto_goTypes = []any{
	(Enum)(0), // 0: goproto.proto.enums.Enum
}
var file_internal_testprotos_enums_enums_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_testprotos_enums_enums_proto_init() }
func file_internal_testprotos_enums_enums_proto_init() {
	if File_internal_testprotos_enums_enums_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_testprotos_enums_enums_proto_rawDesc), len(file_internal_testprotos_enums_enums_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_testprotos_enums_enums_proto_goTypes,
		DependencyIndexes: file_internal_testprotos_enums_enums_proto_depIdxs,
		EnumInfos:         file_internal_testprotos_enums_enums_proto_enumTypes,
	}.Build()
	File_internal_testprotos_enums_enums_proto = out.File
	file_internal_testprotos_enums_enums_proto_goTypes = nil
	file_internal_testprotos_enums_enums_proto_depIdxs = nil
}
