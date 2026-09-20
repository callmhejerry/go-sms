package classes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/callmhejerry/sms/internal/academic"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func NewService(queries *store.Queries, pool *pgxpool.Pool) *Service {
	return &Service{
		queries: queries,
		pool:    pool,
	}
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
		return nil, academic.TranslateAcademicError(err)
	}
	return &class, nil
}

func (service *Service) GetClassById(ctx context.Context, tenantId, classId uuid.UUID) (*store.Class, error) {
	class, err := service.queries.GetClassByID(ctx, store.GetClassByIDParams{
		ID:       pgtype.UUID{Bytes: classId, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			message := fmt.Sprintf("class with id %s. Not found", classId.String())
			return nil, apierror.NotFound(message)
		}
		return nil, academic.TranslateAcademicError(err)
	}
	return &class, nil
}

func (service *Service) ListClasses(ctx context.Context, tenantId uuid.UUID) ([]store.Class, error) {
	classes, err := service.queries.ListClasses(ctx, pgtype.UUID{Bytes: tenantId, Valid: true})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
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
		return nil, academic.TranslateAcademicError(err)
	}

	classArm, err := service.queries.CreateClassArm(ctx, store.CreateClassArmParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  class.ID,
		Name:     name,
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return &classArm, nil
}

func (service *Service) ListClassArms(ctx context.Context, tenantId uuid.UUID, classId uuid.UUID) ([]store.ClassArm, error) {
	classArms, err := service.queries.ListClassArms(ctx, store.ListClassArmsParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  pgtype.UUID{Bytes: classId, Valid: true},
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return classArms, nil
}
