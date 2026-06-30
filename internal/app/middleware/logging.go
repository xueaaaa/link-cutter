package middleware

import (
	"link-cutter/internal/app/util"
	"net/http"

	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("req_id", util.GetRequestId(r)),
			)

			next.ServeHTTP(w, r)
		})
	}
}
