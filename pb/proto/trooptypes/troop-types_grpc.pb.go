// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/trooptypes/troop-types.proto

package trooptypes

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	V1TroopTypes_GetTroopTypes_FullMethodName   = "/trooptypes.V1TroopTypes/GetTroopTypes"
	V1TroopTypes_GetTroopType_FullMethodName    = "/trooptypes.V1TroopTypes/GetTroopType"
	V1TroopTypes_CreateTroopType_FullMethodName = "/trooptypes.V1TroopTypes/CreateTroopType"
	V1TroopTypes_UpdateTroopType_FullMethodName = "/trooptypes.V1TroopTypes/UpdateTroopType"
	V1TroopTypes_DeleteTroopType_FullMethodName = "/trooptypes.V1TroopTypes/DeleteTroopType"
)

// V1TroopTypesClient is the client API for V1TroopTypes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type V1TroopTypesClient interface {
	GetTroopTypes(ctx context.Context, in *GetTroopTypesRequest, opts ...grpc.CallOption) (*GetTroopTypesReply, error)
	GetTroopType(ctx context.Context, in *GetTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error)
	CreateTroopType(ctx context.Context, in *CreateTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error)
	UpdateTroopType(ctx context.Context, in *UpdateTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error)
	DeleteTroopType(ctx context.Context, in *DeleteTroopTypeRequest, opts ...grpc.CallOption) (*EmptyReply, error)
}

type v1TroopTypesClient struct {
	cc grpc.ClientConnInterface
}

func NewV1TroopTypesClient(cc grpc.ClientConnInterface) V1TroopTypesClient {
	return &v1TroopTypesClient{cc}
}

func (c *v1TroopTypesClient) GetTroopTypes(ctx context.Context, in *GetTroopTypesRequest, opts ...grpc.CallOption) (*GetTroopTypesReply, error) {
	out := new(GetTroopTypesReply)
	err := c.cc.Invoke(ctx, V1TroopTypes_GetTroopTypes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v1TroopTypesClient) GetTroopType(ctx context.Context, in *GetTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error) {
	out := new(TroopType)
	err := c.cc.Invoke(ctx, V1TroopTypes_GetTroopType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v1TroopTypesClient) CreateTroopType(ctx context.Context, in *CreateTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error) {
	out := new(TroopType)
	err := c.cc.Invoke(ctx, V1TroopTypes_CreateTroopType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v1TroopTypesClient) UpdateTroopType(ctx context.Context, in *UpdateTroopTypeRequest, opts ...grpc.CallOption) (*TroopType, error) {
	out := new(TroopType)
	err := c.cc.Invoke(ctx, V1TroopTypes_UpdateTroopType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *v1TroopTypesClient) DeleteTroopType(ctx context.Context, in *DeleteTroopTypeRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, V1TroopTypes_DeleteTroopType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// V1TroopTypesServer is the server API for V1TroopTypes service.
// All implementations must embed UnimplementedV1TroopTypesServer
// for forward compatibility
type V1TroopTypesServer interface {
	GetTroopTypes(context.Context, *GetTroopTypesRequest) (*GetTroopTypesReply, error)
	GetTroopType(context.Context, *GetTroopTypeRequest) (*TroopType, error)
	CreateTroopType(context.Context, *CreateTroopTypeRequest) (*TroopType, error)
	UpdateTroopType(context.Context, *UpdateTroopTypeRequest) (*TroopType, error)
	DeleteTroopType(context.Context, *DeleteTroopTypeRequest) (*EmptyReply, error)
	mustEmbedUnimplementedV1TroopTypesServer()
}

// UnimplementedV1TroopTypesServer must be embedded to have forward compatible implementations.
type UnimplementedV1TroopTypesServer struct {
}

func (UnimplementedV1TroopTypesServer) GetTroopTypes(context.Context, *GetTroopTypesRequest) (*GetTroopTypesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTroopTypes not implemented")
}
func (UnimplementedV1TroopTypesServer) GetTroopType(context.Context, *GetTroopTypeRequest) (*TroopType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTroopType not implemented")
}
func (UnimplementedV1TroopTypesServer) CreateTroopType(context.Context, *CreateTroopTypeRequest) (*TroopType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTroopType not implemented")
}
func (UnimplementedV1TroopTypesServer) UpdateTroopType(context.Context, *UpdateTroopTypeRequest) (*TroopType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTroopType not implemented")
}
func (UnimplementedV1TroopTypesServer) DeleteTroopType(context.Context, *DeleteTroopTypeRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTroopType not implemented")
}
func (UnimplementedV1TroopTypesServer) mustEmbedUnimplementedV1TroopTypesServer() {}

// UnsafeV1TroopTypesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to V1TroopTypesServer will
// result in compilation errors.
type UnsafeV1TroopTypesServer interface {
	mustEmbedUnimplementedV1TroopTypesServer()
}

func RegisterV1TroopTypesServer(s grpc.ServiceRegistrar, srv V1TroopTypesServer) {
	s.RegisterService(&V1TroopTypes_ServiceDesc, srv)
}

func _V1TroopTypes_GetTroopTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTroopTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1TroopTypesServer).GetTroopTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1TroopTypes_GetTroopTypes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1TroopTypesServer).GetTroopTypes(ctx, req.(*GetTroopTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _V1TroopTypes_GetTroopType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTroopTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1TroopTypesServer).GetTroopType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1TroopTypes_GetTroopType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1TroopTypesServer).GetTroopType(ctx, req.(*GetTroopTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _V1TroopTypes_CreateTroopType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTroopTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1TroopTypesServer).CreateTroopType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1TroopTypes_CreateTroopType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1TroopTypesServer).CreateTroopType(ctx, req.(*CreateTroopTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _V1TroopTypes_UpdateTroopType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTroopTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1TroopTypesServer).UpdateTroopType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1TroopTypes_UpdateTroopType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1TroopTypesServer).UpdateTroopType(ctx, req.(*UpdateTroopTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _V1TroopTypes_DeleteTroopType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTroopTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1TroopTypesServer).DeleteTroopType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1TroopTypes_DeleteTroopType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1TroopTypesServer).DeleteTroopType(ctx, req.(*DeleteTroopTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// V1TroopTypes_ServiceDesc is the grpc.ServiceDesc for V1TroopTypes service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var V1TroopTypes_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trooptypes.V1TroopTypes",
	HandlerType: (*V1TroopTypesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTroopTypes",
			Handler:    _V1TroopTypes_GetTroopTypes_Handler,
		},
		{
			MethodName: "GetTroopType",
			Handler:    _V1TroopTypes_GetTroopType_Handler,
		},
		{
			MethodName: "CreateTroopType",
			Handler:    _V1TroopTypes_CreateTroopType_Handler,
		},
		{
			MethodName: "UpdateTroopType",
			Handler:    _V1TroopTypes_UpdateTroopType_Handler,
		},
		{
			MethodName: "DeleteTroopType",
			Handler:    _V1TroopTypes_DeleteTroopType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/trooptypes/troop-types.proto",
}