package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"quotes/internal/drivers"
	"quotes/internal/dtos"
	"quotes/internal/models"
)

type QuoteService struct {
	driver *drivers.QuoteDriver
}

func NewQuoteService(driver *drivers.QuoteDriver) *QuoteService {
	return &QuoteService{driver: driver}
}

func (s *QuoteService) CreateQuote(ctx context.Context, quoteDto dtos.QuoteDto) (*dtos.QuoteDto, error) {
	id := generateUuid()

	quote := models.Quote{Id: id, Author: *quoteDto.Author, Text: *quoteDto.Text}
	err := s.driver.CreateQuote(ctx, quote)
	if err != nil {
		return nil, err
	}

	quoteDto.Id = &id

	return &quoteDto, nil
}

func (s *QuoteService) DeleteQuote(ctx context.Context, id pgtype.UUID) error {
	_, err := s.driver.GetQuoteById(ctx, id)
	if err != nil {
		return err
	}

	err = s.driver.DeleteQuote(ctx, id)
	return err
}

func (s *QuoteService) GetAllQuotes(ctx context.Context) ([]dtos.QuoteDto, error) {
	quotes, err := s.driver.GetAllQuotes(ctx)
	if err != nil {
		return nil, err
	}

	quoteDtos := make([]dtos.QuoteDto, len(quotes))
	for i, quote := range quotes {
		quoteDtos[i] = dtos.QuoteDto{Id: &quote.Id, Author: &quote.Author, Text: &quote.Text}
	}

	return quoteDtos, nil
}

func (s *QuoteService) GetQuoteByAuthor(ctx context.Context, author string) ([]dtos.QuoteDto, error) {
	quotes, err := s.driver.GetQuoteByAuthor(ctx, author)
	if err != nil {
		return nil, err
	}

	quoteDtos := make([]dtos.QuoteDto, len(quotes))
	for i, quote := range quotes {
		quoteDtos[i] = dtos.QuoteDto{Id: &quote.Id, Author: &quote.Author, Text: &quote.Text}
	}

	return quoteDtos, nil
}

func (s *QuoteService) GetRandomQuote(ctx context.Context) (*dtos.QuoteDto, error) {
	quote, err := s.driver.GetRandomQuote(ctx)
	if err != nil {
		return nil, err
	}

	quoteDto := &dtos.QuoteDto{Id: &quote.Id, Author: &quote.Author, Text: &quote.Text}
	return quoteDto, nil
}

func generateUuid() pgtype.UUID {
	newUuid := uuid.New()

	pgUuid := pgtype.UUID{
		Bytes: newUuid,
		Valid: true,
	}

	return pgUuid
}
