package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

const TenantIdKey contextKey = "tenant-id-key"

type RequireTenantMiddleware struct {
	logger  *slog.Logger
	queries *store.Queries
}

var ErrRequiresTenantId = apierror.New(
	"requires_tenant_id",
	"Unauthorized. Requires tenant ID",
	http.StatusForbidden, nil,
	nil)

func NewRequireTenantMiddleWare(
	logger *slog.Logger,
	queries *store.Queries,
) *RequireTenantMiddleware {
	return &RequireTenantMiddleware{
		logger:  logger,
		queries: queries,
	}
}

func (m *RequireTenantMiddleware) RequireTenantId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r.Context())

		if claims == nil {
			apierror.WriteError(w, apierror.ErrUnauthorized, m.logger)
		}

		tenantId, err := uuid.Parse(r.Header.Get("X-Tenant-ID"))

		if err != nil {
			apierror.WriteError(w, ErrRequiresTenantId, m.logger)
			return
		}

		hasUser, err := m.queries.HasUser(r.Context(), store.HasUserParams{
			TenantID: tenantId,
			UserID:   claims.UserID,
		})

		if err != nil || !hasUser {
			apierror.WriteError(w, ErrRequiresTenantId, m.logger)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), TenantIdKey, tenantId)))
	})
}

func GetTenantId(ctx context.Context) (uuid.UUID, bool, error) {
	tenantId, ok := ctx.Value(TenantIdKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, false, ErrRequiresTenantId
	}
	return tenantId, true, nil
}
