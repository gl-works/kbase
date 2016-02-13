// Code generated by protoc-gen-go.
// source: rpc/rpc.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rpc/rpc.proto

It has these top-level messages:
	KeyResponse
	Offer
	Revoke
	Poll
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type KeyResponse struct {
	Values map[string][]byte `protobuf:"bytes,1,rep,name=values" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *KeyResponse) Reset()                    { *m = KeyResponse{} }
func (m *KeyResponse) String() string            { return proto.CompactTextString(m) }
func (*KeyResponse) ProtoMessage()               {}
func (*KeyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *KeyResponse) GetValues() map[string][]byte {
	if m != nil {
		return m.Values
	}
	return nil
}

type Offer struct {
	Space  string `protobuf:"bytes,1,opt,name=space" json:"space,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Server string `protobuf:"bytes,3,opt,name=server" json:"server,omitempty"`
	Value  []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Offer) Reset()                    { *m = Offer{} }
func (m *Offer) String() string            { return proto.CompactTextString(m) }
func (*Offer) ProtoMessage()               {}
func (*Offer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Revoke struct {
	Space  string `protobuf:"bytes,1,opt,name=space" json:"space,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Server string `protobuf:"bytes,3,opt,name=server" json:"server,omitempty"`
}

func (m *Revoke) Reset()                    { *m = Revoke{} }
func (m *Revoke) String() string            { return proto.CompactTextString(m) }
func (*Revoke) ProtoMessage()               {}
func (*Revoke) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Poll struct {
	Space string `protobuf:"bytes,1,opt,name=space" json:"space,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
}

func (m *Poll) Reset()                    { *m = Poll{} }
func (m *Poll) String() string            { return proto.CompactTextString(m) }
func (*Poll) ProtoMessage()               {}
func (*Poll) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*KeyResponse)(nil), "rpc.KeyResponse")
	proto.RegisterType((*Offer)(nil), "rpc.Offer")
	proto.RegisterType((*Revoke)(nil), "rpc.Revoke")
	proto.RegisterType((*Poll)(nil), "rpc.Poll")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for KeyService service

type KeyServiceClient interface {
	HandleOffer(ctx context.Context, in *Offer, opts ...grpc.CallOption) (*KeyResponse, error)
	HandleRevoke(ctx context.Context, in *Revoke, opts ...grpc.CallOption) (*KeyResponse, error)
	HandlePoll(ctx context.Context, in *Poll, opts ...grpc.CallOption) (*KeyResponse, error)
}

type keyServiceClient struct {
	cc *grpc.ClientConn
}

func NewKeyServiceClient(cc *grpc.ClientConn) KeyServiceClient {
	return &keyServiceClient{cc}
}

func (c *keyServiceClient) HandleOffer(ctx context.Context, in *Offer, opts ...grpc.CallOption) (*KeyResponse, error) {
	out := new(KeyResponse)
	err := grpc.Invoke(ctx, "/rpc.KeyService/HandleOffer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyServiceClient) HandleRevoke(ctx context.Context, in *Revoke, opts ...grpc.CallOption) (*KeyResponse, error) {
	out := new(KeyResponse)
	err := grpc.Invoke(ctx, "/rpc.KeyService/HandleRevoke", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyServiceClient) HandlePoll(ctx context.Context, in *Poll, opts ...grpc.CallOption) (*KeyResponse, error) {
	out := new(KeyResponse)
	err := grpc.Invoke(ctx, "/rpc.KeyService/HandlePoll", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for KeyService service

type KeyServiceServer interface {
	HandleOffer(context.Context, *Offer) (*KeyResponse, error)
	HandleRevoke(context.Context, *Revoke) (*KeyResponse, error)
	HandlePoll(context.Context, *Poll) (*KeyResponse, error)
}

func RegisterKeyServiceServer(s *grpc.Server, srv KeyServiceServer) {
	s.RegisterService(&_KeyService_serviceDesc, srv)
}

func _KeyService_HandleOffer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Offer)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(KeyServiceServer).HandleOffer(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _KeyService_HandleRevoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Revoke)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(KeyServiceServer).HandleRevoke(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _KeyService_HandlePoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Poll)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(KeyServiceServer).HandlePoll(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _KeyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.KeyService",
	HandlerType: (*KeyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleOffer",
			Handler:    _KeyService_HandleOffer_Handler,
		},
		{
			MethodName: "HandleRevoke",
			Handler:    _KeyService_HandleRevoke_Handler,
		},
		{
			MethodName: "HandlePoll",
			Handler:    _KeyService_HandlePoll_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x2a, 0x48, 0xd6,
	0x07, 0x62, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x66, 0x20, 0x53, 0xa9, 0x8e, 0x8b, 0xdb,
	0x3b, 0xb5, 0x32, 0x28, 0xb5, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x55, 0xc8, 0x84, 0x8b, 0xad, 0x2c,
	0x31, 0xa7, 0x34, 0xb5, 0x58, 0x82, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x46, 0x0f, 0xa4, 0x1e,
	0x49, 0x85, 0x5e, 0x18, 0x58, 0xda, 0x35, 0xaf, 0xa4, 0xa8, 0x32, 0x08, 0xaa, 0x56, 0xca, 0x92,
	0x8b, 0x1b, 0x49, 0x58, 0x48, 0x80, 0x8b, 0x39, 0x3b, 0xb5, 0x12, 0x68, 0x02, 0xa3, 0x06, 0x67,
	0x10, 0x88, 0x29, 0x24, 0xc2, 0xc5, 0x0a, 0x56, 0x2a, 0xc1, 0x04, 0x14, 0xe3, 0x09, 0x82, 0x70,
	0xac, 0x98, 0x2c, 0x18, 0x95, 0x62, 0xb9, 0x58, 0xfd, 0xd3, 0xd2, 0x52, 0x8b, 0x40, 0x4a, 0x8a,
	0x0b, 0x12, 0x93, 0x53, 0xa1, 0xda, 0x20, 0x1c, 0x98, 0x51, 0x4c, 0x08, 0xa3, 0xc4, 0xb8, 0xd8,
	0x8a, 0x53, 0x8b, 0xca, 0x52, 0x8b, 0x24, 0x98, 0xc1, 0x82, 0x50, 0x1e, 0xc2, 0x0a, 0x16, 0x24,
	0x2b, 0x94, 0x3c, 0xb8, 0xd8, 0x82, 0x52, 0xcb, 0xf2, 0xb3, 0x53, 0x29, 0x35, 0x5f, 0x49, 0x8f,
	0x8b, 0x25, 0x20, 0x3f, 0x27, 0x87, 0x58, 0x73, 0x8c, 0x66, 0x32, 0x72, 0x71, 0x01, 0xc3, 0x2d,
	0x18, 0xa8, 0x3b, 0x13, 0xa8, 0x40, 0x97, 0x8b, 0xdb, 0x23, 0x31, 0x2f, 0x25, 0x27, 0x15, 0xe2,
	0x5b, 0x2e, 0x70, 0xb8, 0x82, 0xd9, 0x52, 0x02, 0xe8, 0x61, 0xac, 0xc4, 0x20, 0xa4, 0xcf, 0xc5,
	0x03, 0x51, 0x0e, 0x75, 0x3d, 0x37, 0x58, 0x0d, 0x84, 0x83, 0x55, 0x83, 0x36, 0x17, 0x17, 0x44,
	0x03, 0xd8, 0x91, 0x9c, 0x60, 0x15, 0x20, 0x26, 0x36, 0xc5, 0x49, 0x6c, 0xe0, 0x04, 0x60, 0x0c,
	0x08, 0x00, 0x00, 0xff, 0xff, 0x53, 0xdc, 0x54, 0x7e, 0x11, 0x02, 0x00, 0x00,
}