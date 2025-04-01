// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: posts.proto

package main

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Posts_CreatePost_FullMethodName    = "/posts.Posts/CreatePost"
	Posts_CreateComment_FullMethodName = "/posts.Posts/CreateComment"
	Posts_ModifyPost_FullMethodName    = "/posts.Posts/ModifyPost"
	Posts_ModifyComment_FullMethodName = "/posts.Posts/ModifyComment"
	Posts_DeletePost_FullMethodName    = "/posts.Posts/DeletePost"
	Posts_DeleteComment_FullMethodName = "/posts.Posts/DeleteComment"
	Posts_GetPosts_FullMethodName      = "/posts.Posts/GetPosts"
	Posts_GetComments_FullMethodName   = "/posts.Posts/GetComments"
)

// PostsClient is the client API for Posts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostsClient interface {
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ModifyPost(ctx context.Context, in *ModifyPostRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ModifyComment(ctx context.Context, in *ModifyCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetPosts(ctx context.Context, in *GetPostsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PostInfo], error)
	GetComments(ctx context.Context, in *GetCommentsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CommentInfo], error)
}

type postsClient struct {
	cc grpc.ClientConnInterface
}

func NewPostsClient(cc grpc.ClientConnInterface) PostsClient {
	return &postsClient{cc}
}

func (c *postsClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_CreatePost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_CreateComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) ModifyPost(ctx context.Context, in *ModifyPostRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_ModifyPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) ModifyComment(ctx context.Context, in *ModifyCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_ModifyComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_DeletePost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Posts_DeleteComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postsClient) GetPosts(ctx context.Context, in *GetPostsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[PostInfo], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Posts_ServiceDesc.Streams[0], Posts_GetPosts_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetPostsRequest, PostInfo]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Posts_GetPostsClient = grpc.ServerStreamingClient[PostInfo]

func (c *postsClient) GetComments(ctx context.Context, in *GetCommentsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CommentInfo], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Posts_ServiceDesc.Streams[1], Posts_GetComments_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetCommentsRequest, CommentInfo]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Posts_GetCommentsClient = grpc.ServerStreamingClient[CommentInfo]

// PostsServer is the server API for Posts service.
// All implementations must embed UnimplementedPostsServer
// for forward compatibility.
type PostsServer interface {
	CreatePost(context.Context, *CreatePostRequest) (*empty.Empty, error)
	CreateComment(context.Context, *CreateCommentRequest) (*empty.Empty, error)
	ModifyPost(context.Context, *ModifyPostRequest) (*empty.Empty, error)
	ModifyComment(context.Context, *ModifyCommentRequest) (*empty.Empty, error)
	DeletePost(context.Context, *DeletePostRequest) (*empty.Empty, error)
	DeleteComment(context.Context, *DeleteCommentRequest) (*empty.Empty, error)
	GetPosts(*GetPostsRequest, grpc.ServerStreamingServer[PostInfo]) error
	GetComments(*GetCommentsRequest, grpc.ServerStreamingServer[CommentInfo]) error
	mustEmbedUnimplementedPostsServer()
}

// UnimplementedPostsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPostsServer struct{}

func (UnimplementedPostsServer) CreatePost(context.Context, *CreatePostRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedPostsServer) CreateComment(context.Context, *CreateCommentRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedPostsServer) ModifyPost(context.Context, *ModifyPostRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyPost not implemented")
}
func (UnimplementedPostsServer) ModifyComment(context.Context, *ModifyCommentRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyComment not implemented")
}
func (UnimplementedPostsServer) DeletePost(context.Context, *DeletePostRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (UnimplementedPostsServer) DeleteComment(context.Context, *DeleteCommentRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (UnimplementedPostsServer) GetPosts(*GetPostsRequest, grpc.ServerStreamingServer[PostInfo]) error {
	return status.Errorf(codes.Unimplemented, "method GetPosts not implemented")
}
func (UnimplementedPostsServer) GetComments(*GetCommentsRequest, grpc.ServerStreamingServer[CommentInfo]) error {
	return status.Errorf(codes.Unimplemented, "method GetComments not implemented")
}
func (UnimplementedPostsServer) mustEmbedUnimplementedPostsServer() {}
func (UnimplementedPostsServer) testEmbeddedByValue()               {}

// UnsafePostsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostsServer will
// result in compilation errors.
type UnsafePostsServer interface {
	mustEmbedUnimplementedPostsServer()
}

func RegisterPostsServer(s grpc.ServiceRegistrar, srv PostsServer) {
	// If the following call pancis, it indicates UnimplementedPostsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Posts_ServiceDesc, srv)
}

func _Posts_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_CreatePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_CreateComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).CreateComment(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_ModifyPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).ModifyPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_ModifyPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).ModifyPost(ctx, req.(*ModifyPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_ModifyComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).ModifyComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_ModifyComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).ModifyComment(ctx, req.(*ModifyCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_DeletePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).DeletePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_DeletePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).DeletePost(ctx, req.(*DeletePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostsServer).DeleteComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Posts_DeleteComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostsServer).DeleteComment(ctx, req.(*DeleteCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Posts_GetPosts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetPostsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PostsServer).GetPosts(m, &grpc.GenericServerStream[GetPostsRequest, PostInfo]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Posts_GetPostsServer = grpc.ServerStreamingServer[PostInfo]

func _Posts_GetComments_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetCommentsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PostsServer).GetComments(m, &grpc.GenericServerStream[GetCommentsRequest, CommentInfo]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Posts_GetCommentsServer = grpc.ServerStreamingServer[CommentInfo]

// Posts_ServiceDesc is the grpc.ServiceDesc for Posts service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Posts_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "posts.Posts",
	HandlerType: (*PostsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePost",
			Handler:    _Posts_CreatePost_Handler,
		},
		{
			MethodName: "CreateComment",
			Handler:    _Posts_CreateComment_Handler,
		},
		{
			MethodName: "ModifyPost",
			Handler:    _Posts_ModifyPost_Handler,
		},
		{
			MethodName: "ModifyComment",
			Handler:    _Posts_ModifyComment_Handler,
		},
		{
			MethodName: "DeletePost",
			Handler:    _Posts_DeletePost_Handler,
		},
		{
			MethodName: "DeleteComment",
			Handler:    _Posts_DeleteComment_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPosts",
			Handler:       _Posts_GetPosts_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetComments",
			Handler:       _Posts_GetComments_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "posts.proto",
}
