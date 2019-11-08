// Code created by protoc-gen-gogo. Edited to make code compile.
// source: null.proto

package null

import (
	encoding_binary "encoding/binary"
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	"github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func (d *Decimal) Reset()      { *d = Decimal{} }
func (*Decimal) ProtoMessage() {}
func (*Decimal) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{0}
}

func (d *Decimal) XXX_Unmarshal(b []byte) error {
	return d.Unmarshal(b)
}

func (d *Decimal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Decimal.Marshal(b, d, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := d.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (d *Decimal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Decimal.Merge(d, src)
}

func (d *Decimal) XXX_Size() int {
	return d.Size()
}

func (d *Decimal) XXX_DiscardUnknown() {
	xxx_messageInfo_Decimal.DiscardUnknown(d)
}

var xxx_messageInfo_Decimal proto.InternalMessageInfo

func (m *Bool) Reset()      { *m = Bool{} }
func (*Bool) ProtoMessage() {}
func (*Bool) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{1}
}

func (m *Bool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Bool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Bool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Bool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bool.Merge(m, src)
}

func (m *Bool) XXX_Size() int {
	return m.Size()
}

func (m *Bool) XXX_DiscardUnknown() {
	xxx_messageInfo_Bool.DiscardUnknown(m)
}

var xxx_messageInfo_Bool proto.InternalMessageInfo

func (m *Float64) Reset()      { *m = Float64{} }
func (*Float64) ProtoMessage() {}
func (*Float64) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{2}
}

func (m *Float64) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Float64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Float64.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Float64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Float64.Merge(m, src)
}

func (m *Float64) XXX_Size() int {
	return m.Size()
}

func (m *Float64) XXX_DiscardUnknown() {
	xxx_messageInfo_Float64.DiscardUnknown(m)
}

var xxx_messageInfo_Float64 proto.InternalMessageInfo

func (m *Int64) Reset()      { *m = Int64{} }
func (*Int64) ProtoMessage() {}
func (*Int64) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{3}
}

func (m *Int64) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Int64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Int64.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Int64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Int64.Merge(m, src)
}

func (m *Int64) XXX_Size() int {
	return m.Size()
}

func (m *Int64) XXX_DiscardUnknown() {
	xxx_messageInfo_Int64.DiscardUnknown(m)
}

var xxx_messageInfo_Int64 proto.InternalMessageInfo

func (m *Int32) Reset()      { *m = Int32{} }
func (*Int32) ProtoMessage() {}
func (*Int32) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{4}
}

func (m *Int32) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Int32) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Int32.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Int32) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Int32.Merge(m, src)
}

func (m *Int32) XXX_Size() int {
	return m.Size()
}

func (m *Int32) XXX_DiscardUnknown() {
	xxx_messageInfo_Int32.DiscardUnknown(m)
}

var xxx_messageInfo_Int32 proto.InternalMessageInfo

func (m *Int16) Reset()      { *m = Int16{} }
func (*Int16) ProtoMessage() {}
func (*Int16) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{5}
}

func (m *Int16) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Int16) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Int16.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Int16) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Int16.Merge(m, src)
}

func (m *Int16) XXX_Size() int {
	return m.Size()
}

func (m *Int16) XXX_DiscardUnknown() {
	xxx_messageInfo_Int16.DiscardUnknown(m)
}

var xxx_messageInfo_Int16 proto.InternalMessageInfo

func (m *Int8) Reset()      { *m = Int8{} }
func (*Int8) ProtoMessage() {}
func (*Int8) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{6}
}

func (m *Int8) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Int8) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Int8.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Int8) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Int8.Merge(m, src)
}

func (m *Int8) XXX_Size() int {
	return m.Size()
}

func (m *Int8) XXX_DiscardUnknown() {
	xxx_messageInfo_Int8.DiscardUnknown(m)
}

var xxx_messageInfo_Int8 proto.InternalMessageInfo

func (m *Uint64) Reset()      { *m = Uint64{} }
func (*Uint64) ProtoMessage() {}
func (*Uint64) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{7}
}

func (m *Uint64) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Uint64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Uint64.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Uint64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Uint64.Merge(m, src)
}

func (m *Uint64) XXX_Size() int {
	return m.Size()
}

func (m *Uint64) XXX_DiscardUnknown() {
	xxx_messageInfo_Uint64.DiscardUnknown(m)
}

var xxx_messageInfo_Uint64 proto.InternalMessageInfo

func (m *Uint32) Reset()      { *m = Uint32{} }
func (*Uint32) ProtoMessage() {}
func (*Uint32) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{8}
}

func (m *Uint32) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Uint32) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Uint32.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Uint32) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Uint32.Merge(m, src)
}

func (m *Uint32) XXX_Size() int {
	return m.Size()
}

func (m *Uint32) XXX_DiscardUnknown() {
	xxx_messageInfo_Uint32.DiscardUnknown(m)
}

var xxx_messageInfo_Uint32 proto.InternalMessageInfo

func (m *Uint16) Reset()      { *m = Uint16{} }
func (*Uint16) ProtoMessage() {}
func (*Uint16) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{9}
}

func (m *Uint16) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Uint16) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Uint16.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Uint16) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Uint16.Merge(m, src)
}

func (m *Uint16) XXX_Size() int {
	return m.Size()
}

func (m *Uint16) XXX_DiscardUnknown() {
	xxx_messageInfo_Uint16.DiscardUnknown(m)
}

var xxx_messageInfo_Uint16 proto.InternalMessageInfo

func (a *Uint8) Reset()      { *a = Uint8{} }
func (*Uint8) ProtoMessage() {}
func (*Uint8) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{10}
}

func (a *Uint8) XXX_Unmarshal(b []byte) error {
	return a.Unmarshal(b)
}

func (a *Uint8) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Uint8.Marshal(b, a, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := a.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (a *Uint8) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Uint8.Merge(a, src)
}

func (a *Uint8) XXX_Size() int {
	return a.Size()
}

func (a *Uint8) XXX_DiscardUnknown() {
	xxx_messageInfo_Uint8.DiscardUnknown(a)
}

var xxx_messageInfo_Uint8 proto.InternalMessageInfo

func (a *String) Reset()      { *a = String{} }
func (*String) ProtoMessage() {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{11}
}

func (a *String) XXX_Unmarshal(b []byte) error {
	return a.Unmarshal(b)
}

func (a *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_String.Marshal(b, a, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := a.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (a *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(a, src)
}

func (a *String) XXX_Size() int {
	return a.Size()
}

func (a *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(a)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (a *Time) Reset()      { *a = Time{} }
func (*Time) ProtoMessage() {}
func (*Time) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5db73f817afc81, []int{12}
}

func (a *Time) XXX_Unmarshal(b []byte) error {
	return a.Unmarshal(b)
}

func (a *Time) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Time.Marshal(b, a, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := a.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (a *Time) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Time.Merge(a, src)
}

func (a *Time) XXX_Size() int {
	return a.Size()
}

func (a *Time) XXX_DiscardUnknown() {
	xxx_messageInfo_Time.DiscardUnknown(a)
}

var xxx_messageInfo_Time proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Decimal)(nil), "null.Decimal")
	proto.RegisterType((*Bool)(nil), "null.Bool")
	proto.RegisterType((*Float64)(nil), "null.Float64")
	proto.RegisterType((*Int64)(nil), "null.Int64")
	proto.RegisterType((*Int32)(nil), "null.Int32")
	proto.RegisterType((*Int16)(nil), "null.Int16")
	proto.RegisterType((*Int8)(nil), "null.Int8")
	proto.RegisterType((*Uint64)(nil), "null.Uint64")
	proto.RegisterType((*Uint32)(nil), "null.Uint32")
	proto.RegisterType((*Uint16)(nil), "null.Uint16")
	proto.RegisterType((*Uint8)(nil), "null.Uint8")
	proto.RegisterType((*String)(nil), "null.String")
	proto.RegisterType((*Time)(nil), "null.Time")
}

func init() { proto.RegisterFile("null.proto", fileDescriptor_bf5db73f817afc81) }

var fileDescriptor_bf5db73f817afc81 = []byte{
	// 512 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0x31, 0x6f, 0xda, 0x40,
	0x14, 0x80, 0x39, 0x62, 0x1b, 0x78, 0x2d, 0x8b, 0x55, 0xb5, 0x16, 0xad, 0x0c, 0xa2, 0x0b, 0x4b,
	0x9d, 0x80, 0x2b, 0x44, 0x57, 0x52, 0x55, 0x8a, 0x14, 0x45, 0xe8, 0x08, 0x1d, 0xba, 0x54, 0xc6,
	0xb9, 0xb8, 0x27, 0x9d, 0x7d, 0xd4, 0x3e, 0xe7, 0x77, 0xe4, 0x37, 0x74, 0xcc, 0xd4, 0x9f, 0xc1,
	0xd8, 0xb1, 0x53, 0xdb, 0x90, 0x3f, 0xd1, 0x31, 0xf2, 0xdd, 0x19, 0x18, 0x62, 0xb6, 0xf7, 0xbd,
	0x77, 0xdf, 0xf3, 0xbd, 0x93, 0x1f, 0x40, 0x92, 0x33, 0xe6, 0xad, 0x52, 0x2e, 0xb8, 0x6d, 0x14,
	0x71, 0xe7, 0x5d, 0x44, 0xc5, 0xb7, 0x7c, 0xe9, 0x85, 0x3c, 0x3e, 0x8e, 0x78, 0xc4, 0x8f, 0x65,
	0x71, 0x99, 0x5f, 0x4b, 0x92, 0x20, 0x23, 0x25, 0x75, 0xba, 0x11, 0xe7, 0x11, 0x23, 0xbb, 0x53,
	0x82, 0xc6, 0x24, 0x13, 0x41, 0xbc, 0x52, 0x07, 0xfa, 0x77, 0x08, 0x1a, 0x1f, 0x49, 0x48, 0xe3,
	0x80, 0xd9, 0x6f, 0xa1, 0xbd, 0x4a, 0x49, 0x48, 0x33, 0xca, 0x93, 0xaf, 0x99, 0x48, 0x1d, 0xd4,
	0x43, 0x83, 0x16, 0x7e, 0xbe, 0x4d, 0xce, 0x45, 0x6a, 0xbf, 0x81, 0xd6, 0x96, 0x9d, 0x7a, 0x0f,
	0x0d, 0x0c, 0xbc, 0x4b, 0xd8, 0x2f, 0xc0, 0xcc, 0xc2, 0x80, 0x11, 0xe7, 0xa8, 0x87, 0x06, 0x26,
	0x56, 0x60, 0x77, 0xa0, 0x99, 0x90, 0x28, 0x10, 0xf4, 0x86, 0x38, 0x46, 0x0f, 0x0d, 0x9a, 0x78,
	0xcb, 0x85, 0x71, 0x13, 0x30, 0x7a, 0xe5, 0x98, 0xb2, 0xa0, 0xa0, 0xc8, 0x7e, 0xcf, 0xb9, 0x20,
	0x8e, 0xa5, 0xb2, 0x12, 0xfa, 0x27, 0x60, 0x4c, 0x39, 0x67, 0xb6, 0x0d, 0xc6, 0x92, 0x73, 0x26,
	0xef, 0xd7, 0xc4, 0x32, 0xde, 0xf5, 0xa9, 0xef, 0xf5, 0xe9, 0x7f, 0x80, 0xc6, 0x27, 0xc6, 0x03,
	0x31, 0x7e, 0x6f, 0x3b, 0xd0, 0xb8, 0x56, 0xa1, 0xf4, 0x10, 0x2e, 0xb1, 0x42, 0xf5, 0xc1, 0x3c,
	0x4b, 0x74, 0x99, 0x26, 0xa5, 0x76, 0x84, 0x15, 0x1c, 0x94, 0xfc, 0x91, 0x96, 0xfc, 0x91, 0x94,
	0x4c, 0xac, 0xe0, 0xa0, 0x34, 0x1c, 0x6b, 0x69, 0x38, 0xde, 0x93, 0x54, 0xf6, 0x09, 0xe9, 0x04,
	0x8c, 0xb3, 0x44, 0x4c, 0x8a, 0xb7, 0xa0, 0x89, 0x98, 0x68, 0x45, 0xc6, 0x15, 0xc6, 0x18, 0xac,
	0x85, 0xba, 0xfb, 0x4b, 0xb0, 0xf2, 0xdd, 0x48, 0x06, 0xd6, 0x74, 0xd8, 0xf3, 0x47, 0xa5, 0xa7,
	0xa7, 0x6a, 0x63, 0x4d, 0x87, 0xbd, 0xe1, 0xb8, 0xf4, 0xf4, 0x60, 0xda, 0xab, 0x9c, 0xcc, 0x07,
	0x73, 0x51, 0x8e, 0x91, 0x6f, 0x67, 0x6b, 0x63, 0x05, 0x15, 0xd2, 0x08, 0xac, 0xb9, 0x48, 0x69,
	0x12, 0x15, 0x0f, 0x72, 0x15, 0x88, 0x40, 0xff, 0xbc, 0x32, 0xae, 0x70, 0x3e, 0x83, 0x71, 0x49,
	0x63, 0x62, 0x4f, 0xc0, 0x28, 0xd6, 0x42, 0x1a, 0xcf, 0x46, 0x1d, 0x4f, 0xed, 0x8c, 0x57, 0xee,
	0x8c, 0x77, 0x59, 0xee, 0xcc, 0xb4, 0xb9, 0xfe, 0xd3, 0xad, 0xdd, 0xfe, 0xed, 0x22, 0x2c, 0x8d,
	0xa7, 0xfb, 0x4e, 0xd3, 0xf5, 0xbd, 0x5b, 0xfb, 0x7d, 0xef, 0xd6, 0xd6, 0x1b, 0x17, 0xfd, 0xda,
	0xb8, 0xe8, 0xdf, 0xc6, 0x45, 0xb7, 0x0f, 0x6e, 0xed, 0xe7, 0x83, 0x5b, 0x83, 0xd7, 0x21, 0x8f,
	0xbd, 0x90, 0xa7, 0x24, 0x13, 0x3c, 0xdd, 0xfb, 0x4a, 0xb1, 0xde, 0xd3, 0xd6, 0x45, 0xce, 0xd8,
	0xac, 0x48, 0xcd, 0xd0, 0x17, 0xb9, 0xf1, 0xff, 0x11, 0xfa, 0x51, 0xb7, 0x4e, 0xe7, 0x17, 0x8b,
	0xf3, 0xf3, 0xbb, 0xfa, 0xab, 0x53, 0x9e, 0x92, 0xb9, 0x14, 0x67, 0xa5, 0x58, 0x28, 0x4b, 0x4b,
	0xf6, 0xf1, 0x1f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x18, 0xc5, 0x3b, 0x54, 0x32, 0x04, 0x00, 0x00,
}

func (d *Decimal) Marshal() (dAtA []byte, err error) {
	size := d.Size()
	dAtA = make([]byte, size)
	n, err := d.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (d *Decimal) MarshalTo(dAtA []byte) (int, error) {
	size := d.Size()
	return d.MarshalToSizedBuffer(dAtA[:size])
}

func (d *Decimal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if d.Quote {
		i--
		if d.Quote {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if d.Valid {
		i--
		if d.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if d.Negative {
		i--
		if d.Negative {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if d.Scale != 0 {
		i = encodeVarintNull(dAtA, i, uint64(d.Scale))
		i--
		dAtA[i] = 0x18
	}
	if d.Precision != 0 {
		i = encodeVarintNull(dAtA, i, uint64(d.Precision))
		i--
		dAtA[i] = 0x10
	}
	if len(d.PrecisionStr) > 0 {
		i -= len(d.PrecisionStr)
		copy(dAtA[i:], d.PrecisionStr)
		i = encodeVarintNull(dAtA, i, uint64(len(d.PrecisionStr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Bool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Bool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Bool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Bool {
		i--
		if m.Bool {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Float64) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Float64) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Float64) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Float64 != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Float64))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *Int64) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Int64) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Int64) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Int64 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Int64))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Int32) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Int32) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Int32) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Int32 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Int32))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Int16) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Int16) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Int16) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Int16 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Int16))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Int8) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Int8) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Int8) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Int8 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Int8))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Uint64) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Uint64) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Uint64) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Uint64 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Uint64))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Uint32) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Uint32) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Uint32) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Uint32 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Uint32))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Uint16) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Uint16) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Uint16) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valid {
		i--
		if m.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Uint16 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(m.Uint16))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (a *Uint8) Marshal() (dAtA []byte, err error) {
	size := a.Size()
	dAtA = make([]byte, size)
	n, err := a.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (a *Uint8) MarshalTo(dAtA []byte) (int, error) {
	size := a.Size()
	return a.MarshalToSizedBuffer(dAtA[:size])
}

func (a *Uint8) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if a.Valid {
		i--
		if a.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if a.Uint8 != 0 {
		i = encodeVarintNull(dAtA, i, uint64(a.Uint8))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (a *String) Marshal() (dAtA []byte, err error) {
	size := a.Size()
	dAtA = make([]byte, size)
	n, err := a.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (a *String) MarshalTo(dAtA []byte) (int, error) {
	size := a.Size()
	return a.MarshalToSizedBuffer(dAtA[:size])
}

func (a *String) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if a.Valid {
		i--
		if a.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(a.Data) > 0 {
		i -= len(a.Data)
		copy(dAtA[i:], a.Data)
		i = encodeVarintNull(dAtA, i, uint64(len(a.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (a *Time) Marshal() (dAtA []byte, err error) {
	size := a.Size()
	dAtA = make([]byte, size)
	n, err := a.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (a *Time) MarshalTo(dAtA []byte) (int, error) {
	size := a.Size()
	return a.MarshalToSizedBuffer(dAtA[:size])
}

func (a *Time) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if a.Valid {
		i--
		if a.Valid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(a.Time, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(a.Time):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintNull(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintNull(dAtA []byte, offset int, v uint64) int {
	offset -= sovNull(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func (d Decimal) Size() (n int) {
	var l int
	_ = l
	l = len(d.PrecisionStr)
	if l > 0 {
		n += 1 + l + sovNull(uint64(l))
	}
	if d.Precision != 0 {
		n += 1 + sovNull(uint64(d.Precision))
	}
	if d.Scale != 0 {
		n += 1 + sovNull(uint64(d.Scale))
	}
	if d.Negative {
		n += 2
	}
	if d.Valid {
		n += 2
	}
	if d.Quote {
		n += 2
	}
	return n
}

func (m Bool) Size() (n int) {
	var l int
	_ = l
	if m.Bool {
		n += 2
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Float64) Size() (n int) {
	var l int
	_ = l
	if m.Float64 != 0 {
		n += 9
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Int64) Size() (n int) {
	var l int
	_ = l
	if m.Int64 != 0 {
		n += 1 + sovNull(uint64(m.Int64))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Int32) Size() (n int) {
	var l int
	_ = l
	if m.Int32 != 0 {
		n += 1 + sovNull(uint64(m.Int32))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Int16) Size() (n int) {
	var l int
	_ = l
	if m.Int16 != 0 {
		n += 1 + sovNull(uint64(m.Int16))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Int8) Size() (n int) {
	var l int
	_ = l
	if m.Int8 != 0 {
		n += 1 + sovNull(uint64(m.Int8))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Uint64) Size() (n int) {
	var l int
	_ = l
	if m.Uint64 != 0 {
		n += 1 + sovNull(uint64(m.Uint64))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Uint32) Size() (n int) {
	var l int
	_ = l
	if m.Uint32 != 0 {
		n += 1 + sovNull(uint64(m.Uint32))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (m Uint16) Size() (n int) {
	var l int
	_ = l
	if m.Uint16 != 0 {
		n += 1 + sovNull(uint64(m.Uint16))
	}
	if m.Valid {
		n += 2
	}
	return n
}

func (a Uint8) Size() (n int) {
	var l int
	_ = l
	if a.Uint8 != 0 {
		n += 1 + sovNull(uint64(a.Uint8))
	}
	if a.Valid {
		n += 2
	}
	return n
}

func (a String) Size() (n int) {
	var l int
	_ = l
	l = len(a.Data)
	if l > 0 {
		n += 1 + l + sovNull(uint64(l))
	}
	if a.Valid {
		n += 2
	}
	return n
}

func (a Time) Size() (n int) {
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdTime(a.Time)
	n += 1 + l + sovNull(uint64(l))
	if a.Valid {
		n += 2
	}
	return n
}

func sovNull(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func sozNull(x uint64) (n int) {
	return sovNull(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func (d *Decimal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Decimal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Decimal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrecisionStr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthNull
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNull
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			d.PrecisionStr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Precision", wireType)
			}
			d.Precision = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				d.Precision |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Scale", wireType)
			}
			d.Scale = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				d.Scale |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Negative", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			d.Negative = bool(v != 0)
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			d.Valid = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quote", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			d.Quote = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Bool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Bool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Bool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bool", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Bool = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Float64) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Float64: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Float64: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Float64", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Float64 = float64(math.Float64frombits(v))
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Int64) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Int64: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Int64: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int64", wireType)
			}
			m.Int64 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int64 |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Int32) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Int32: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Int32: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int32", wireType)
			}
			m.Int32 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int32 |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Int16) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Int16: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Int16: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int16", wireType)
			}
			m.Int16 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int16 |= int16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Int8) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Int8: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Int8: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int8", wireType)
			}
			m.Int8 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int8 |= int8(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Uint64) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Uint64: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Uint64: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uint64", wireType)
			}
			m.Uint64 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uint64 |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Uint32) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Uint32: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Uint32: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uint32", wireType)
			}
			m.Uint32 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uint32 |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *Uint16) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Uint16: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Uint16: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uint16", wireType)
			}
			m.Uint16 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uint16 |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (a *Uint8) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Uint8: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Uint8: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uint8", wireType)
			}
			a.Uint8 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				a.Uint8 |= uint8(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			a.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (a *String) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: String: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: String: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthNull
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNull
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			a.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			a.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (a *Time) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNull
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Time: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Time: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthNull
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNull
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&a.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNull
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			a.Valid = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipNull(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNull
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func skipNull(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNull
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowNull
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowNull
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthNull
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthNull
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowNull
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipNull(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthNull
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthNull = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNull   = fmt.Errorf("proto: integer overflow")
)
