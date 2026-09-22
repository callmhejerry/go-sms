package grading

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func (service *GradingService) RecordStudentScore(
	ctx context.Context,
	tenantId uuid.UUID,
	recordedBy *uuid.UUID,
	request RecordScoreRequest,
) (*store.Score, *apierror.AppError) {
	score := decimal.NewFromFloat(request.Score)

	assessment, err := service.gradingRepository.GetAssessmentTypeById(
		ctx, tenantId,
		request.AssessmentTypeID,
	)

	if err != nil {
		return nil, err
	}

	maxScore := assessment.MaxScore

	if score.GreaterThan(maxScore) {
		return nil, ErrScoreGreaterThanMaxScore
	}

	return service.gradingRepository.RecordScore(
		ctx, tenantId,
		request.StudentID,
		request.SubjectID,
		request.AcademicSessionID,
		request.ClassArmID,
		request.AssessmentTypeID,
		request.Score,
		recordedBy,
	)
}

func (service *GradingService) GetStudentScores(
	ctx context.Context,
	tenantId,
	studentId,
	academicSessionId uuid.UUID,
) ([]store.GetScoresByStudentRow, *apierror.AppError) {
	return service.GetStudentScores(ctx, tenantId, studentId, academicSessionId)
}
