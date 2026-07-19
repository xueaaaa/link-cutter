package util

import (
	"link-cutter/internal/user/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueToken(key string, user model.User) (string, error) {
	claims := model.Claims{
		UserId:   user.Id,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(key)
}
