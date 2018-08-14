// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto.proto

// +build csall proto

package modification

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func (m *Observers) Reset()         { *m = Observers{} }
func (m *Observers) String() string { return proto.CompactTextString(m) }
func (*Observers) ProtoMessage()    {}
func (*Observers) Descriptor() ([]byte, []int) {
	return fileDescriptor_proto_72f976f22bbab417, []int{0}
}
func (m *Observers) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Observers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Observers.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Observers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Observers.Merge(dst, src)
}
func (m *Observers) XXX_Size() int {
	return m.Size()
}
func (m *Observers) XXX_DiscardUnknown() {
	xxx_messageInfo_Observers.DiscardUnknown(m)
}

var xxx_messageInfo_Observers proto.InternalMessageInfo

func (*Observers) XXX_MessageName() string {
	return "modification.Observers"
}
func (m *Observer) Reset()         { *m = Observer{} }
func (m *Observer) String() string { return proto.CompactTextString(m) }
func (*Observer) ProtoMessage()    {}
func (*Observer) Descriptor() ([]byte, []int) {
	return fileDescriptor_proto_72f976f22bbab417, []int{1}
}
func (m *Observer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Observer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Observer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Observer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Observer.Merge(dst, src)
}
func (m *Observer) XXX_Size() int {
	return m.Size()
}
func (m *Observer) XXX_DiscardUnknown() {
	xxx_messageInfo_Observer.DiscardUnknown(m)
}

var xxx_messageInfo_Observer proto.InternalMessageInfo

func (*Observer) XXX_MessageName() string {
	return "modification.Observer"
}
func init() {
	proto.RegisterType((*Observers)(nil), "modification.Observers")
	proto.RegisterType((*Observer)(nil), "modification.Observer")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProtoObserverServiceClient is the client API for ProtoObserverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProtoObserverServiceClient interface {
	Register(ctx context.Context, in *Observers, opts ...grpc.CallOption) (*types.Empty, error)
	Deregister(ctx context.Context, in *Observers, opts ...grpc.CallOption) (*types.Empty, error)
}

type protoObserverServiceClient struct {
	cc *grpc.ClientConn
}

func NewProtoObserverServiceClient(cc *grpc.ClientConn) ProtoObserverServiceClient {
	return &protoObserverServiceClient{cc}
}

func (c *protoObserverServiceClient) Register(ctx context.Context, in *Observers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/modification.ProtoObserverService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protoObserverServiceClient) Deregister(ctx context.Context, in *Observers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/modification.ProtoObserverService/Deregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProtoObserverService service

type ProtoObserverServiceServer interface {
	Register(context.Context, *Observers) (*types.Empty, error)
	Deregister(context.Context, *Observers) (*types.Empty, error)
}

func RegisterProtoObserverServiceServer(s *grpc.Server, srv ProtoObserverServiceServer) {
	s.RegisterService(&_ProtoObserverService_serviceDesc, srv)
}

func _ProtoObserverService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Observers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoObserverServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/modification.ProtoObserverService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoObserverServiceServer).Register(ctx, req.(*Observers))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoObserverService_Deregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Observers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoObserverServiceServer).Deregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/modification.ProtoObserverService/Deregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoObserverServiceServer).Deregister(ctx, req.(*Observers))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProtoObserverService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "modification.ProtoObserverService",
	HandlerType: (*ProtoObserverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ProtoObserverService_Register_Handler,
		},
		{
			MethodName: "Deregister",
			Handler:    _ProtoObserverService_Deregister_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto.proto",
}

func (m *Observers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Observers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Collection) > 0 {
		for _, msg := range m.Collection {
			dAtA[i] = 0xa
			i++
			i = encodeVarintProto(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Observer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Observer) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Route) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintProto(dAtA, i, uint64(len(m.Route)))
		i += copy(dAtA[i:], m.Route)
	}
	if len(m.Event) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintProto(dAtA, i, uint64(len(m.Event)))
		i += copy(dAtA[i:], m.Event)
	}
	if len(m.Type) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintProto(dAtA, i, uint64(len(m.Type)))
		i += copy(dAtA[i:], m.Type)
	}
	if len(m.Condition) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintProto(dAtA, i, uint64(len(m.Condition)))
		i += copy(dAtA[i:], m.Condition)
	}
	return i, nil
}

func encodeVarintProto(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Observers) Size() (n int) {
	var l int
	_ = l
	if len(m.Collection) > 0 {
		for _, e := range m.Collection {
			l = e.Size()
			n += 1 + l + sovProto(uint64(l))
		}
	}
	return n
}

func (m *Observer) Size() (n int) {
	var l int
	_ = l
	l = len(m.Route)
	if l > 0 {
		n += 1 + l + sovProto(uint64(l))
	}
	l = len(m.Event)
	if l > 0 {
		n += 1 + l + sovProto(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovProto(uint64(l))
	}
	l = len(m.Condition)
	if l > 0 {
		n += 1 + l + sovProto(uint64(l))
	}
	return n
}

func sovProto(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozProto(x uint64) (n int) {
	return sovProto(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Observers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProto
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Observers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Observers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Collection", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProto
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProto
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Collection = append(m.Collection, &Observer{})
			if err := m.Collection[len(m.Collection)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProto(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProto
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
func (m *Observer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProto
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Observer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Observer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Route", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProto
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProto
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Route = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Event", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProto
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProto
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Event = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProto
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProto
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Condition", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProto
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProto
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Condition = append(m.Condition[:0], dAtA[iNdEx:postIndex]...)
			if m.Condition == nil {
				m.Condition = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProto(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProto
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
func skipProto(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProto
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
					return 0, ErrIntOverflowProto
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
					return 0, ErrIntOverflowProto
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthProto
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowProto
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
				next, err := skipProto(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
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
	ErrInvalidLengthProto = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProto   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("proto.proto", fileDescriptor_proto_72f976f22bbab417) }

var fileDescriptor_proto_72f976f22bbab417 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xc1, 0x4e, 0x32, 0x31,
	0x14, 0x85, 0xa7, 0x3f, 0xf0, 0x87, 0x29, 0xb8, 0x69, 0x0c, 0x4e, 0xd0, 0x74, 0x08, 0x2b, 0xa2,
	0xb1, 0x24, 0x98, 0xb8, 0x33, 0x26, 0x20, 0x6b, 0x4d, 0x75, 0xe5, 0xce, 0x19, 0x2e, 0x63, 0x13,
	0xa0, 0x93, 0x4e, 0x21, 0xe1, 0x2d, 0x5c, 0xe9, 0x6b, 0xf8, 0x18, 0x2c, 0x7d, 0x02, 0xa2, 0xe5,
	0x45, 0x4c, 0x3b, 0x99, 0x0c, 0x26, 0xae, 0xdc, 0x4c, 0xe6, 0x9e, 0xf3, 0xb5, 0xf7, 0xde, 0x53,
	0xdc, 0x48, 0x95, 0xd4, 0x92, 0xb9, 0x2f, 0x69, 0xce, 0xe5, 0x44, 0x4c, 0x45, 0xfc, 0xa4, 0x85,
	0x5c, 0xb4, 0x8f, 0x13, 0x29, 0x93, 0x19, 0xf4, 0x9d, 0x17, 0x2d, 0xa7, 0x7d, 0x98, 0xa7, 0x7a,
	0x9d, 0xa3, 0xed, 0xf3, 0x44, 0xe8, 0xe7, 0x65, 0xc4, 0x62, 0x39, 0xef, 0x27, 0x32, 0x91, 0x25,
	0x65, 0x2b, 0x57, 0xb8, 0xbf, 0x1c, 0xef, 0x8e, 0xb0, 0x7f, 0x1b, 0x65, 0xa0, 0x56, 0xa0, 0x32,
	0x72, 0x89, 0x71, 0x2c, 0x67, 0x33, 0x88, 0x6d, 0x9b, 0x00, 0x75, 0x2a, 0xbd, 0xc6, 0xa0, 0xc5,
	0xf6, 0x7b, 0xb3, 0x02, 0xe6, 0x7b, 0x64, 0xf7, 0x0d, 0xe1, 0x7a, 0x61, 0x90, 0x10, 0xd7, 0x94,
	0x5c, 0x6a, 0x08, 0x50, 0x07, 0xf5, 0xfc, 0xa1, 0x6f, 0xb6, 0x61, 0x8d, 0x5b, 0x81, 0xe7, 0xba,
	0x05, 0x60, 0x05, 0x0b, 0x1d, 0xfc, 0x2b, 0x81, 0xb1, 0x15, 0x78, 0xae, 0x93, 0x13, 0x5c, 0xd5,
	0xeb, 0x14, 0x82, 0x8a, 0xf3, 0xeb, 0x66, 0x1b, 0x56, 0x1f, 0xd6, 0x29, 0x70, 0xa7, 0x92, 0x33,
	0xec, 0xc7, 0x72, 0x31, 0x11, 0x6e, 0xc6, 0x6a, 0x07, 0xf5, 0x9a, 0xc3, 0x03, 0xb3, 0x0d, 0xfd,
	0x51, 0x21, 0xf2, 0xd2, 0x1f, 0xbc, 0x22, 0x7c, 0x78, 0x67, 0x17, 0x2d, 0xc6, 0xbb, 0x07, 0xb5,
	0x12, 0x31, 0x90, 0x2b, 0x5c, 0xe7, 0x90, 0x88, 0x4c, 0x83, 0x22, 0x47, 0xbf, 0xaf, 0x98, 0xb5,
	0x5b, 0x2c, 0x4f, 0x9a, 0x15, 0x19, 0xb2, 0xb1, 0x4d, 0xba, 0xeb, 0x91, 0x6b, 0x8c, 0x6f, 0x40,
	0xfd, 0xfd, 0x82, 0xe1, 0xe9, 0xe6, 0x8b, 0x7a, 0x1b, 0x43, 0xd1, 0x87, 0xa1, 0xe8, 0xd3, 0x50,
	0xf4, 0xb2, 0xa3, 0xde, 0xfb, 0x8e, 0x7a, 0x9b, 0x1d, 0x45, 0x8f, 0x3f, 0xde, 0x3b, 0xfa, 0xef,
	0x4e, 0x5f, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0xa0, 0xa2, 0x22, 0x06, 0x13, 0x02, 0x00, 0x00,
}