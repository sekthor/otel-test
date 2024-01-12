package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/sekthor/otel-test/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	spanexp, _ := telemetry.SpanExporter(ctx)
	tp, _ := telemetry.TraceProvider("book-service", spanexp)

	propagator := telemetry.Propagator()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagator)

	tracer := tp.Tracer("booksvc")

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
	}

	Request(ctx, tracer, client)

	time.Sleep(100 * time.Second)
}

func Request(ctx context.Context, tracer trace.Tracer, client http.Client) {
	ctx, span := tracer.Start(ctx, "hello")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/authors/asdf", nil)
	if err != nil {
		fmt.Println("request prep failed: " + err.Error())
		return
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
}
