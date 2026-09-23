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

	GetAssessmentTypeById(
		ctx context.Context,
		tenantId uuid.UUID,
		assessmentId uuid.UUID,
	) (*store.AssessmentType, *apierror.AppError)

	ListAssessmentTypes(ctx context.Context, tenantID uuid.UUID) ([]store.AssessmentType, *apierror.AppError)

	RecordScore(
		ctx context.Context,
		tenantId,
		studentId,
		subjectId,
		academicSessionId,
		classArmId,
		assessmentTypeId uuid.UUID,
		score float64,
		recordedBy *uuid.UUID,
	) (*store.Score, *apierror.AppError)

	GetStudentScores(
		ctx context.Context,
		tenantId, studentId, academicSessionId uuid.UUID,
	) ([]store.GetScoresByStudentRow, *apierror.AppError)

	GetScoresByClassArmAndSubject(
		ctx context.Context,
		tenantId,
		classArmId,
		subjectId,
		assessmentTypeId,
		academicSessionId uuid.UUID,
	) ([]store.GetScoresByClassArmAndSubjectRow, *apierror.AppError)

	GetScoresForComputation(
		ctx context.Context,
		tenantId,
		classArmId,
		academicSessionId uuid.UUID,
	) ([]store.GetScoresForComputationRow, *apierror.AppError)

	UpsertResult(
		ctx context.Context,
		tenantId uuid.UUID,
		request CreateResultRequest,
	) (*store.Result, *apierror.AppError)

	GetStudentResults(
		ctx context.Context,
		tenantId uuid.UUID,
		request GetStudentResultRequest,
	) ([]store.GetResultsByStudentRow, *apierror.AppError)
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

func (repo *gradingRepositoryImpl) RecordScore(
	ctx context.Context,
	tenantId,
	studentId,
	subjectId,
	academicSessionId,
	classArmId,
	assessmentTypeId uuid.UUID,
	score float64,
	recordedBy *uuid.UUID,
) (*store.Score, *apierror.AppError) {
	scoreDecimal := decimal.NewFromFloat(score)

	row, err := repo.queries.UpsertScore(
		ctx, store.UpsertScoreParams{
			TenantID:          tenantId,
			StudentID:         studentId,
			SubjectID:         subjectId,
			AssessmentTypeID:  assessmentTypeId,
			AcademicSessionID: academicSessionId,
			Score:             scoreDecimal,
			ClassArmID:        classArmId,
			RecordedBy:        recordedBy,
		},
	)
	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return &row, nil
}

func (repo *gradingRepositoryImpl) GetStudentScores(
	ctx context.Context,
	tenantId, studentId, academicSessionId uuid.UUID,
) ([]store.GetScoresByStudentRow, *apierror.AppError) {
	scores, err := repo.queries.GetScoresByStudent(
		ctx, store.GetScoresByStudentParams{
			TenantID:          tenantId,
			StudentID:         studentId,
			AcademicSessionID: academicSessionId,
		},
	)
	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return scores, nil
}

func (repo *gradingRepositoryImpl) GetScoresByClassArmAndSubject(
	ctx context.Context,
	tenantId,
	classArmId,
	subjectId,
	assessmentTypeId,
	academicSessionId uuid.UUID,
) ([]store.GetScoresByClassArmAndSubjectRow, *apierror.AppError) {
	scores, err := repo.queries.GetScoresByClassArmAndSubject(
		ctx,
		store.GetScoresByClassArmAndSubjectParams{
			TenantID:          tenantId,
			ClassArmID:        classArmId,
			SubjectID:         subjectId,
			AssessmentTypeID:  assessmentTypeId,
			AcademicSessionID: academicSessionId,
		},
	)
	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return scores, nil
}

func (repo *gradingRepositoryImpl) GetAssessmentTypeById(
	ctx context.Context,
	tenantId,
	assessmentTypeId uuid.UUID,
) (*store.AssessmentType, *apierror.AppError) {
	assessmentType, err := repo.queries.GetAssessmentTypeByID(
		ctx,
		store.GetAssessmentTypeByIDParams{
			ID:       assessmentTypeId,
			TenantID: tenantId,
		},
	)
	if err != nil {
		return nil, translateAssessmentTypeError(err)
	}
	return &assessmentType, nil
}

func (repo *gradingRepositoryImpl) GetScoresForComputation(
	ctx context.Context,
	tenantId,
	classArmId,
	academicSessionId uuid.UUID,
) ([]store.GetScoresForComputationRow, *apierror.AppError) {
	scores, err := repo.queries.GetScoresForComputation(
		ctx, store.GetScoresForComputationParams{
			TenantID:          tenantId,
			ClassArmID:        classArmId,
			AcademicSessionID: academicSessionId,
		},
	)
	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return scores, nil
}

func (repo *gradingRepositoryImpl) UpsertResult(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateResultRequest,
) (*store.Result, *apierror.AppError) {
	result, err := repo.queries.UpsertResult(ctx, store.UpsertResultParams{
		TenantID:          tenantId,
		StudentID:         request.StudentID,
		SubjectID:         request.SubjectID,
		ClassArmID:        request.ClassArmID,
		AcademicSessionID: request.AcademicSessionID,
		TotalScore:        request.TotalScore,
		MaxTotal:          request.TotalScore,
		Percentage:        request.Percentage,
		Grade:             request.Grade,
		Remark:            request.Remark,
	})

	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return &result, nil
}

func (repo *gradingRepositoryImpl) GetStudentResults(
	ctx context.Context,
	tenantId uuid.UUID,
	request GetStudentResultRequest,
) ([]store.GetResultsByStudentRow, *apierror.AppError) {
	rows, err := repo.queries.GetResultsByStudent(ctx, store.GetResultsByStudentParams{
		TenantID:          tenantId,
		StudentID:         request.StudentID,
		AcademicSessionID: request.AcademicSessionID,
	})

	if err != nil {
		return nil, translateScoresTypeError(err)
	}
	return rows, nil
}
