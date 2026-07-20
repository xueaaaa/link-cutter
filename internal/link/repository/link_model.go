package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type LinkModel struct {
	Id             pgtype.UUID
	UserId         pgtype.UUID
	ShortId        string
	Origin         string
	CreationDate   time.Time
	LastAccessDate *time.Time
}
