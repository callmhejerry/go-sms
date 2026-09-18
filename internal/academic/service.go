package academic

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *store.Queries
}

func NewService(queries *store.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (service *Service) CreateAcademicSession(ctx context.Context, tenantId uuid.UUID, request CreateAcademicSessionRequest) (*store.AcademicSession, error) {

	name := strings.TrimSpace(request.Name)

	if name == "" {
		return nil, apierror.Validation("name is required")
	}
	if request.StartDate.IsZero() || request.EndDate.IsZero() {
		return nil, apierror.Validation("start_date and end_date is required")
	}

	if request.EndDate.Before(request.StartDate) {
		return nil, apierror.Validation("end_date cannot be before start_date")
	}
	if request.IsCurrent {
		_ = service.queries.SetCurrentAcademicSession(ctx, pgtype.UUID{
			Bytes: tenantId,
			Valid: true,
		})
	}

	session, err := service.queries.CreateAcademicSession(ctx, store.CreateAcademicSessionParams{
		TenantID: pgtype.UUID{
			Bytes: tenantId,
			Valid: true,
		},
		Name: name,
		StartDate: pgtype.Date{
			Time:  request.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  request.EndDate,
			Valid: true,
		},
		IsCurrent: request.IsCurrent,
	})

	if err != nil {
		return nil, database.TranslateError(err)
	}
	return &session, nil
}

func (service *Service) ListAcademicSession(ctx context.Context, tenantId uuid.UUID) ([]store.AcademicSession, error) {
	sessions, err := service.queries.ListAcademicSessions(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})

	if err != nil {
		return nil, database.TranslateError(err)
	}
	return sessions, nil
}

func (service *Service) GetCurrentAcademicSession(ctx context.Context, tenantId uuid.UUID) (*store.AcademicSession, error) {
	currentSession, err := service.queries.GetCurrentAcademicSession(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})
	if err != nil {
		return nil, database.TranslateError(err)
	}
	return &currentSession, nil
}
