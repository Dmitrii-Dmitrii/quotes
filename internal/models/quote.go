package models

import "github.com/jackc/pgx/v5/pgtype"

type Quote struct {
	Id     pgtype.UUID
	Author string
	Text   string
}
