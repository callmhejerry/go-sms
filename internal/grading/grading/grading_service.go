package grading

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type GradingService struct {
	gradingRepository GradingRepository
}

func NewGradingService(gradingRepository GradingRepository) *GradingService {
	return &GradingService{
		gradingRepository: gradingRepository,
	}
}

func (service *GradingService) CreateAssessmentType(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateAssessmentTypeRequest,
) (*store.AssessmentType, *apierror.AppError) {
	name := strings.TrimSpace(strings.ToLower(request.Name))
	if name == "" {
		return nil, apierror.Validation("name is requried")
	}

	if request.MaxScore <= 0 {
		return nil, apierror.Validation("max_score must be greater than zero")
	}
	if request.Weight != nil && *request.Weight <= 0 {
		return nil, apierror.Validation("weight must be greater than zero")
	}

	return service.gradingRepository.CreateAssessmentType(
		ctx,
		tenantId,
		name, request.MaxScore,
		request.Weight, request.IsExam,
	)
}

func (service *GradingService) ListAssessmentTypes(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.AssessmentType, *apierror.AppError) {
	return service.gradingRepository.ListAssessmentTypes(ctx, tenantId)
}
