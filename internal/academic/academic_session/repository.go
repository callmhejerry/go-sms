package academicsession

import (
	"context"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type AcademicRepository interface {
	CreateAcademicSession(
		ctx context.Context,
		tenantId uuid.UUID,
		name string,
		startDate time.Time,
		endDate time.Time,
		isCurrent bool,
	) (*store.AcademicSession, apierror.AppError)
}

type AcademicRepositoryImpl struct {
	queries *store.Queries
}

func NewRepository(queries *store.Queries) *AcademicRepositoryImpl {
	return &AcademicRepositoryImpl{
		queries: queries,
	}
}

func (repo *AcademicRepositoryImpl) CreateAcademicSession(
	ctx context.Context,
	tenantId uuid.UUID,
	name string,
	startDate time.Time,
	endDate time.Time,
	isCurrent bool,
) (*store.AcademicSession, *apierror.AppError) {

	newAcademicSession, err := repo.queries.CreateAcademicSession(ctx, store.CreateAcademicSessionParams{
		TenantID:  tenantId,
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		IsCurrent: isCurrent,
	})

	if err != nil {
		return nil, translateError(err)
	}

	return &newAcademicSession, nil
}
