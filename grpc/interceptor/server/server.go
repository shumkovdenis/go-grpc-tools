package server

import (
	"context"

	"github.com/releaseband/go-grpc-tools/trace"
	"github.com/releaseband/go-grpc-tools/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Trace() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			carrier := MetadataCarrier(md)
			ctx = trace.ExtractBinaryTrace(ctx, carrier)
		}
		return handler(ctx, req)
	}
}

func TraceLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = zerolog.WithTrace(ctx)
		return handler(ctx, req)
	}
}
