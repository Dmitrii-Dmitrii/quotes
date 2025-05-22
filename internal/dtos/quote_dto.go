package dtos

import "github.com/jackc/pgx/v5/pgtype"

type QuoteDto struct {
	Id     *pgtype.UUID `json:"id"`
	Author *string      `json:"author"`
	Text   *string      `json:"text"`
}
