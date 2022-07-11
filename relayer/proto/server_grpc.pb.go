// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// RelayHubClient is the client API for RelayHub service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelayHubClient interface {
	SignCheckpointProof(ctx context.Context, in *SignCheckpointProofRequest, opts ...grpc.CallOption) (RelayHub_SignCheckpointProofClient, error)
}

type relayHubClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayHubClient(cc grpc.ClientConnInterface) RelayHubClient {
	return &relayHubClient{cc}
}

func (c *relayHubClient) SignCheckpointProof(ctx context.Context, in *SignCheckpointProofRequest, opts ...grpc.CallOption) (RelayHub_SignCheckpointProofClient, error) {
	stream, err := c.cc.NewStream(ctx, &RelayHub_ServiceDesc.Streams[0], "/com.ankr.RelayHub/SignCheckpointProof", opts...)
	if err != nil {
		return nil, err
	}
	x := &relayHubSignCheckpointProofClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RelayHub_SignCheckpointProofClient interface {
	Recv() (*SignCheckpointProofReply, error)
	grpc.ClientStream
}

type relayHubSignCheckpointProofClient struct {
	grpc.ClientStream
}

func (x *relayHubSignCheckpointProofClient) Recv() (*SignCheckpointProofReply, error) {
	m := new(SignCheckpointProofReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RelayHubServer is the server API for RelayHub service.
// All implementations must embed UnimplementedRelayHubServer
// for forward compatibility
type RelayHubServer interface {
	SignCheckpointProof(*SignCheckpointProofRequest, RelayHub_SignCheckpointProofServer) error
	mustEmbedUnimplementedRelayHubServer()
}

// UnimplementedRelayHubServer must be embedded to have forward compatible implementations.
type UnimplementedRelayHubServer struct {
}

func (UnimplementedRelayHubServer) SignCheckpointProof(*SignCheckpointProofRequest, RelayHub_SignCheckpointProofServer) error {
	return status.Errorf(codes.Unimplemented, "method SignCheckpointProof not implemented")
}
func (UnimplementedRelayHubServer) mustEmbedUnimplementedRelayHubServer() {}

// UnsafeRelayHubServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelayHubServer will
// result in compilation errors.
type UnsafeRelayHubServer interface {
	mustEmbedUnimplementedRelayHubServer()
}

func RegisterRelayHubServer(s grpc.ServiceRegistrar, srv RelayHubServer) {
	s.RegisterService(&RelayHub_ServiceDesc, srv)
}

func _RelayHub_SignCheckpointProof_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SignCheckpointProofRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RelayHubServer).SignCheckpointProof(m, &relayHubSignCheckpointProofServer{stream})
}

type RelayHub_SignCheckpointProofServer interface {
	Send(*SignCheckpointProofReply) error
	grpc.ServerStream
}

type relayHubSignCheckpointProofServer struct {
	grpc.ServerStream
}

func (x *relayHubSignCheckpointProofServer) Send(m *SignCheckpointProofReply) error {
	return x.ServerStream.SendMsg(m)
}

// RelayHub_ServiceDesc is the grpc.ServiceDesc for RelayHub service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelayHub_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.ankr.RelayHub",
	HandlerType: (*RelayHubServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SignCheckpointProof",
			Handler:       _RelayHub_SignCheckpointProof_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "relayer/proto/server.proto",
}