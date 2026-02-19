package indexers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// OTELProviders holds the OTEL provider instances for shutdown.
type OTELProviders struct {
	log    *sdklog.LoggerProvider
	metric *sdkmetric.MeterProvider
	trace  *sdktrace.TracerProvider
}

// InitOTEL initializes OpenTelemetry providers for the indexer service.
// Uses OTEL_EXPORTER_OTLP_ENDPOINT from environment, defaults to localhost:4318.
func InitOTEL(ctx context.Context, serviceName string) (*OTELProviders, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	p := &OTELProviders{}

	// Log provider
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(endpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating log exporter: %w", err)
	}
	p.log = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)

	// Metric provider
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		p.Shutdown(ctx)
		return nil, fmt.Errorf("creating metric exporter: %w", err)
	}
	p.metric = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(p.metric)

	// Trace provider
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		p.Shutdown(ctx)
		return nil, fmt.Errorf("creating trace exporter: %w", err)
	}
	p.trace = sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter)),
	)
	otel.SetTracerProvider(p.trace)

	return p, nil
}

// Shutdown gracefully shuts down all providers.
func (p *OTELProviders) Shutdown(ctx context.Context) error {
	var errs []error
	if p.log != nil {
		if err := p.log.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if p.metric != nil {
		if err := p.metric.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if p.trace != nil {
		if err := p.trace.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
