package drivers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"quotes/internal/models"
)

type QuoteDriverInterface interface {
	CreateQuote(ctx context.Context, quote *models.Quote) error
	DeleteQuote(ctx context.Context, id pgtype.UUID) error
	GetAllQuotes(ctx context.Context) ([]models.Quote, error)
	GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error)
	GetRandomQuote(ctx context.Context) (*models.Quote, error)
	GetQuoteById(ctx context.Context, id pgtype.UUID) (*models.Quote, error)
}
