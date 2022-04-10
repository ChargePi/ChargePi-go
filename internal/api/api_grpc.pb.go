// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// ChargePointClient is the client API for ChargePoint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChargePointClient interface {
	GetConnectorStatus(ctx context.Context, opts ...grpc.CallOption) (ChargePoint_GetConnectorStatusClient, error)
	StartTransaction(ctx context.Context, in *StartTransactionRequest, opts ...grpc.CallOption) (*StartTransactionResponse, error)
	StopTransaction(ctx context.Context, in *StopTransactionRequest, opts ...grpc.CallOption) (*StopTransactionResponse, error)
	HandleCharging(ctx context.Context, in *HandleChargingRequest, opts ...grpc.CallOption) (*HandleChargingResponse, error)
}

type chargePointClient struct {
	cc grpc.ClientConnInterface
}

func NewChargePointClient(cc grpc.ClientConnInterface) ChargePointClient {
	return &chargePointClient{cc}
}

func (c *chargePointClient) GetConnectorStatus(ctx context.Context, opts ...grpc.CallOption) (ChargePoint_GetConnectorStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChargePoint_ServiceDesc.Streams[0], "/api.ChargePoint/GetConnectorStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &chargePointGetConnectorStatusClient{stream}
	return x, nil
}

type ChargePoint_GetConnectorStatusClient interface {
	Send(*GetConnectorStatusRequest) error
	Recv() (*GetConnectorStatusResponse, error)
	grpc.ClientStream
}

type chargePointGetConnectorStatusClient struct {
	grpc.ClientStream
}

func (x *chargePointGetConnectorStatusClient) Send(m *GetConnectorStatusRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chargePointGetConnectorStatusClient) Recv() (*GetConnectorStatusResponse, error) {
	m := new(GetConnectorStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chargePointClient) StartTransaction(ctx context.Context, in *StartTransactionRequest, opts ...grpc.CallOption) (*StartTransactionResponse, error) {
	out := new(StartTransactionResponse)
	err := c.cc.Invoke(ctx, "/api.ChargePoint/StartTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chargePointClient) StopTransaction(ctx context.Context, in *StopTransactionRequest, opts ...grpc.CallOption) (*StopTransactionResponse, error) {
	out := new(StopTransactionResponse)
	err := c.cc.Invoke(ctx, "/api.ChargePoint/StopTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chargePointClient) HandleCharging(ctx context.Context, in *HandleChargingRequest, opts ...grpc.CallOption) (*HandleChargingResponse, error) {
	out := new(HandleChargingResponse)
	err := c.cc.Invoke(ctx, "/api.ChargePoint/HandleCharging", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChargePointServer is the server API for ChargePoint service.
// All implementations must embed UnimplementedChargePointServer
// for forward compatibility
type ChargePointServer interface {
	GetConnectorStatus(ChargePoint_GetConnectorStatusServer) error
	StartTransaction(context.Context, *StartTransactionRequest) (*StartTransactionResponse, error)
	StopTransaction(context.Context, *StopTransactionRequest) (*StopTransactionResponse, error)
	HandleCharging(context.Context, *HandleChargingRequest) (*HandleChargingResponse, error)
	mustEmbedUnimplementedChargePointServer()
}

// UnimplementedChargePointServer must be embedded to have forward compatible implementations.
type UnimplementedChargePointServer struct {
}

func (UnimplementedChargePointServer) GetConnectorStatus(ChargePoint_GetConnectorStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method GetConnectorStatus not implemented")
}
func (UnimplementedChargePointServer) StartTransaction(context.Context, *StartTransactionRequest) (*StartTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTransaction not implemented")
}
func (UnimplementedChargePointServer) StopTransaction(context.Context, *StopTransactionRequest) (*StopTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTransaction not implemented")
}
func (UnimplementedChargePointServer) HandleCharging(context.Context, *HandleChargingRequest) (*HandleChargingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleCharging not implemented")
}
func (UnimplementedChargePointServer) mustEmbedUnimplementedChargePointServer() {}

// UnsafeChargePointServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChargePointServer will
// result in compilation errors.
type UnsafeChargePointServer interface {
	mustEmbedUnimplementedChargePointServer()
}

func RegisterChargePointServer(s grpc.ServiceRegistrar, srv ChargePointServer) {
	s.RegisterService(&ChargePoint_ServiceDesc, srv)
}

func _ChargePoint_GetConnectorStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChargePointServer).GetConnectorStatus(&chargePointGetConnectorStatusServer{stream})
}

type ChargePoint_GetConnectorStatusServer interface {
	Send(*GetConnectorStatusResponse) error
	Recv() (*GetConnectorStatusRequest, error)
	grpc.ServerStream
}

type chargePointGetConnectorStatusServer struct {
	grpc.ServerStream
}

func (x *chargePointGetConnectorStatusServer) Send(m *GetConnectorStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chargePointGetConnectorStatusServer) Recv() (*GetConnectorStatusRequest, error) {
	m := new(GetConnectorStatusRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChargePoint_StartTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChargePointServer).StartTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ChargePoint/StartTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChargePointServer).StartTransaction(ctx, req.(*StartTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChargePoint_StopTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChargePointServer).StopTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ChargePoint/StopTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChargePointServer).StopTransaction(ctx, req.(*StopTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChargePoint_HandleCharging_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleChargingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChargePointServer).HandleCharging(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ChargePoint/HandleCharging",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChargePointServer).HandleCharging(ctx, req.(*HandleChargingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChargePoint_ServiceDesc is the grpc.ServiceDesc for ChargePoint service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChargePoint_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ChargePoint",
	HandlerType: (*ChargePointServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartTransaction",
			Handler:    _ChargePoint_StartTransaction_Handler,
		},
		{
			MethodName: "StopTransaction",
			Handler:    _ChargePoint_StopTransaction_Handler,
		},
		{
			MethodName: "HandleCharging",
			Handler:    _ChargePoint_HandleCharging_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetConnectorStatus",
			Handler:       _ChargePoint_GetConnectorStatus_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internal/models/api/api.proto",
}