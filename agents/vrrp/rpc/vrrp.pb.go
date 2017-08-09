//
// Copyright 2017 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by protoc-gen-go.
// source: vrrp.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	vrrp.proto

It has these top-level messages:
	Reply
	VifEntry
	VifInfo
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
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResultCode int32

const (
	ResultCode_SUCCESS ResultCode = 0
	ResultCode_FAILURE ResultCode = 1
)

var ResultCode_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILURE",
}
var ResultCode_value = map[string]int32{
	"SUCCESS": 0,
	"FAILURE": 1,
}

func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}
func (ResultCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Reply struct {
	Code ResultCode `protobuf:"varint,1,opt,name=code,enum=rpc.ResultCode" json:"code,omitempty"`
}

func (m *Reply) Reset()                    { *m = Reply{} }
func (m *Reply) String() string            { return proto.CompactTextString(m) }
func (*Reply) ProtoMessage()               {}
func (*Reply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Reply) GetCode() ResultCode {
	if m != nil {
		return m.Code
	}
	return ResultCode_SUCCESS
}

type VifEntry struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *VifEntry) Reset()                    { *m = VifEntry{} }
func (m *VifEntry) String() string            { return proto.CompactTextString(m) }
func (*VifEntry) ProtoMessage()               {}
func (*VifEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *VifEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *VifEntry) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type VifInfo struct {
	N       uint64      `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
	Entries []*VifEntry `protobuf:"bytes,2,rep,name=entries" json:"entries,omitempty"`
}

func (m *VifInfo) Reset()                    { *m = VifInfo{} }
func (m *VifInfo) String() string            { return proto.CompactTextString(m) }
func (*VifInfo) ProtoMessage()               {}
func (*VifInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *VifInfo) GetN() uint64 {
	if m != nil {
		return m.N
	}
	return 0
}

func (m *VifInfo) GetEntries() []*VifEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*Reply)(nil), "rpc.Reply")
	proto.RegisterType((*VifEntry)(nil), "rpc.VifEntry")
	proto.RegisterType((*VifInfo)(nil), "rpc.VifInfo")
	proto.RegisterEnum("rpc.ResultCode", ResultCode_name, ResultCode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Vrrp service

type VrrpClient interface {
	GetVifInfo(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*VifInfo, error)
	ToMaster(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*Reply, error)
	ToBackup(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*Reply, error)
}

type vrrpClient struct {
	cc *grpc.ClientConn
}

func NewVrrpClient(cc *grpc.ClientConn) VrrpClient {
	return &vrrpClient{cc}
}

func (c *vrrpClient) GetVifInfo(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*VifInfo, error) {
	out := new(VifInfo)
	err := grpc.Invoke(ctx, "/rpc.Vrrp/GetVifInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vrrpClient) ToMaster(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Vrrp/ToMaster", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vrrpClient) ToBackup(ctx context.Context, in *VifInfo, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Vrrp/ToBackup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Vrrp service

type VrrpServer interface {
	GetVifInfo(context.Context, *VifInfo) (*VifInfo, error)
	ToMaster(context.Context, *VifInfo) (*Reply, error)
	ToBackup(context.Context, *VifInfo) (*Reply, error)
}

func RegisterVrrpServer(s *grpc.Server, srv VrrpServer) {
	s.RegisterService(&_Vrrp_serviceDesc, srv)
}

func _Vrrp_GetVifInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VifInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VrrpServer).GetVifInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Vrrp/GetVifInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VrrpServer).GetVifInfo(ctx, req.(*VifInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vrrp_ToMaster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VifInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VrrpServer).ToMaster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Vrrp/ToMaster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VrrpServer).ToMaster(ctx, req.(*VifInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vrrp_ToBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VifInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VrrpServer).ToBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Vrrp/ToBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VrrpServer).ToBackup(ctx, req.(*VifInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Vrrp_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Vrrp",
	HandlerType: (*VrrpServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVifInfo",
			Handler:    _Vrrp_GetVifInfo_Handler,
		},
		{
			MethodName: "ToMaster",
			Handler:    _Vrrp_ToMaster_Handler,
		},
		{
			MethodName: "ToBackup",
			Handler:    _Vrrp_ToBackup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vrrp.proto",
}

func init() { proto.RegisterFile("vrrp.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0xbb, 0x6d, 0x34, 0x75, 0xea, 0x9f, 0x32, 0xa7, 0xe0, 0xa9, 0x44, 0xa8, 0xa5, 0x48,
	0x0e, 0xf1, 0x0b, 0xa8, 0x21, 0x4a, 0x41, 0x2f, 0x1b, 0x9b, 0x7b, 0x4c, 0x36, 0x10, 0x8c, 0xbb,
	0xcb, 0x64, 0x2b, 0xe4, 0xea, 0x27, 0x97, 0x6c, 0x5d, 0x45, 0x0f, 0xde, 0xe6, 0xcd, 0xfb, 0xed,
	0xbe, 0xe1, 0x01, 0xbc, 0x13, 0xe9, 0x48, 0x93, 0x32, 0x0a, 0x27, 0xa4, 0xcb, 0xf0, 0x0a, 0x0e,
	0xb8, 0xd0, 0x6d, 0x8f, 0x17, 0xe0, 0x95, 0xaa, 0x12, 0x01, 0x5b, 0xb0, 0xd5, 0x69, 0x7c, 0x16,
	0x91, 0x2e, 0x23, 0x2e, 0xba, 0x5d, 0x6b, 0x12, 0x55, 0x09, 0x6e, 0xcd, 0x30, 0x86, 0x69, 0xde,
	0xd4, 0xa9, 0x34, 0xd4, 0x23, 0x82, 0x27, 0x8b, 0xb7, 0xfd, 0x83, 0x23, 0x6e, 0xe7, 0x61, 0x57,
	0x54, 0x15, 0x05, 0xe3, 0xfd, 0x6e, 0x98, 0xc3, 0x1b, 0xf0, 0xf3, 0xa6, 0xde, 0xc8, 0x5a, 0xe1,
	0x31, 0x30, 0x69, 0x79, 0x8f, 0x33, 0x89, 0x97, 0xe0, 0x0b, 0x69, 0xa8, 0x11, 0x5d, 0x30, 0x5e,
	0x4c, 0x56, 0xb3, 0xf8, 0xc4, 0x86, 0xba, 0x00, 0xee, 0xdc, 0xf5, 0x12, 0xe0, 0xe7, 0x12, 0x9c,
	0x81, 0x9f, 0x6d, 0x93, 0x24, 0xcd, 0xb2, 0xf9, 0x68, 0x10, 0xf7, 0xb7, 0x9b, 0xc7, 0x2d, 0x4f,
	0xe7, 0x2c, 0xfe, 0x60, 0xe0, 0xe5, 0x44, 0x1a, 0xd7, 0x00, 0x0f, 0xc2, 0x7c, 0xa7, 0xba, 0x6f,
	0x07, 0x75, 0xfe, 0x4b, 0x85, 0x23, 0x5c, 0xc2, 0xf4, 0x59, 0x3d, 0x15, 0x9d, 0x11, 0xf4, 0x87,
	0x84, 0xaf, 0x0e, 0x74, 0xdb, 0x3b, 0xee, 0xae, 0x28, 0x5f, 0x77, 0xfa, 0x3f, 0xee, 0xe5, 0xd0,
	0x96, 0x7b, 0xfd, 0x19, 0x00, 0x00, 0xff, 0xff, 0x17, 0x59, 0xec, 0xb8, 0x6a, 0x01, 0x00, 0x00,
}
