package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"quotes/internal/dtos"
	"quotes/internal/services"
)

type QuoteController struct {
	service *services.QuoteService
}

func NewQuoteController(service *services.QuoteService) *QuoteController {
	return &QuoteController{service: service}
}

func (c *QuoteController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/quotes", c.handleQuotes).Methods("GET", "POST")
	router.HandleFunc("/quotes/random", c.getRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id}", c.deleteQuote).Methods("DELETE")
}

func (c *QuoteController) handleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getQuotes(w, r)
	case "POST":
		c.createQuote(w, r)
	}
}

func (c *QuoteController) createQuote(w http.ResponseWriter, r *http.Request) {
	var quoteDto dtos.QuoteDto

	if err := json.NewDecoder(r.Body).Decode(&quoteDto); err != nil {
		c.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if quoteDto.Author == nil || strings.TrimSpace(*quoteDto.Author) == "" {
		c.writeErrorResponse(w, "Author is required", http.StatusBadRequest)
		return
	}

	if quoteDto.Text == nil || strings.TrimSpace(*quoteDto.Text) == "" {
		c.writeErrorResponse(w, "Text is required", http.StatusBadRequest)
		return
	}

	createdQuote, err := c.service.CreateQuote(r.Context(), quoteDto)
	if err != nil {
		c.writeErrorResponse(w, "Failed to create quote", http.StatusInternalServerError)
		return
	}

	c.writeJSONResponse(w, createdQuote, http.StatusCreated)
}

func (c *QuoteController) getQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []dtos.QuoteDto
	var err error

	if author != "" {
		quotes, err = c.service.GetQuoteByAuthor(r.Context(), author)
	} else {
		quotes, err = c.service.GetAllQuotes(r.Context())
	}

	if err != nil {
		c.writeErrorResponse(w, "Failed to retrieve quotes", http.StatusInternalServerError)
		return
	}

	c.writeJSONResponse(w, quotes, http.StatusOK)
}

func (c *QuoteController) getRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := c.service.GetRandomQuote(r.Context())
	if err != nil {
		c.writeErrorResponse(w, "Failed to retrieve random quote", http.StatusInternalServerError)
		return
	}

	c.writeJSONResponse(w, quote, http.StatusOK)
}

func (c *QuoteController) deleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		c.writeErrorResponse(w, "Quote ID is required", http.StatusBadRequest)
		return
	}

	pgUuid, err := c.parseUUID(idStr)
	if err != nil {
		c.writeErrorResponse(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteQuote(r.Context(), pgUuid)
	if err != nil {
		c.writeErrorResponse(w, "Failed to delete quote", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *QuoteController) parseUUID(uuidStr string) (pgtype.UUID, error) {
	uuidStr = strings.TrimSpace(uuidStr)

	if len(uuidStr) != 36 ||
		uuidStr[8] != '-' || uuidStr[13] != '-' ||
		uuidStr[18] != '-' || uuidStr[23] != '-' {
		return pgtype.UUID{}, &InvalidUUIDError{UUID: uuidStr}
	}

	var pgUuid pgtype.UUID
	err := pgUuid.Scan(uuidStr)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgUuid, nil
}

func (c *QuoteController) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (c *QuoteController) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errorResponse)
}

type InvalidUUIDError struct {
	UUID string
}

func (e *InvalidUUIDError) Error() string {
	return "invalid UUID format: " + e.UUID
}
