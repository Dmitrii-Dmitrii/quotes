package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"quotes/internal/dtos"
)

type QuoteServiceInterface interface {
	CreateQuote(ctx context.Context, quoteDto dtos.QuoteDto) (*dtos.QuoteDto, error)
	DeleteQuote(ctx context.Context, id pgtype.UUID) error
	GetAllQuotes(ctx context.Context) ([]dtos.QuoteDto, error)
	GetQuotesByAuthor(ctx context.Context, author string) ([]dtos.QuoteDto, error)
	GetRandomQuote(ctx context.Context) (*dtos.QuoteDto, error)
}
