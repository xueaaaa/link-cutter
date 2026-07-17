package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserModel struct {
	Id             pgtype.UUID
	Email          string
	Username       string
	Password       string
	CreationDate   time.Time
	LastAccessDate *time.Time
}
