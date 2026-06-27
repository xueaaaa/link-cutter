package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Link struct {
	Id             pgtype.UUID
	ShortId        string    `validate:"required,len=8"`
	Origin         string    `validate:"required,url,max=4096"`
	CreationDate   time.Time `validate:"required"`
	LastAccessDate *time.Time
}
