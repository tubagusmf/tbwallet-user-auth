// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: pb/kycdoc/kycdoc.proto

package kycdoc

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
	KycdocService_GetKycdocByUserID_FullMethodName = "/kycdoc.KycdocService/GetKycdocByUserID"
	KycdocService_GetKycStatus_FullMethodName      = "/kycdoc.KycdocService/GetKycStatus"
)

// KycdocServiceClient is the client API for KycdocService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KycdocServiceClient interface {
	GetKycdocByUserID(ctx context.Context, in *GetKycdocByUserIDRequest, opts ...grpc.CallOption) (*GetKycdocByUserIDResponse, error)
	GetKycStatus(ctx context.Context, in *GetKycStatusRequest, opts ...grpc.CallOption) (*GetKycStatusResponse, error)
}

type kycdocServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKycdocServiceClient(cc grpc.ClientConnInterface) KycdocServiceClient {
	return &kycdocServiceClient{cc}
}

func (c *kycdocServiceClient) GetKycdocByUserID(ctx context.Context, in *GetKycdocByUserIDRequest, opts ...grpc.CallOption) (*GetKycdocByUserIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKycdocByUserIDResponse)
	err := c.cc.Invoke(ctx, KycdocService_GetKycdocByUserID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kycdocServiceClient) GetKycStatus(ctx context.Context, in *GetKycStatusRequest, opts ...grpc.CallOption) (*GetKycStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKycStatusResponse)
	err := c.cc.Invoke(ctx, KycdocService_GetKycStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KycdocServiceServer is the server API for KycdocService service.
// All implementations must embed UnimplementedKycdocServiceServer
// for forward compatibility.
type KycdocServiceServer interface {
	GetKycdocByUserID(context.Context, *GetKycdocByUserIDRequest) (*GetKycdocByUserIDResponse, error)
	GetKycStatus(context.Context, *GetKycStatusRequest) (*GetKycStatusResponse, error)
	mustEmbedUnimplementedKycdocServiceServer()
}

// UnimplementedKycdocServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKycdocServiceServer struct{}

func (UnimplementedKycdocServiceServer) GetKycdocByUserID(context.Context, *GetKycdocByUserIDRequest) (*GetKycdocByUserIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKycdocByUserID not implemented")
}
func (UnimplementedKycdocServiceServer) GetKycStatus(context.Context, *GetKycStatusRequest) (*GetKycStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKycStatus not implemented")
}
func (UnimplementedKycdocServiceServer) mustEmbedUnimplementedKycdocServiceServer() {}
func (UnimplementedKycdocServiceServer) testEmbeddedByValue()                       {}

// UnsafeKycdocServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KycdocServiceServer will
// result in compilation errors.
type UnsafeKycdocServiceServer interface {
	mustEmbedUnimplementedKycdocServiceServer()
}

func RegisterKycdocServiceServer(s grpc.ServiceRegistrar, srv KycdocServiceServer) {
	// If the following call pancis, it indicates UnimplementedKycdocServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KycdocService_ServiceDesc, srv)
}

func _KycdocService_GetKycdocByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKycdocByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KycdocServiceServer).GetKycdocByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KycdocService_GetKycdocByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KycdocServiceServer).GetKycdocByUserID(ctx, req.(*GetKycdocByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KycdocService_GetKycStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKycStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KycdocServiceServer).GetKycStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KycdocService_GetKycStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KycdocServiceServer).GetKycStatus(ctx, req.(*GetKycStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KycdocService_ServiceDesc is the grpc.ServiceDesc for KycdocService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KycdocService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kycdoc.KycdocService",
	HandlerType: (*KycdocServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKycdocByUserID",
			Handler:    _KycdocService_GetKycdocByUserID_Handler,
		},
		{
			MethodName: "GetKycStatus",
			Handler:    _KycdocService_GetKycStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/kycdoc/kycdoc.proto",
}
