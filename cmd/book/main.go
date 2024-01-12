package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/otel-test/service"
	"github.com/sekthor/otel-test/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func main() {

	// otel setup
	ctx := context.Background()

	// create the exporter for traces
	// it specifies where the indiviual span data is sent to (here stdout)
	spanexp, err := telemetry.SpanExporter(ctx)
	if err != nil {
		log.Fatalf("could not create span exporter")
	}

	// create the traceprovider, with a given exporter for spans
	tp, err := telemetry.TraceProvider("book-service", spanexp)
	if err != nil {
		log.Fatal("could not create tracer provider: " + err.Error())
	}

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

	svc := service.BookService{
		Tracer: tracer,
		Client: client,
	}

	router := gin.New()
	router.GET("books/:id", svc.GetBookByID)

	err = router.Run("0.0.0.0:8081")
	if err != nil {
		log.Fatal("failed to run server")
	}
}
