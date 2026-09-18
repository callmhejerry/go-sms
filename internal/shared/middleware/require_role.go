package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/callmhejerry/sms/internal/identity"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/google/uuid"
)

type RoleChecker interface {
	UserHasRole(ctx context.Context, userId uuid.UUID, roleName identity.RoleName) (bool, error)
}

func RequireRole(checker RoleChecker, roleName identity.RoleName) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())

			fmt.Printf("Gotten claims %v\n", claims)
			if claims == nil {
				apierror.WriteError(w, apierror.ErrUnauthorized, nil)
				return
			}
			hasRole, err := checker.UserHasRole(r.Context(), claims.UserID, roleName)

			if err != nil {
				apierror.WriteError(w, apierror.ErrInternal, nil)
				return
			}

			if !hasRole {
				apierror.WriteError(w, apierror.ErrForbidden, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
