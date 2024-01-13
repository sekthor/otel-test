package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/otel-test/service"
	"github.com/sekthor/otel-test/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	tp, err := telemetry.TraceProvider("author-service", spanexp)
	if err != nil {
		log.Fatal("could not create tracer provider: " + err.Error())
	}

	propagator := telemetry.Propagator()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagator)

	tracer := tp.Tracer("authorsvc")

	svc := service.AuthorService{
		Tracer: tracer,
	}

	router := gin.New()
	router.Use(otelgin.Middleware("bookservice-otelgin"))
	router.GET("authors/:id", svc.GetAuthorByID)

	err = router.Run()
	if err != nil {
		log.Fatal("failed to run server")
	}
}
