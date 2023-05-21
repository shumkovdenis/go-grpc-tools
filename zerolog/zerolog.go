package zerolog

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shumkovdenis/go-micro-tools/trace"
)

func WithTrace(ctx context.Context) context.Context {
	sc := trace.SpanContextFromContext(ctx)

	logger := log.With().
		Str("trace_id", sc.TraceID().String()).
		Str("span_id", sc.SpanID().String()).
		Logger()

	return logger.WithContext(ctx)
}
