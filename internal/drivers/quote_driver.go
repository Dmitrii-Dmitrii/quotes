package drivers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"quotes/internal/models"
)

type QuoteDriver struct {
	adapter Adapter
}

func NewQuoteDriver(adapter Adapter) *QuoteDriver {
	return &QuoteDriver{adapter: adapter}
}

func (d *QuoteDriver) CreateQuote(ctx context.Context, quote models.Quote) error {
	_, err := d.adapter.Exec(
		ctx,
		queryCreateQuote,
		quote.Id,
		quote.Author,
		quote.Text,
	)

	return err
}

func (d *QuoteDriver) DeleteQuote(ctx context.Context, id pgtype.UUID) error {
	_, err := d.adapter.Exec(ctx, queryDeleteQuote, id)

	return err
}

func (d *QuoteDriver) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	rows, err := d.adapter.Query(ctx, queryGetAllQuotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		quote := models.Quote{}

		err = rows.Scan(&quote.Id, &quote.Author, &quote.Text)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (d *QuoteDriver) GetQuoteByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	rows, err := d.adapter.Query(ctx, queryGetQuoteByAuthor, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		quote := models.Quote{Author: author}

		err = rows.Scan(&quote.Id, &quote.Text)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (d *QuoteDriver) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	quote := models.Quote{}

	err := d.adapter.QueryRow(ctx, queryGetRandomQuote).Scan(&quote.Id, &quote.Author, &quote.Text)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (d *QuoteDriver) GetQuoteById(ctx context.Context, id pgtype.UUID) (*models.Quote, error) {
	quote := models.Quote{Id: id}

	err := d.adapter.QueryRow(ctx, queryGetQuoteById, id).Scan(&quote.Author, &quote.Text)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}
