// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: cellphone_service.proto

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

// CellphoneServiceClient is the client API for CellphoneService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CellphoneServiceClient interface {
	// Unary RPC
	// 添加一条手机信息
	CreateCellphone(ctx context.Context, in *CreateCellphoneRequest, opts ...grpc.CallOption) (*CreateCellphoneResponse, error)
	// Server streaming RPC
	// 查找符合条件的手机
	SearchCellphone(ctx context.Context, in *FilterCondition, opts ...grpc.CallOption) (CellphoneService_SearchCellphoneClient, error)
	// Client streaming RPC
	// 客户端上传字节流数据（上传手机封面图片）
	UploadCellphoneCover(ctx context.Context, opts ...grpc.CallOption) (CellphoneService_UploadCellphoneCoverClient, error)
	// Bidirectional stream RPC
	// 客户端购买手机，服务端返回购买手机的平均价格
	BuyCellphone(ctx context.Context, opts ...grpc.CallOption) (CellphoneService_BuyCellphoneClient, error)
}

type cellphoneServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCellphoneServiceClient(cc grpc.ClientConnInterface) CellphoneServiceClient {
	return &cellphoneServiceClient{cc}
}

func (c *cellphoneServiceClient) CreateCellphone(ctx context.Context, in *CreateCellphoneRequest, opts ...grpc.CallOption) (*CreateCellphoneResponse, error) {
	out := new(CreateCellphoneResponse)
	err := c.cc.Invoke(ctx, "/pb.CellphoneService/CreateCellphone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cellphoneServiceClient) SearchCellphone(ctx context.Context, in *FilterCondition, opts ...grpc.CallOption) (CellphoneService_SearchCellphoneClient, error) {
	stream, err := c.cc.NewStream(ctx, &CellphoneService_ServiceDesc.Streams[0], "/pb.CellphoneService/SearchCellphone", opts...)
	if err != nil {
		return nil, err
	}
	x := &cellphoneServiceSearchCellphoneClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CellphoneService_SearchCellphoneClient interface {
	Recv() (*Cellphone, error)
	grpc.ClientStream
}

type cellphoneServiceSearchCellphoneClient struct {
	grpc.ClientStream
}

func (x *cellphoneServiceSearchCellphoneClient) Recv() (*Cellphone, error) {
	m := new(Cellphone)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cellphoneServiceClient) UploadCellphoneCover(ctx context.Context, opts ...grpc.CallOption) (CellphoneService_UploadCellphoneCoverClient, error) {
	stream, err := c.cc.NewStream(ctx, &CellphoneService_ServiceDesc.Streams[1], "/pb.CellphoneService/UploadCellphoneCover", opts...)
	if err != nil {
		return nil, err
	}
	x := &cellphoneServiceUploadCellphoneCoverClient{stream}
	return x, nil
}

type CellphoneService_UploadCellphoneCoverClient interface {
	Send(*UploadCellphoneCoverRequest) error
	CloseAndRecv() (*UploadCellphoneCoverResponse, error)
	grpc.ClientStream
}

type cellphoneServiceUploadCellphoneCoverClient struct {
	grpc.ClientStream
}

func (x *cellphoneServiceUploadCellphoneCoverClient) Send(m *UploadCellphoneCoverRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cellphoneServiceUploadCellphoneCoverClient) CloseAndRecv() (*UploadCellphoneCoverResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadCellphoneCoverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cellphoneServiceClient) BuyCellphone(ctx context.Context, opts ...grpc.CallOption) (CellphoneService_BuyCellphoneClient, error) {
	stream, err := c.cc.NewStream(ctx, &CellphoneService_ServiceDesc.Streams[2], "/pb.CellphoneService/BuyCellphone", opts...)
	if err != nil {
		return nil, err
	}
	x := &cellphoneServiceBuyCellphoneClient{stream}
	return x, nil
}

type CellphoneService_BuyCellphoneClient interface {
	Send(*BuyCellphoneRequest) error
	Recv() (*BuyCellphoneResponse, error)
	grpc.ClientStream
}

type cellphoneServiceBuyCellphoneClient struct {
	grpc.ClientStream
}

func (x *cellphoneServiceBuyCellphoneClient) Send(m *BuyCellphoneRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cellphoneServiceBuyCellphoneClient) Recv() (*BuyCellphoneResponse, error) {
	m := new(BuyCellphoneResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CellphoneServiceServer is the server API for CellphoneService service.
// All implementations must embed UnimplementedCellphoneServiceServer
// for forward compatibility
type CellphoneServiceServer interface {
	// Unary RPC
	// 添加一条手机信息
	CreateCellphone(context.Context, *CreateCellphoneRequest) (*CreateCellphoneResponse, error)
	// Server streaming RPC
	// 查找符合条件的手机
	SearchCellphone(*FilterCondition, CellphoneService_SearchCellphoneServer) error
	// Client streaming RPC
	// 客户端上传字节流数据（上传手机封面图片）
	UploadCellphoneCover(CellphoneService_UploadCellphoneCoverServer) error
	// Bidirectional stream RPC
	// 客户端购买手机，服务端返回购买手机的平均价格
	BuyCellphone(CellphoneService_BuyCellphoneServer) error
	mustEmbedUnimplementedCellphoneServiceServer()
}

// UnimplementedCellphoneServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCellphoneServiceServer struct {
}

func (UnimplementedCellphoneServiceServer) CreateCellphone(context.Context, *CreateCellphoneRequest) (*CreateCellphoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCellphone not implemented")
}
func (UnimplementedCellphoneServiceServer) SearchCellphone(*FilterCondition, CellphoneService_SearchCellphoneServer) error {
	return status.Errorf(codes.Unimplemented, "method SearchCellphone not implemented")
}
func (UnimplementedCellphoneServiceServer) UploadCellphoneCover(CellphoneService_UploadCellphoneCoverServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadCellphoneCover not implemented")
}
func (UnimplementedCellphoneServiceServer) BuyCellphone(CellphoneService_BuyCellphoneServer) error {
	return status.Errorf(codes.Unimplemented, "method BuyCellphone not implemented")
}
func (UnimplementedCellphoneServiceServer) mustEmbedUnimplementedCellphoneServiceServer() {}

// UnsafeCellphoneServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CellphoneServiceServer will
// result in compilation errors.
type UnsafeCellphoneServiceServer interface {
	mustEmbedUnimplementedCellphoneServiceServer()
}

func RegisterCellphoneServiceServer(s grpc.ServiceRegistrar, srv CellphoneServiceServer) {
	s.RegisterService(&CellphoneService_ServiceDesc, srv)
}

func _CellphoneService_CreateCellphone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCellphoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellphoneServiceServer).CreateCellphone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CellphoneService/CreateCellphone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellphoneServiceServer).CreateCellphone(ctx, req.(*CreateCellphoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CellphoneService_SearchCellphone_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FilterCondition)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CellphoneServiceServer).SearchCellphone(m, &cellphoneServiceSearchCellphoneServer{stream})
}

type CellphoneService_SearchCellphoneServer interface {
	Send(*Cellphone) error
	grpc.ServerStream
}

type cellphoneServiceSearchCellphoneServer struct {
	grpc.ServerStream
}

func (x *cellphoneServiceSearchCellphoneServer) Send(m *Cellphone) error {
	return x.ServerStream.SendMsg(m)
}

func _CellphoneService_UploadCellphoneCover_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CellphoneServiceServer).UploadCellphoneCover(&cellphoneServiceUploadCellphoneCoverServer{stream})
}

type CellphoneService_UploadCellphoneCoverServer interface {
	SendAndClose(*UploadCellphoneCoverResponse) error
	Recv() (*UploadCellphoneCoverRequest, error)
	grpc.ServerStream
}

type cellphoneServiceUploadCellphoneCoverServer struct {
	grpc.ServerStream
}

func (x *cellphoneServiceUploadCellphoneCoverServer) SendAndClose(m *UploadCellphoneCoverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cellphoneServiceUploadCellphoneCoverServer) Recv() (*UploadCellphoneCoverRequest, error) {
	m := new(UploadCellphoneCoverRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CellphoneService_BuyCellphone_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CellphoneServiceServer).BuyCellphone(&cellphoneServiceBuyCellphoneServer{stream})
}

type CellphoneService_BuyCellphoneServer interface {
	Send(*BuyCellphoneResponse) error
	Recv() (*BuyCellphoneRequest, error)
	grpc.ServerStream
}

type cellphoneServiceBuyCellphoneServer struct {
	grpc.ServerStream
}

func (x *cellphoneServiceBuyCellphoneServer) Send(m *BuyCellphoneResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cellphoneServiceBuyCellphoneServer) Recv() (*BuyCellphoneRequest, error) {
	m := new(BuyCellphoneRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CellphoneService_ServiceDesc is the grpc.ServiceDesc for CellphoneService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CellphoneService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CellphoneService",
	HandlerType: (*CellphoneServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCellphone",
			Handler:    _CellphoneService_CreateCellphone_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SearchCellphone",
			Handler:       _CellphoneService_SearchCellphone_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadCellphoneCover",
			Handler:       _CellphoneService_UploadCellphoneCover_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BuyCellphone",
			Handler:       _CellphoneService_BuyCellphone_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "cellphone_service.proto",
}
