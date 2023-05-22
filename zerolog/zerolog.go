package zerolog

import (
	"context"

	"github.com/releaseband/go-micro-tools/trace"
	"github.com/rs/zerolog/log"
)

func WithTrace(ctx context.Context) context.Context {
	sc := trace.SpanContextFromContext(ctx)

	logger := log.With().
		Str("trace_id", sc.TraceID().String()).
		Str("span_id", sc.SpanID().String()).
		Logger()

	return logger.WithContext(ctx)
}
