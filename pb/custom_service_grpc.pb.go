// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: custom_service.proto

package pb

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

// CustomServiceClient is the client API for CustomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomServiceClient interface {
	MetadataCarryTest(ctx context.Context, in *CustomRequest, opts ...grpc.CallOption) (*CustomResponse, error)
}

type customServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomServiceClient(cc grpc.ClientConnInterface) CustomServiceClient {
	return &customServiceClient{cc}
}

func (c *customServiceClient) MetadataCarryTest(ctx context.Context, in *CustomRequest, opts ...grpc.CallOption) (*CustomResponse, error) {
	out := new(CustomResponse)
	err := c.cc.Invoke(ctx, "/pb.CustomService/MetadataCarryTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomServiceServer is the server API for CustomService service.
// All implementations must embed UnimplementedCustomServiceServer
// for forward compatibility
type CustomServiceServer interface {
	MetadataCarryTest(context.Context, *CustomRequest) (*CustomResponse, error)
	mustEmbedUnimplementedCustomServiceServer()
}

// UnimplementedCustomServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCustomServiceServer struct {
}

func (UnimplementedCustomServiceServer) MetadataCarryTest(context.Context, *CustomRequest) (*CustomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MetadataCarryTest not implemented")
}
func (UnimplementedCustomServiceServer) mustEmbedUnimplementedCustomServiceServer() {}

// UnsafeCustomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomServiceServer will
// result in compilation errors.
type UnsafeCustomServiceServer interface {
	mustEmbedUnimplementedCustomServiceServer()
}

func RegisterCustomServiceServer(s grpc.ServiceRegistrar, srv CustomServiceServer) {
	s.RegisterService(&CustomService_ServiceDesc, srv)
}

func _CustomService_MetadataCarryTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomServiceServer).MetadataCarryTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CustomService/MetadataCarryTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomServiceServer).MetadataCarryTest(ctx, req.(*CustomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CustomService_ServiceDesc is the grpc.ServiceDesc for CustomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CustomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CustomService",
	HandlerType: (*CustomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MetadataCarryTest",
			Handler:    _CustomService_MetadataCarryTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "custom_service.proto",
}