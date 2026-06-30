package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Link struct {
	Id             pgtype.UUID `json:"id"`
	ShortId        string      `validate:"required,len=8" json:"shortId"`
	Origin         string      `validate:"required,url,max=4096" json:"origin"`
	CreationDate   time.Time   `validate:"required" json:"creationDate"`
	LastAccessDate *time.Time  `json:"lastAccessDate"`
}
