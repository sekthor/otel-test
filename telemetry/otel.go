package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func SpanExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	endpoint := os.Getenv("OTELENDPOINT")

	// no endpoint -> traces to stdout
	if endpoint == "" {
		return stdouttrace.New()
	}

	return otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	)
}

func TraceProvider(serviceName string, exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {

	// use an empty resource
	// res, err := resource.New(context.TODO())

	// create a resource with otel defaults + the schemaurl and our service name
	res, err := resource.Merge(
		// default: opentelemetry metadata -> sdk infos
		resource.Default(),

		// addtional attributes
		resource.NewWithAttributes(

			// schemaURL: url to the opentelemetry schema used
			semconv.SchemaURL,

			// the global name of our service
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	), nil
}

func Propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
