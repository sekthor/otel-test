package telemetry

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Formatter(c *gin.Context) otelgin.SpanNameFormatter {
	return func(*http.Request) string {
		return fmt.Sprintf("%s%s", "Handle", c.HandlerName())
	}
}

func Middleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		option := otelgin.WithSpanNameFormatter(Formatter(c))
		otelgin.Middleware(service, option)
	}
}
