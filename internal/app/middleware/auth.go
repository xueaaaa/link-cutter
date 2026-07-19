package middleware

import (
	"context"
	errors2 "link-cutter/internal/app/errors"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type contextKey string

const userContextKey contextKey = "user"

func Auth(key string, logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Info(errors2.ErrMissingAuthHeader.Error(),
					zap.String("req_id", util.GetRequestId(r)),
				)
				util.WriteError(
					w,
					http.StatusUnauthorized,
					errors2.ErrMissingAuthHeader.Error(),
					util.GetRequestId(r),
				)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Info(errors2.ErrInvalidAuthHeader.Error(),
					zap.String("req_id", util.GetRequestId(r)),
				)
				util.WriteError(
					w,
					http.StatusUnauthorized,
					errors2.ErrInvalidAuthHeader.Error(),
					util.GetRequestId(r),
				)
			}

			claims := model.Claims{}
			token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors2.ErrUnexpectedSigningMethod
				}
				return key, nil
			})
			if err != nil || !token.Valid {
				logger.Info(err.Error(),
					zap.String("req_id", util.GetRequestId(r)),
				)
				util.WriteError(
					w,
					http.StatusUnauthorized,
					err.Error(),
					util.GetRequestId(r),
				)
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
