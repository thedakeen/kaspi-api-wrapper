// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: device/device.proto

package kaspiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DeviceService_GetTradePoints_FullMethodName         = "/kaspi.api.v1.DeviceService/GetTradePoints"
	DeviceService_RegisterDevice_FullMethodName         = "/kaspi.api.v1.DeviceService/RegisterDevice"
	DeviceService_DeleteDevice_FullMethodName           = "/kaspi.api.v1.DeviceService/DeleteDevice"
	DeviceService_GetTradePointsEnhanced_FullMethodName = "/kaspi.api.v1.DeviceService/GetTradePointsEnhanced"
	DeviceService_RegisterDeviceEnhanced_FullMethodName = "/kaspi.api.v1.DeviceService/RegisterDeviceEnhanced"
	DeviceService_DeleteDeviceEnhanced_FullMethodName   = "/kaspi.api.v1.DeviceService/DeleteDeviceEnhanced"
)

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	// Basic/Standard scheme methods
	GetTradePoints(ctx context.Context, in *GetTradePointsRequest, opts ...grpc.CallOption) (*GetTradePointsResponse, error)
	RegisterDevice(ctx context.Context, in *RegisterDeviceRequest, opts ...grpc.CallOption) (*RegisterDeviceResponse, error)
	DeleteDevice(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*DeleteDeviceResponse, error)
	// Enhanced scheme methods
	GetTradePointsEnhanced(ctx context.Context, in *GetTradePointsEnhancedRequest, opts ...grpc.CallOption) (*GetTradePointsResponse, error)
	RegisterDeviceEnhanced(ctx context.Context, in *RegisterDeviceEnhancedRequest, opts ...grpc.CallOption) (*RegisterDeviceResponse, error)
	DeleteDeviceEnhanced(ctx context.Context, in *DeleteDeviceEnhancedRequest, opts ...grpc.CallOption) (*DeleteDeviceResponse, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) GetTradePoints(ctx context.Context, in *GetTradePointsRequest, opts ...grpc.CallOption) (*GetTradePointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTradePointsResponse)
	err := c.cc.Invoke(ctx, DeviceService_GetTradePoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) RegisterDevice(ctx context.Context, in *RegisterDeviceRequest, opts ...grpc.CallOption) (*RegisterDeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterDeviceResponse)
	err := c.cc.Invoke(ctx, DeviceService_RegisterDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) DeleteDevice(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*DeleteDeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDeviceResponse)
	err := c.cc.Invoke(ctx, DeviceService_DeleteDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetTradePointsEnhanced(ctx context.Context, in *GetTradePointsEnhancedRequest, opts ...grpc.CallOption) (*GetTradePointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTradePointsResponse)
	err := c.cc.Invoke(ctx, DeviceService_GetTradePointsEnhanced_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) RegisterDeviceEnhanced(ctx context.Context, in *RegisterDeviceEnhancedRequest, opts ...grpc.CallOption) (*RegisterDeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterDeviceResponse)
	err := c.cc.Invoke(ctx, DeviceService_RegisterDeviceEnhanced_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) DeleteDeviceEnhanced(ctx context.Context, in *DeleteDeviceEnhancedRequest, opts ...grpc.CallOption) (*DeleteDeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDeviceResponse)
	err := c.cc.Invoke(ctx, DeviceService_DeleteDeviceEnhanced_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility.
type DeviceServiceServer interface {
	// Basic/Standard scheme methods
	GetTradePoints(context.Context, *GetTradePointsRequest) (*GetTradePointsResponse, error)
	RegisterDevice(context.Context, *RegisterDeviceRequest) (*RegisterDeviceResponse, error)
	DeleteDevice(context.Context, *DeleteDeviceRequest) (*DeleteDeviceResponse, error)
	// Enhanced scheme methods
	GetTradePointsEnhanced(context.Context, *GetTradePointsEnhancedRequest) (*GetTradePointsResponse, error)
	RegisterDeviceEnhanced(context.Context, *RegisterDeviceEnhancedRequest) (*RegisterDeviceResponse, error)
	DeleteDeviceEnhanced(context.Context, *DeleteDeviceEnhancedRequest) (*DeleteDeviceResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDeviceServiceServer struct{}

func (UnimplementedDeviceServiceServer) GetTradePoints(context.Context, *GetTradePointsRequest) (*GetTradePointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTradePoints not implemented")
}
func (UnimplementedDeviceServiceServer) RegisterDevice(context.Context, *RegisterDeviceRequest) (*RegisterDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDevice not implemented")
}
func (UnimplementedDeviceServiceServer) DeleteDevice(context.Context, *DeleteDeviceRequest) (*DeleteDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDevice not implemented")
}
func (UnimplementedDeviceServiceServer) GetTradePointsEnhanced(context.Context, *GetTradePointsEnhancedRequest) (*GetTradePointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTradePointsEnhanced not implemented")
}
func (UnimplementedDeviceServiceServer) RegisterDeviceEnhanced(context.Context, *RegisterDeviceEnhancedRequest) (*RegisterDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDeviceEnhanced not implemented")
}
func (UnimplementedDeviceServiceServer) DeleteDeviceEnhanced(context.Context, *DeleteDeviceEnhancedRequest) (*DeleteDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDeviceEnhanced not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}
func (UnimplementedDeviceServiceServer) testEmbeddedByValue()                       {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	// If the following call pancis, it indicates UnimplementedDeviceServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_GetTradePoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTradePointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetTradePoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_GetTradePoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetTradePoints(ctx, req.(*GetTradePointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_RegisterDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).RegisterDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_RegisterDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).RegisterDevice(ctx, req.(*RegisterDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_DeleteDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).DeleteDevice(ctx, req.(*DeleteDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetTradePointsEnhanced_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTradePointsEnhancedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetTradePointsEnhanced(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_GetTradePointsEnhanced_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetTradePointsEnhanced(ctx, req.(*GetTradePointsEnhancedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_RegisterDeviceEnhanced_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterDeviceEnhancedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).RegisterDeviceEnhanced(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_RegisterDeviceEnhanced_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).RegisterDeviceEnhanced(ctx, req.(*RegisterDeviceEnhancedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_DeleteDeviceEnhanced_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceEnhancedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).DeleteDeviceEnhanced(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_DeleteDeviceEnhanced_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).DeleteDeviceEnhanced(ctx, req.(*DeleteDeviceEnhancedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kaspi.api.v1.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTradePoints",
			Handler:    _DeviceService_GetTradePoints_Handler,
		},
		{
			MethodName: "RegisterDevice",
			Handler:    _DeviceService_RegisterDevice_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _DeviceService_DeleteDevice_Handler,
		},
		{
			MethodName: "GetTradePointsEnhanced",
			Handler:    _DeviceService_GetTradePointsEnhanced_Handler,
		},
		{
			MethodName: "RegisterDeviceEnhanced",
			Handler:    _DeviceService_RegisterDeviceEnhanced_Handler,
		},
		{
			MethodName: "DeleteDeviceEnhanced",
			Handler:    _DeviceService_DeleteDeviceEnhanced_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "device/device.proto",
}
