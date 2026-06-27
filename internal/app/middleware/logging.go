package middleware

import (
	"net/http"

	middleware2 "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("req_id", r.Context().Value(middleware2.RequestIDKey).(string)),
			)

			next.ServeHTTP(w, r)
		})
	}
}
