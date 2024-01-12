package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/otel-test/model"
	"go.opentelemetry.io/otel/trace"
)

type BookService struct {
	Tracer trace.Tracer
	Client http.Client
}

func (s *BookService) GetBookByID(c *gin.Context) {
	//var err error
	ctx, span := s.Tracer.Start(context.Background(), "GetAuthorByID")
	defer span.End()

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformatted or missing id"})
		return
	}

	book := s.FetchBookByID(ctx, id)
	author, err := s.GetAuthor(ctx, book.Author)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book":   &book,
		"author": &author,
	})
}

func (s *BookService) FetchBookByID(ctx context.Context, id string) model.Book {
	_, span := s.Tracer.Start(ctx, "FetchAuthorByID")
	defer span.End()

	return model.Book{
		ID:     id,
		Title:  "Title of the Book",
		Author: "a4369f5e-074e-458a-b6a7-dcc650e79066",
	}
}

func (s *BookService) GetAuthor(ctx context.Context, authorId string) (model.Author, error) {
	var author model.Author

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/authors/"+authorId, nil)
	if err != nil {
		return author, err
	}
	_, err = s.Client.Do(req)
	if err != nil {
		return author, err
	}
	//defer resp.Body.Close()
	return author, nil
}
