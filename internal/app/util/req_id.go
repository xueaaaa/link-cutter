package util

import (
	"net/http"

	middleware2 "github.com/go-chi/chi/v5/middleware"
)

func GetRequestId(r *http.Request) string {
	return r.Context().Value(middleware2.RequestIDKey).(string)
}
