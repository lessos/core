// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0-devel
// 	protoc        v3.5.0
// source: operator.proto

package inapi

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PbOpLogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" toml:"name,omitempty"` // struct:object_slice_key
	Status  string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty" toml:"status,omitempty"`
	Updated uint64 `protobuf:"varint,3,opt,name=updated,proto3" json:"updated,omitempty" toml:"updated,omitempty"` // unix time in milliseconds
	Message string `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty" toml:"message,omitempty"`
}

func (x *PbOpLogEntry) Reset() {
	*x = PbOpLogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbOpLogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbOpLogEntry) ProtoMessage() {}

func (x *PbOpLogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_operator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbOpLogEntry.ProtoReflect.Descriptor instead.
func (*PbOpLogEntry) Descriptor() ([]byte, []int) {
	return file_operator_proto_rawDescGZIP(), []int{0}
}

func (x *PbOpLogEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PbOpLogEntry) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *PbOpLogEntry) GetUpdated() uint64 {
	if x != nil {
		return x.Updated
	}
	return 0
}

func (x *PbOpLogEntry) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PbOpLogSets struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" toml:"name,omitempty"` // struct:object_slice_key
	Version uint32          `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty" toml:"version,omitempty"`
	Items   []*PbOpLogEntry `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty" toml:"items,omitempty"`
}

func (x *PbOpLogSets) Reset() {
	*x = PbOpLogSets{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbOpLogSets) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbOpLogSets) ProtoMessage() {}

func (x *PbOpLogSets) ProtoReflect() protoreflect.Message {
	mi := &file_operator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbOpLogSets.ProtoReflect.Descriptor instead.
func (*PbOpLogSets) Descriptor() ([]byte, []int) {
	return file_operator_proto_rawDescGZIP(), []int{1}
}

func (x *PbOpLogSets) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PbOpLogSets) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *PbOpLogSets) GetItems() []*PbOpLogEntry {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_operator_proto protoreflect.FileDescriptor

var file_operator_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x22, 0x6e, 0x0a, 0x0c, 0x50, 0x62, 0x4f, 0x70, 0x4c,
	0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x66, 0x0a, 0x0b, 0x50, 0x62, 0x4f, 0x70, 0x4c,
	0x6f, 0x67, 0x53, 0x65, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x62, 0x4f, 0x70,
	0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x42,
	0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_operator_proto_rawDescOnce sync.Once
	file_operator_proto_rawDescData = file_operator_proto_rawDesc
)

func file_operator_proto_rawDescGZIP() []byte {
	file_operator_proto_rawDescOnce.Do(func() {
		file_operator_proto_rawDescData = protoimpl.X.CompressGZIP(file_operator_proto_rawDescData)
	})
	return file_operator_proto_rawDescData
}

var file_operator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_operator_proto_goTypes = []interface{}{
	(*PbOpLogEntry)(nil), // 0: inapi.PbOpLogEntry
	(*PbOpLogSets)(nil),  // 1: inapi.PbOpLogSets
}
var file_operator_proto_depIdxs = []int32{
	0, // 0: inapi.PbOpLogSets.items:type_name -> inapi.PbOpLogEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_operator_proto_init() }
func file_operator_proto_init() {
	if File_operator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_operator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbOpLogEntry); i {
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
		file_operator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbOpLogSets); i {
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
			RawDescriptor: file_operator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_operator_proto_goTypes,
		DependencyIndexes: file_operator_proto_depIdxs,
		MessageInfos:      file_operator_proto_msgTypes,
	}.Build()
	File_operator_proto = out.File
	file_operator_proto_rawDesc = nil
	file_operator_proto_goTypes = nil
	file_operator_proto_depIdxs = nil
}
