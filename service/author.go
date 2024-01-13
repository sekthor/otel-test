package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/otel-test/model"
	"go.opentelemetry.io/otel/trace"
)

type AuthorService struct {
	Tracer trace.Tracer
}

func (s *AuthorService) GetAuthorByID(c *gin.Context) {
	ctx, span := s.Tracer.Start(c.Request.Context(), "GetAuthorByID")
	defer span.End()

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformatted or missing id"})
		return
	}

	// call repository function
	author := s.FetchAuthorByID(ctx, id)

	c.JSON(http.StatusOK, &author)
}

func (s *AuthorService) FetchAuthorByID(ctx context.Context, id string) model.Author {
	_, span := s.Tracer.Start(ctx, "FetchAuthorByID")
	defer span.End()

	return model.Author{
		ID:        id,
		Firstname: "Max",
		Lastname:  "Muster",
	}
}
