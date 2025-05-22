package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"quotes/internal/dtos"
)

type MockQuoteService struct {
	mock.Mock
}

func (m *MockQuoteService) CreateQuote(ctx context.Context, dto dtos.QuoteDto) (*dtos.QuoteDto, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.QuoteDto), args.Error(1)
}

func (m *MockQuoteService) GetAllQuotes(ctx context.Context) ([]dtos.QuoteDto, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dtos.QuoteDto), args.Error(1)
}

func (m *MockQuoteService) GetQuotesByAuthor(ctx context.Context, author string) ([]dtos.QuoteDto, error) {
	args := m.Called(ctx, author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dtos.QuoteDto), args.Error(1)
}

func (m *MockQuoteService) GetRandomQuote(ctx context.Context) (*dtos.QuoteDto, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.QuoteDto), args.Error(1)
}

func (m *MockQuoteService) DeleteQuote(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateQuote(t *testing.T) {
	mockService := &MockQuoteService{}
	controller := NewQuoteController(mockService)

	author := "author"
	text := "text"
	inputDto := dtos.QuoteDto{
		Author: &author,
		Text:   &text,
	}
	expectedDto := dtos.QuoteDto{
		Author: &author,
		Text:   &text,
	}

	mockService.On("CreateQuote", mock.Anything, inputDto).Return(&expectedDto, nil)

	jsonBody, _ := json.Marshal(inputDto)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	controller.createQuote(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var responseDto dtos.QuoteDto
	err := json.Unmarshal(rr.Body.Bytes(), &responseDto)
	assert.NoError(t, err)
	assert.Equal(t, expectedDto, responseDto)
	mockService.AssertExpectations(t)
}

func TestGetAllQuotes(t *testing.T) {
	mockService := &MockQuoteService{}
	controller := NewQuoteController(mockService)

	author := "author"
	text := "text"
	expectedQuotes := []dtos.QuoteDto{
		{
			Author: &author,
			Text:   &text,
		},
	}
	mockService.On("GetAllQuotes", mock.Anything).Return(expectedQuotes, nil)

	req := httptest.NewRequest("GET", "/quotes", nil)
	rr := httptest.NewRecorder()

	controller.getQuotes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var responseQuotes []dtos.QuoteDto
	err := json.Unmarshal(rr.Body.Bytes(), &responseQuotes)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuotes, responseQuotes)

	mockService.AssertExpectations(t)
}

func TestGetQuotesByAuthor(t *testing.T) {
	mockService := &MockQuoteService{}
	controller := NewQuoteController(mockService)

	author := "author"
	text := "text"
	expectedQuotes := []dtos.QuoteDto{
		{
			Author: &author,
			Text:   &text,
		},
	}
	mockService.On("GetQuotesByAuthor", mock.Anything, author).Return(expectedQuotes, nil)

	req := httptest.NewRequest("GET", "/quotes?author="+author, nil)
	rr := httptest.NewRecorder()

	controller.getQuotes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseQuotes []dtos.QuoteDto
	err := json.Unmarshal(rr.Body.Bytes(), &responseQuotes)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuotes, responseQuotes)

	mockService.AssertExpectations(t)
}

func TestGetRandomQuote(t *testing.T) {
	mockService := &MockQuoteService{}
	controller := NewQuoteController(mockService)

	author := "author"
	text := "text"
	expectedQuote := dtos.QuoteDto{
		Author: &author,
		Text:   &text,
	}
	mockService.On("GetRandomQuote", mock.Anything).Return(&expectedQuote, nil)

	req := httptest.NewRequest("GET", "/quotes/random", nil)
	rr := httptest.NewRecorder()

	controller.getRandomQuote(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var responseQuote dtos.QuoteDto
	err := json.Unmarshal(rr.Body.Bytes(), &responseQuote)
	assert.NoError(t, err)
	assert.Equal(t, expectedQuote, responseQuote)

	mockService.AssertExpectations(t)
}

func TestDeleteQuote(t *testing.T) {
	mockService := &MockQuoteService{}
	controller := NewQuoteController(mockService)

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}

	mockService.On("DeleteQuote", mock.Anything, id).Return(nil)

	req := httptest.NewRequest("DELETE", "/quotes/"+idBytes.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": idBytes.String()})
	rr := httptest.NewRecorder()

	controller.deleteQuote(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())

	mockService.AssertExpectations(t)
}
