package util

import (
	"context"
	"net/http"

	middleware2 "github.com/go-chi/chi/v5/middleware"
)

func GetRequestId(r *http.Request) string {
	return r.Context().Value(middleware2.RequestIDKey).(string)
}

func PutRequestId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, middleware2.RequestIDKey, reqId)
}

func GetRequestIdFromContext(ctx context.Context) string {
	return ctx.Value(middleware2.RequestIDKey).(string)
}
