package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	Enabled          bool
	ServiceName      string
	Environment      string
	APMServerURL     string
	ElasticsearchURL string
}

type Tracer struct {
	config *Config
	tracer *sdktrace.TracerProvider
}

func NewTracer(config *Config) (*Tracer, error) {
	if !config.Enabled {
		return &Tracer{config: config}, nil
	}

	// Create OTLP HTTP exporter for APM Server
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(config.APMServerURL),
		otlptracehttp.WithInsecure(), // For development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1.0)), // 100% sampling for dev
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	return &Tracer{
		config: config,
		tracer: tp,
	}, nil
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.tracer != nil {
		return t.tracer.Shutdown(ctx)
	}
	return nil
}

func (t *Tracer) IsEnabled() bool {
	return t.config.Enabled
}
