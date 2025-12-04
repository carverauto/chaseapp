package observability

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"chaseapp.tv/api/internal/config"
)

// SetupTracing configures a global tracer provider. Returns a shutdown func.
func SetupTracing(ctx context.Context, cfg config.ObservabilityConfig, logger *slog.Logger) (func(context.Context) error, error) {
	if cfg.OTLPEndpoint == "" {
		logger.Warn("OTLP endpoint not configured; tracing disabled")
		return func(context.Context) error { return nil }, nil
	}

	clientOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
	}
	if cfg.OTLPInsecure {
		clientOpts = append(clientOpts, otlptracehttp.WithInsecure())
	}
	if cfg.OTLPHeaders != "" {
		headers := map[string]string{}
		for _, kv := range strings.Split(cfg.OTLPHeaders, ",") {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		clientOpts = append(clientOpts, otlptracehttp.WithHeaders(headers))
	}

	exporter, err := otlptracehttp.New(ctx, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("create otlp exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}
