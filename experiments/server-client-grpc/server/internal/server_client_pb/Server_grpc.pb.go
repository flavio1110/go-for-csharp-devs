// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: Server.proto

package server_client_pb

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

// MessageExchagerServiceClient is the client API for MessageExchagerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageExchagerServiceClient interface {
	Connect(ctx context.Context, opts ...grpc.CallOption) (MessageExchagerService_ConnectClient, error)
}

type messageExchagerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageExchagerServiceClient(cc grpc.ClientConnInterface) MessageExchagerServiceClient {
	return &messageExchagerServiceClient{cc}
}

func (c *messageExchagerServiceClient) Connect(ctx context.Context, opts ...grpc.CallOption) (MessageExchagerService_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageExchagerService_ServiceDesc.Streams[0], "/main_server.MessageExchagerService/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageExchagerServiceConnectClient{stream}
	return x, nil
}

type MessageExchagerService_ConnectClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type messageExchagerServiceConnectClient struct {
	grpc.ClientStream
}

func (x *messageExchagerServiceConnectClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageExchagerServiceConnectClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageExchagerServiceServer is the server API for MessageExchagerService service.
// All implementations must embed UnimplementedMessageExchagerServiceServer
// for forward compatibility
type MessageExchagerServiceServer interface {
	Connect(MessageExchagerService_ConnectServer) error
	mustEmbedUnimplementedMessageExchagerServiceServer()
}

// UnimplementedMessageExchagerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageExchagerServiceServer struct {
}

func (UnimplementedMessageExchagerServiceServer) Connect(MessageExchagerService_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedMessageExchagerServiceServer) mustEmbedUnimplementedMessageExchagerServiceServer() {
}

// UnsafeMessageExchagerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageExchagerServiceServer will
// result in compilation errors.
type UnsafeMessageExchagerServiceServer interface {
	mustEmbedUnimplementedMessageExchagerServiceServer()
}

func RegisterMessageExchagerServiceServer(s grpc.ServiceRegistrar, srv MessageExchagerServiceServer) {
	s.RegisterService(&MessageExchagerService_ServiceDesc, srv)
}

func _MessageExchagerService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageExchagerServiceServer).Connect(&messageExchagerServiceConnectServer{stream})
}

type MessageExchagerService_ConnectServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type messageExchagerServiceConnectServer struct {
	grpc.ServerStream
}

func (x *messageExchagerServiceConnectServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageExchagerServiceConnectServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageExchagerService_ServiceDesc is the grpc.ServiceDesc for MessageExchagerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageExchagerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main_server.MessageExchagerService",
	HandlerType: (*MessageExchagerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _MessageExchagerService_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "Server.proto",
}