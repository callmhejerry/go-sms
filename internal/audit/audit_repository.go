package audit

import (
	"context"
	"encoding/json"

	"github.com/callmhejerry/sms/internal/shared/store"
)

type AuditRepository interface {
	Log(ctx context.Context, e Entry)
}

type auditRepositoryImpl struct {
	queries *store.Queries
}

func NewAuditRepositoryImpl(queries *store.Queries) AuditRepository {
	return &auditRepositoryImpl{
		queries: queries,
	}
}

func (repo *auditRepositoryImpl) Log(ctx context.Context, e Entry) {
	var meta []byte
	if e.Metadata != nil {
		meta, _ = json.Marshal(e.Metadata)
	}

	// Fire-and-forget style (we don't want audit failures to break the main flow)
	_, _ = repo.queries.CreateAuditLog(ctx, store.CreateAuditLogParams{
		TenantID:     e.TenantID,
		UserID:       e.UserID,
		Action:       e.Action,
		ResourceType: &e.ResourceType,
		ResourceID:   e.ResourceID,
		Metadata:     meta,
		IpAddress:    &e.IPAddress,
		UserAgent:    &e.UserAgent,
		RequestID:    &e.RequestID,
	})
}
