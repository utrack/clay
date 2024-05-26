// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sum

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

// SummatorClient is the client API for Summator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SummatorClient interface {
	Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error)
}

type summatorClient struct {
	cc grpc.ClientConnInterface
}

func NewSummatorClient(cc grpc.ClientConnInterface) SummatorClient {
	return &summatorClient{cc}
}

func (c *summatorClient) Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error) {
	out := new(SumResponse)
	err := c.cc.Invoke(ctx, "/sumpb.Summator/Sum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SummatorServer is the server API for Summator service.
// All implementations must embed UnimplementedSummatorServer
// for forward compatibility
type SummatorServer interface {
	Sum(context.Context, *SumRequest) (*SumResponse, error)
	mustEmbedUnimplementedSummatorServer()
}

// UnimplementedSummatorServer must be embedded to have forward compatible implementations.
type UnimplementedSummatorServer struct {
}

func (UnimplementedSummatorServer) Sum(context.Context, *SumRequest) (*SumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (UnimplementedSummatorServer) mustEmbedUnimplementedSummatorServer() {}

// UnsafeSummatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SummatorServer will
// result in compilation errors.
type UnsafeSummatorServer interface {
	mustEmbedUnimplementedSummatorServer()
}

func RegisterSummatorServer(s grpc.ServiceRegistrar, srv SummatorServer) {
	s.RegisterService(&Summator_ServiceDesc, srv)
}

func _Summator_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SummatorServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sumpb.Summator/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SummatorServer).Sum(ctx, req.(*SumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Summator_ServiceDesc is the grpc.ServiceDesc for Summator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Summator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sumpb.Summator",
	HandlerType: (*SummatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _Summator_Sum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/sum.proto",
}
