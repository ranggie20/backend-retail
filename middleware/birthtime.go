package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/online-bnsp/backend/constant"
)

func BirthTime(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := time.Now()
		ctx = context.WithValue(ctx, constant.ContextBirthTime, t)

		requestID := r.Header.Get(middleware.RequestIDHeader)
		if requestID == "" {
			ctx = context.WithValue(ctx, middleware.RequestIDKey, uuid.New())
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
