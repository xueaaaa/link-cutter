package middleware

import (
	"context"
	"errors"
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

func parseToken(key string, r *http.Request) (*jwt.Token, model.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, model.Claims{}, errors2.ErrMissingAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, model.Claims{}, errors2.ErrInvalidAuthHeader
	}

	claims := model.Claims{}
	token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors2.ErrUnexpectedSigningMethod
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, model.Claims{}, err
	}
	return token, claims, nil
}

func Auth(key string, logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := parseToken(key, r)

			if err != nil || token == nil || !token.Valid {
				msg := "invalid token"
				if err != nil {
					msg = err.Error()
				}
				logger.Info(msg,
					zap.String("req_id", util.GetRequestId(r)),
				)
				util.WriteError(
					w,
					http.StatusUnauthorized,
					msg,
					util.GetRequestId(r),
				)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuth(key string, logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := parseToken(key, r)

			switch {
			case err == nil:
				ctx := context.WithValue(r.Context(), userContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			case errors.Is(err, errors2.ErrMissingAuthHeader):
				next.ServeHTTP(w, r)
				return
			default:
				logger.Info(err.Error(),
					zap.String("req_id", util.GetRequestId(r)),
				)
				util.WriteError(
					w,
					http.StatusUnauthorized,
					err.Error(),
					util.GetRequestId(r),
				)
				return
			}
		})
	}
}

func ClaimsFromContext(ctx context.Context) (model.Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(model.Claims)
	return claims, ok
}
