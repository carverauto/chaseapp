package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var propagator = propagation.TraceContext{}

// Tracing starts a span for each incoming request and injects trace info into the context.
func Tracing(next http.Handler) http.Handler {
	tracer := otel.Tracer("chaseapp/api")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path, trace.WithSpanKind(trace.SpanKindServer))
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.target", r.URL.Path),
		)
		defer span.End()

		// Propagate traceparent header downstream
		propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
