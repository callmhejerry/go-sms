package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestId)

		fmt.Printf("Setting request id  to be %v\n", requestId)
		ctx := context.WithValue(r.Context(), RequestIDKey, requestId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HELPER TO GET REQUEST ID FROM THE CONTEXT
func GetRequestId(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
