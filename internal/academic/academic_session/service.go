package academicsession

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/callmhejerry/sms/internal/academic"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries            *store.Queries
	pool               *pgxpool.Pool
	AcademicRepository *AcademicRepository
}

func NewService(queries *store.Queries, pool *pgxpool.Pool) *Service {
	return &Service{
		queries: queries,
		pool:    pool,
	}
}

func (service *Service) CreateAcademicSession(ctx context.Context, tenantId uuid.UUID, request CreateAcademicSessionRequest) (*store.AcademicSession, error) {
	name := strings.TrimSpace(request.Name)
	startDate, _ := time.Parse("2026-01-29", request.StartDate)
	endDate, _ := time.Parse("2026-01-30", request.EndDate)

	if startDate.IsZero() || endDate.IsZero() {
		return nil, apierror.Validation("start_date and end_date is required")
	}

	if endDate.Before(startDate) {
		return nil, apierror.Validation("end_date cannot be before start_date")
	}
	if request.IsCurrent {
		_ = service.queries.SetCurrentAcademicSession(ctx, tenantId)
	}

	session, err := service.queries.CreateAcademicSession(ctx, store.CreateAcademicSessionParams{
		TenantID:  tenantId,
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		IsCurrent: request.IsCurrent,
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return &session, nil
}

func (service *Service) ListAcademicSession(ctx context.Context, tenantId uuid.UUID) ([]store.AcademicSession, error) {
	sessions, err := service.queries.ListAcademicSessions(ctx, tenantId)

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return sessions, nil
}

func (service *Service) GetCurrentAcademicSession(ctx context.Context, tenantId uuid.UUID) (*store.AcademicSession, error) {
	currentSession, err := service.queries.GetCurrentAcademicSession(ctx, tenantId)
	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return &currentSession, nil
}

func (service *Service) GetAcademicSessionById(ctx context.Context, tenantId, academicSessionId uuid.UUID) (*store.AcademicSession, error) {
	academicSession, err := service.queries.GetAcademicSessionById(ctx, store.GetAcademicSessionByIdParams{
		ID:       academicSessionId,
		TenantID: tenantId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			message := fmt.Sprintf("Academic session with id %s. Not found", academicSessionId.String())
			return nil, apierror.NotFound(message)
		}
		return nil, academic.TranslateAcademicError(err)
	}

	return &academicSession, nil
}
