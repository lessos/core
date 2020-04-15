// Code generated by protoc-gen-go. DO NOT EDIT.
// source: inapi/base.proto

package inapi

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ErrorMeta struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty" toml:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty" toml:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" toml:"-"`
	XXX_unrecognized     []byte   `json:"-" toml:"-"`
	XXX_sizecache        int32    `json:"-" toml:"-"`
}

func (m *ErrorMeta) Reset()         { *m = ErrorMeta{} }
func (m *ErrorMeta) String() string { return proto.CompactTextString(m) }
func (*ErrorMeta) ProtoMessage()    {}
func (*ErrorMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61595a173c150d1, []int{0}
}

func (m *ErrorMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorMeta.Unmarshal(m, b)
}
func (m *ErrorMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorMeta.Marshal(b, m, deterministic)
}
func (m *ErrorMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorMeta.Merge(m, src)
}
func (m *ErrorMeta) XXX_Size() int {
	return xxx_messageInfo_ErrorMeta.Size(m)
}
func (m *ErrorMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorMeta.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorMeta proto.InternalMessageInfo

func (m *ErrorMeta) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ErrorMeta) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type TypeMeta struct {
	Kind                 string     `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty" toml:"kind,omitempty"`
	Error                *ErrorMeta `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty" toml:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-" toml:"-"`
	XXX_unrecognized     []byte     `json:"-" toml:"-"`
	XXX_sizecache        int32      `json:"-" toml:"-"`
}

func (m *TypeMeta) Reset()         { *m = TypeMeta{} }
func (m *TypeMeta) String() string { return proto.CompactTextString(m) }
func (*TypeMeta) ProtoMessage()    {}
func (*TypeMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61595a173c150d1, []int{1}
}

func (m *TypeMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TypeMeta.Unmarshal(m, b)
}
func (m *TypeMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TypeMeta.Marshal(b, m, deterministic)
}
func (m *TypeMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypeMeta.Merge(m, src)
}
func (m *TypeMeta) XXX_Size() int {
	return xxx_messageInfo_TypeMeta.Size(m)
}
func (m *TypeMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TypeMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TypeMeta proto.InternalMessageInfo

func (m *TypeMeta) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *TypeMeta) GetError() *ErrorMeta {
	if m != nil {
		return m.Error
	}
	return nil
}

type ObjectMeta struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" toml:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" toml:"name,omitempty"`
	Created              uint64   `protobuf:"varint,3,opt,name=created,proto3" json:"created,omitempty" toml:"created,omitempty"`
	Updated              uint64   `protobuf:"varint,4,opt,name=updated,proto3" json:"updated,omitempty" toml:"updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" toml:"-"`
	XXX_unrecognized     []byte   `json:"-" toml:"-"`
	XXX_sizecache        int32    `json:"-" toml:"-"`
}

func (m *ObjectMeta) Reset()         { *m = ObjectMeta{} }
func (m *ObjectMeta) String() string { return proto.CompactTextString(m) }
func (*ObjectMeta) ProtoMessage()    {}
func (*ObjectMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61595a173c150d1, []int{2}
}

func (m *ObjectMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObjectMeta.Unmarshal(m, b)
}
func (m *ObjectMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObjectMeta.Marshal(b, m, deterministic)
}
func (m *ObjectMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObjectMeta.Merge(m, src)
}
func (m *ObjectMeta) XXX_Size() int {
	return xxx_messageInfo_ObjectMeta.Size(m)
}
func (m *ObjectMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_ObjectMeta.DiscardUnknown(m)
}

var xxx_messageInfo_ObjectMeta proto.InternalMessageInfo

func (m *ObjectMeta) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ObjectMeta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ObjectMeta) GetCreated() uint64 {
	if m != nil {
		return m.Created
	}
	return 0
}

func (m *ObjectMeta) GetUpdated() uint64 {
	if m != nil {
		return m.Updated
	}
	return 0
}

type Label struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" toml:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty" toml:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" toml:"-"`
	XXX_unrecognized     []byte   `json:"-" toml:"-"`
	XXX_sizecache        int32    `json:"-" toml:"-"`
}

func (m *Label) Reset()         { *m = Label{} }
func (m *Label) String() string { return proto.CompactTextString(m) }
func (*Label) ProtoMessage()    {}
func (*Label) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61595a173c150d1, []int{3}
}

func (m *Label) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Label.Unmarshal(m, b)
}
func (m *Label) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Label.Marshal(b, m, deterministic)
}
func (m *Label) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Label.Merge(m, src)
}
func (m *Label) XXX_Size() int {
	return xxx_messageInfo_Label.Size(m)
}
func (m *Label) XXX_DiscardUnknown() {
	xxx_messageInfo_Label.DiscardUnknown(m)
}

var xxx_messageInfo_Label proto.InternalMessageInfo

func (m *Label) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Label) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*ErrorMeta)(nil), "inapi.ErrorMeta")
	proto.RegisterType((*TypeMeta)(nil), "inapi.TypeMeta")
	proto.RegisterType((*ObjectMeta)(nil), "inapi.ObjectMeta")
	proto.RegisterType((*Label)(nil), "inapi.Label")
}

func init() { proto.RegisterFile("inapi/base.proto", fileDescriptor_b61595a173c150d1) }

var fileDescriptor_b61595a173c150d1 = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0x49, 0x6d, 0xd4, 0x1d, 0x41, 0x96, 0xe0, 0x21, 0xc7, 0x25, 0x07, 0xd9, 0x53, 0x45,
	0x3d, 0xf9, 0x03, 0xf4, 0xa4, 0x08, 0xc5, 0x3f, 0x30, 0x69, 0x06, 0x89, 0xee, 0x36, 0x21, 0x4d,
	0x05, 0xff, 0xbd, 0x64, 0x6a, 0xca, 0xde, 0xde, 0x9b, 0xc7, 0xfb, 0x5e, 0x08, 0x6c, 0xfd, 0x88,
	0xd1, 0xdf, 0x59, 0x9c, 0xa8, 0x8b, 0x29, 0xe4, 0xa0, 0x24, 0x5f, 0xcc, 0x13, 0x6c, 0x9e, 0x53,
	0x0a, 0xe9, 0x8d, 0x32, 0x2a, 0x05, 0xed, 0x10, 0x1c, 0x69, 0xb1, 0x13, 0xfb, 0x4d, 0xcf, 0x5a,
	0x69, 0xb8, 0x38, 0xd2, 0x34, 0xe1, 0x27, 0xe9, 0x86, 0xcf, 0xd5, 0x9a, 0x17, 0xb8, 0xfc, 0xf8,
	0x8d, 0x54, 0x9b, 0xdf, 0x7e, 0x74, 0xb5, 0x59, 0xb4, 0xba, 0x05, 0x49, 0x05, 0xcd, 0xbd, 0xab,
	0x87, 0x6d, 0xc7, 0x8b, 0xdd, 0x3a, 0xd7, 0x2f, 0xb1, 0x71, 0x00, 0xef, 0xf6, 0x8b, 0x86, 0xcc,
	0xa4, 0x6b, 0x68, 0x7c, 0xe5, 0x34, 0xde, 0x15, 0xf2, 0x88, 0xc7, 0x3a, 0xce, 0xba, 0xbc, 0x69,
	0x48, 0x84, 0x99, 0x9c, 0x3e, 0xdb, 0x89, 0x7d, 0xdb, 0x57, 0x5b, 0x92, 0x39, 0x3a, 0x4e, 0xda,
	0x25, 0xf9, 0xb7, 0xe6, 0x1e, 0xe4, 0x2b, 0x5a, 0x3a, 0xac, 0x40, 0x71, 0x02, 0xbc, 0x01, 0xf9,
	0x83, 0x87, 0xb9, 0xae, 0x2c, 0xc6, 0x9e, 0xf3, 0x4f, 0x3d, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff,
	0x67, 0x87, 0xd7, 0x3d, 0x3d, 0x01, 0x00, 0x00,
}
