package trace

import (
	"context"

	"github.com/shumkovdenis/go-micro-tools/trace/utils"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceparentHeader  = "Traceparent"
	TracestateHeader   = "Tracestate"
	GrpcTraceBinHeader = "Grpc-Trace-Bin"
)

func SpanContextFromContext(ctx context.Context) trace.SpanContext {
	return trace.SpanContextFromContext(ctx)
}

func TraceFromContext(ctx context.Context) (string, string) {
	sc := SpanContextFromContext(ctx)
	traceparent := utils.SpanContextToW3CString(sc)
	tracestate := utils.TraceStateToW3CString(sc)
	return traceparent, tracestate
}

func BinaryTraceFromContext(ctx context.Context) string {
	sc := SpanContextFromContext(ctx)
	grpcTraceBin := utils.BinaryFromSpanContext(sc)
	return string(grpcTraceBin)
}

func ExtractTrace(ctx context.Context, carrier Carrier) context.Context {
	traceparent := carrier.Get(TraceparentHeader)
	sc, _ := utils.SpanContextFromW3CString(traceparent)
	if !sc.IsValid() {
		return ctx
	}

	tracestate := carrier.Get(TracestateHeader)
	ts := utils.TraceStateFromW3CString(tracestate)
	sc = sc.WithTraceState(ts)

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func ExtractBinaryTrace(ctx context.Context, carrier Carrier) context.Context {
	grpcTraceBin := carrier.Get(GrpcTraceBinHeader)
	sc, _ := utils.SpanContextFromBinary([]byte(grpcTraceBin))
	if !sc.IsValid() {
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func InjectTrace(ctx context.Context, carrier Carrier) {
	traceparent, tracestate := TraceFromContext(ctx)
	carrier.Set(TraceparentHeader, traceparent)
	carrier.Set(TracestateHeader, tracestate)
}

func InjectBinaryTrace(ctx context.Context, carrier Carrier) {
	grpcTraceBin := BinaryTraceFromContext(ctx)
	carrier.Set(GrpcTraceBinHeader, grpcTraceBin)
}
