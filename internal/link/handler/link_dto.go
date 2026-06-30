package handler

import "github.com/jackc/pgx/v5/pgtype"

type CreateLinkDTO struct {
	Origin string `json:"origin"`
}

type EditLinkDTO struct {
	Id     pgtype.UUID `json:"id"`
	Origin string      `json:"origin"`
}
