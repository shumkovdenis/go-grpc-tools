package client

import (
	"context"

	"github.com/releaseband/go-grpc-tools/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithMetadata(key, value string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, key, value)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func WithDaprApp(appID string) grpc.UnaryClientInterceptor {
	return WithMetadata("dapr-app-id", appID)
}

func WithTrace() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		traceparent, tracestate := trace.TraceFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx,
			trace.TraceparentHeader, traceparent,
			trace.TracestateHeader, tracestate)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func WithBinaryTrace() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcTraceBin := trace.BinaryTraceFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx,
			trace.GrpcTraceBinHeader, grpcTraceBin)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
