// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/miner.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MinerServiceClient is the client API for MinerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MinerServiceClient interface {
	Mine(ctx context.Context, in *MinerConfig, opts ...grpc.CallOption) (*MinerResources, error)
}

type minerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMinerServiceClient(cc grpc.ClientConnInterface) MinerServiceClient {
	return &minerServiceClient{cc}
}

func (c *minerServiceClient) Mine(ctx context.Context, in *MinerConfig, opts ...grpc.CallOption) (*MinerResources, error) {
	out := new(MinerResources)
	err := c.cc.Invoke(ctx, "/proto.MinerService/Mine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MinerServiceServer is the server API for MinerService service.
// All implementations should embed UnimplementedMinerServiceServer
// for forward compatibility
type MinerServiceServer interface {
	Mine(context.Context, *MinerConfig) (*MinerResources, error)
}

// UnimplementedMinerServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMinerServiceServer struct {
}

func (UnimplementedMinerServiceServer) Mine(context.Context, *MinerConfig) (*MinerResources, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mine not implemented")
}

// UnsafeMinerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MinerServiceServer will
// result in compilation errors.
type UnsafeMinerServiceServer interface {
	mustEmbedUnimplementedMinerServiceServer()
}

func RegisterMinerServiceServer(s grpc.ServiceRegistrar, srv MinerServiceServer) {
	s.RegisterService(&MinerService_ServiceDesc, srv)
}

func _MinerService_Mine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MinerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServiceServer).Mine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MinerService/Mine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServiceServer).Mine(ctx, req.(*MinerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

// MinerService_ServiceDesc is the grpc.ServiceDesc for MinerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MinerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MinerService",
	HandlerType: (*MinerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Mine",
			Handler:    _MinerService_Mine_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/miner.proto",
}

// PropFormatterClient is the client API for PropFormatter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PropFormatterClient interface {
	Format(ctx context.Context, in *anypb.Any, opts ...grpc.CallOption) (*MinerProperty, error)
}

type propFormatterClient struct {
	cc grpc.ClientConnInterface
}

func NewPropFormatterClient(cc grpc.ClientConnInterface) PropFormatterClient {
	return &propFormatterClient{cc}
}

func (c *propFormatterClient) Format(ctx context.Context, in *anypb.Any, opts ...grpc.CallOption) (*MinerProperty, error) {
	out := new(MinerProperty)
	err := c.cc.Invoke(ctx, "/proto.PropFormatter/Format", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PropFormatterServer is the server API for PropFormatter service.
// All implementations should embed UnimplementedPropFormatterServer
// for forward compatibility
type PropFormatterServer interface {
	Format(context.Context, *anypb.Any) (*MinerProperty, error)
}

// UnimplementedPropFormatterServer should be embedded to have forward compatible implementations.
type UnimplementedPropFormatterServer struct {
}

func (UnimplementedPropFormatterServer) Format(context.Context, *anypb.Any) (*MinerProperty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Format not implemented")
}

// UnsafePropFormatterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PropFormatterServer will
// result in compilation errors.
type UnsafePropFormatterServer interface {
	mustEmbedUnimplementedPropFormatterServer()
}

func RegisterPropFormatterServer(s grpc.ServiceRegistrar, srv PropFormatterServer) {
	s.RegisterService(&PropFormatter_ServiceDesc, srv)
}

func _PropFormatter_Format_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PropFormatterServer).Format(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PropFormatter/Format",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PropFormatterServer).Format(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

// PropFormatter_ServiceDesc is the grpc.ServiceDesc for PropFormatter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PropFormatter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.PropFormatter",
	HandlerType: (*PropFormatterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Format",
			Handler:    _PropFormatter_Format_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/miner.proto",
}
