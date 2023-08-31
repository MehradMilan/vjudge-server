// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: api/judge.proto

package judge

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
	CodeJudge_JudgeCode_FullMethodName = "/judge.CodeJudge/JudgeCode"
)

// CodeJudgeClient is the client API for CodeJudge service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeJudgeClient interface {
	JudgeCode(ctx context.Context, in *SubmissionRequest, opts ...grpc.CallOption) (*JudgementReply, error)
}

type codeJudgeClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeJudgeClient(cc grpc.ClientConnInterface) CodeJudgeClient {
	return &codeJudgeClient{cc}
}

func (c *codeJudgeClient) JudgeCode(ctx context.Context, in *SubmissionRequest, opts ...grpc.CallOption) (*JudgementReply, error) {
	out := new(JudgementReply)
	err := c.cc.Invoke(ctx, CodeJudge_JudgeCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeJudgeServer is the server API for CodeJudge service.
// All implementations must embed UnimplementedCodeJudgeServer
// for forward compatibility
type CodeJudgeServer interface {
	JudgeCode(context.Context, *SubmissionRequest) (*JudgementReply, error)
	mustEmbedUnimplementedCodeJudgeServer()
}

// UnimplementedCodeJudgeServer must be embedded to have forward compatible implementations.
type UnimplementedCodeJudgeServer struct {
}

func (UnimplementedCodeJudgeServer) JudgeCode(context.Context, *SubmissionRequest) (*JudgementReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JudgeCode not implemented")
}
func (UnimplementedCodeJudgeServer) mustEmbedUnimplementedCodeJudgeServer() {}

// UnsafeCodeJudgeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeJudgeServer will
// result in compilation errors.
type UnsafeCodeJudgeServer interface {
	mustEmbedUnimplementedCodeJudgeServer()
}

func RegisterCodeJudgeServer(s grpc.ServiceRegistrar, srv CodeJudgeServer) {
	s.RegisterService(&CodeJudge_ServiceDesc, srv)
}

func _CodeJudge_JudgeCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeJudgeServer).JudgeCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeJudge_JudgeCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeJudgeServer).JudgeCode(ctx, req.(*SubmissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CodeJudge_ServiceDesc is the grpc.ServiceDesc for CodeJudge service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CodeJudge_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "judge.CodeJudge",
	HandlerType: (*CodeJudgeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JudgeCode",
			Handler:    _CodeJudge_JudgeCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/judge.proto",
}
