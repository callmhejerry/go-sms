package grading

import (
	"context"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GradingRepository interface {
	CreateAssessmentType(
		ctx context.Context,
		tenantId uuid.UUID,
		name string,
		maxscore float64,
		weight *float64,
		isExam bool,
	) (*store.AssessmentType, *apierror.AppError)

	ListAssessmentTypes(ctx context.Context, tenantID uuid.UUID) ([]store.AssessmentType, *apierror.AppError)
}

type gradingRepositoryImpl struct {
	queries *store.Queries
}

func NewGradingRespositoryImpl(queries *store.Queries) GradingRepository {
	return &gradingRepositoryImpl{
		queries: queries,
	}
}

func (repo *gradingRepositoryImpl) CreateAssessmentType(
	ctx context.Context,
	tenantId uuid.UUID,
	name string,
	maxscore float64,
	weight *float64,
	isExam bool,
) (*store.AssessmentType, *apierror.AppError) {
	maxscoreDecimal := decimal.NewFromFloat(maxscore)
	var weightDecimal decimal.Decimal

	if weight != nil {
		weightDecimal = decimal.NewFromFloat(*weight)
	}

	assessmentType, err := repo.queries.CreateAssessmentType(
		ctx, store.CreateAssessmentTypeParams{
			TenantID: tenantId,
			Name:     name,
			MaxScore: maxscoreDecimal,
			IsExam:   isExam,
			Weight:   &weightDecimal,
		},
	)
	if err != nil {
		return nil, translateAssessmentTypeError(err)
	}
	return &assessmentType, nil
}

func (s *gradingRepositoryImpl) ListAssessmentTypes(ctx context.Context, tenantID uuid.UUID) ([]store.AssessmentType, *apierror.AppError) {
	types, err := s.queries.ListAssessmentTypes(ctx, tenantID)
	if err != nil {
		return nil, translateAssessmentTypeError(err)
	}
	return types, nil
}
