package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
)

const (
	ClaimsKey contextKey = "claims"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				apierror.WriteError(w, apierror.ErrUnauthorized, nil)
				return
			}

			parts := strings.Split(header, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				apierror.WriteError(w, apierror.ErrUnauthorized, nil)
				return
			}
			claims, err := jwtManager.Verify(parts[1])

			if err != nil {
				apierror.WriteError(w, apierror.ErrUnauthorized, nil)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(ClaimsKey).(*auth.Claims)
	return claims
}
