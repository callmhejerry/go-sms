package academic

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
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
		return nil, translateAcademicError(err)
	}
	return &session, nil
}

func (service *Service) ListAcademicSession(ctx context.Context, tenantId uuid.UUID) ([]store.AcademicSession, error) {
	sessions, err := service.queries.ListAcademicSessions(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return sessions, nil
}

func (service *Service) GetCurrentAcademicSession(ctx context.Context, tenantId uuid.UUID) (*store.AcademicSession, error) {
	currentSession, err := service.queries.GetCurrentAcademicSession(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})
	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &currentSession, nil
}

func (service *Service) CreateClass(ctx context.Context, tenantId uuid.UUID, request CreateClassRequest) (*store.Class, error) {
	name := strings.TrimSpace(strings.ToLower(request.Name))

	if name == "" {
		return nil, apierror.Validation("class name is required")
	}

	class, err := service.queries.CreateClass(ctx, store.CreateClassParams{
		TenantID:   pgtype.UUID{Bytes: tenantId, Valid: true},
		Name:       name,
		LevelOrder: int32(request.LevelOrder),
	})
	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &class, nil
}

func (service *Service) ListClasses(ctx context.Context, tenantId uuid.UUID) ([]store.Class, error) {
	classes, err := service.queries.ListClasses(ctx, pgtype.UUID{Bytes: tenantId, Valid: true})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return classes, nil
}

func (service *Service) CreateClassArm(ctx context.Context, tenantId uuid.UUID, request CreateClassArmRequest) (*store.ClassArm, error) {
	name := strings.TrimSpace(strings.ToLower(request.Name))

	if name == "" {
		return nil, apierror.Validation("Class arm name is required")
	}

	class, err := service.queries.GetClassByID(ctx, store.GetClassByIDParams{
		ID:       pgtype.UUID{Bytes: request.ClassID, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}

	classArm, err := service.queries.CreateClassArm(ctx, store.CreateClassArmParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  class.ID,
		Name:     name,
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &classArm, nil
}

func (service *Service) ListClassArms(ctx context.Context, tenantId uuid.UUID, classId uuid.UUID) ([]store.ClassArm, error) {
	classArms, err := service.queries.ListClassArms(ctx, store.ListClassArmsParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  pgtype.UUID{Bytes: classId, Valid: true},
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return classArms, nil
}
