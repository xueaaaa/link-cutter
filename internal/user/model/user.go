package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id             pgtype.UUID `json:"id"`
	Email          string      `validate:"required,email" json:"email"`
	Username       string      `validate:"required,min=4,max=16" json:"username"`
	Password       string      `json:"-"`
	CreationDate   time.Time   `validate:"required" json:"creationDate"`
	LastAccessDate *time.Time  `json:"lastAccessDate"`
}
