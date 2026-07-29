package handler

import "github.com/jackc/pgx/v5/pgtype"

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditUserDTO struct {
	Id       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	Password string      `json:"password"`
}

type AuthUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
