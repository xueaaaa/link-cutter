package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Link struct {
	Id             pgtype.UUID
	ShortId        string
	Origin         string
	CreationDate   time.Time
	LastAccessDate *time.Time
}
