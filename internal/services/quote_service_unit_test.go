package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"quotes/internal/dtos"
	"quotes/internal/models"
	"testing"
)

type MockQuoteDriver struct {
	mock.Mock
}

func (m *MockQuoteDriver) CreateQuote(ctx context.Context, quote *models.Quote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *MockQuoteDriver) DeleteQuote(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuoteDriver) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockQuoteDriver) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	args := m.Called(ctx, author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Quote), args.Error(1)
}

func (m *MockQuoteDriver) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *MockQuoteDriver) GetQuoteById(ctx context.Context, id pgtype.UUID) (*models.Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quote), args.Error(1)
}

func TestCreateQuote(t *testing.T) {
	ctx := context.Background()
	mockDriver := new(MockQuoteDriver)
	quoteService := NewQuoteService(mockDriver)

	author := "author"
	text := "text"
	sendQuoteDto := dtos.QuoteDto{
		Author: &author,
		Text:   &text,
	}

	mockDriver.On("CreateQuote", mock.Anything, mock.Anything).Return(nil)

	getQuoteDto, err := quoteService.CreateQuote(ctx, sendQuoteDto)

	assert.NoError(t, err)
	assert.Equal(t, author, *getQuoteDto.Author)
	assert.Equal(t, text, *getQuoteDto.Text)
	mockDriver.AssertExpectations(t)
}

func TestDeleteQuote(t *testing.T) {
	ctx := context.Background()
	mockDriver := new(MockQuoteDriver)
	quoteService := NewQuoteService(mockDriver)

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}
	author := "author"
	text := "text"

	mockDriver.On("GetQuoteById", mock.Anything, mock.Anything).Return(&models.Quote{
		Id:     id,
		Author: author,
		Text:   text,
	}, nil)
	mockDriver.On("DeleteQuote", mock.Anything, mock.Anything).Return(nil)

	err := quoteService.DeleteQuote(ctx, id)

	assert.NoError(t, err)
	mockDriver.AssertExpectations(t)
}

func TestGetAllQuotes(t *testing.T) {
	ctx := context.Background()
	mockDriver := new(MockQuoteDriver)
	quoteService := NewQuoteService(mockDriver)

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}
	author := "author"
	text := "text"

	mockDriver.On("GetAllQuotes", mock.Anything).Return([]models.Quote{
		{
			Id:     id,
			Author: author,
			Text:   text,
		},
	}, nil)

	quoteDtos, err := quoteService.GetAllQuotes(ctx)
	assert.NoError(t, err)
	assert.Len(t, quoteDtos, 1)
	assert.Equal(t, id, *quoteDtos[0].Id)
	assert.Equal(t, author, *quoteDtos[0].Author)
	assert.Equal(t, text, *quoteDtos[0].Text)
	mockDriver.AssertExpectations(t)
}

func TestGetQuotesByAuthor(t *testing.T) {
	ctx := context.Background()
	mockDriver := new(MockQuoteDriver)
	quoteService := NewQuoteService(mockDriver)

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}
	author := "author"
	text := "text"

	mockDriver.On("GetQuotesByAuthor", mock.Anything, mock.Anything).Return([]models.Quote{
		{
			Id:     id,
			Author: author,
			Text:   text,
		},
	}, nil)

	quoteDtos, err := quoteService.GetQuotesByAuthor(ctx, author)
	assert.NoError(t, err)
	assert.Len(t, quoteDtos, 1)
	assert.Equal(t, id, *quoteDtos[0].Id)
	assert.Equal(t, author, *quoteDtos[0].Author)
	assert.Equal(t, text, *quoteDtos[0].Text)
	mockDriver.AssertExpectations(t)
}

func TestGetRandomQuote(t *testing.T) {
	ctx := context.Background()
	mockDriver := new(MockQuoteDriver)
	quoteService := NewQuoteService(mockDriver)

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}
	author := "author"
	text := "text"

	mockDriver.On("GetRandomQuote", mock.Anything, mock.Anything).Return(&models.Quote{
		Id:     id,
		Author: author,
		Text:   text,
	}, nil)

	quoteDto, err := quoteService.GetRandomQuote(ctx)
	assert.NoError(t, err)
	assert.Equal(t, id, *quoteDto.Id)
	assert.Equal(t, author, *quoteDto.Author)
	assert.Equal(t, text, *quoteDto.Text)
	mockDriver.AssertExpectations(t)
}
