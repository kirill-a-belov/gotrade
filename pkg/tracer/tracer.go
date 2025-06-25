package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)
	span.SetName(name)
	// TODO Connect tracing

	return trace.ContextWithSpan(ctx, span), span
}
