package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Claims struct {
	UserId   pgtype.UUID `json:"userId"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
	jwt.RegisteredClaims
}
