package grading

import (
	"context"
	"strings"
	"time"

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

func (service *GradingService) ComputeResults(
	ctx context.Context, tenantId uuid.UUID,
	request ComputeResultsRequest,
) (int, *apierror.AppError) {
	scores, err := service.gradingRepository.GetScoresForComputation(
		ctx, tenantId, request.ClassArmID, request.AcademicSessionID,
	)

	if err != nil {
		return 0, err
	}

	type key struct {
		StudentID  uuid.UUID
		SubjectID  uuid.UUID
		ClassArmID uuid.UUID
	}

	type agg struct {
		TotalScore decimal.Decimal
		MaxTotal   decimal.Decimal
	}

	grouped := make(map[key]*agg)

	for _, score := range scores {
		k := key{
			StudentID:  score.StudentID,
			SubjectID:  score.SubjectID,
			ClassArmID: score.ClassArmID,
		}

		if _, exists := grouped[k]; !exists {
			grouped[k] = &agg{
				TotalScore: decimal.Zero,
				MaxTotal:   decimal.Zero,
			}
		}

		grouped[k].TotalScore = grouped[k].TotalScore.Add(score.Score)
		grouped[k].MaxTotal = grouped[k].MaxTotal.Add(score.MaxScore)
	}

	count := 0

	for k, a := range grouped {
		var percentage decimal.Decimal
		if a.MaxTotal.GreaterThan(decimal.Zero) {
			percentage = (a.TotalScore.Div(a.MaxTotal)).Mul(decimal.NewFromInt(100)).Round(2)
		}

		grade, remark := calculateGrade(percentage)

		_, err := service.gradingRepository.UpsertResult(
			ctx, tenantId, CreateResultRequest{
				StudentID:         k.StudentID,
				SubjectID:         k.SubjectID,
				ClassArmID:        k.ClassArmID,
				AcademicSessionID: request.AcademicSessionID,
				TotalScore:        a.TotalScore,
				MaxTotal:          a.MaxTotal,
				Percentage:        percentage,
				Grade:             &grade,
				Remark:            &remark,
			},
		)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func (service *GradingService) GetStudentResults(
	ctx context.Context,
	tenantId uuid.UUID,
	request GetStudentResultRequest,
) ([]store.GetResultsByStudentRow, *apierror.AppError) {
	return service.gradingRepository.GetStudentResults(ctx, tenantId, request)
}

func (service *GradingService) GetStudentReportCard(
	ctx context.Context,
	tenantId, studentId, academicSessionId uuid.UUID,
) (*ReportCardResponse, *apierror.AppError) {
	rows, err := service.gradingRepository.GetStudentReportCard(ctx, tenantId, studentId, academicSessionId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, apierror.NotFound("no results found for this student in the selected session")
	}

	// Assemble a clean report card response
	first := rows[0]

	subjects := make([]SubjectResultResponse, 0, len(rows))
	var totalPercentage decimal.Decimal

	for _, row := range rows {
		subject := SubjectResultResponse{
			SubjectName: row.SubjectName,
			SubjectCode: row.SubjectCode,
			TotalScore:  row.TotalScore,
			MaxTotal:    row.MaxTotal,
			Percentage:  row.Percentage,
			Grade:       row.Grade,
			Remark:      row.Remark,
		}
		subjects = append(subjects, subject)
		totalPercentage = totalPercentage.Add(row.Percentage)
	}

	average := decimal.Zero
	if len(rows) > 0 {
		average = totalPercentage.Div(decimal.NewFromInt(int64(len(rows)))).Round(2)
	}

	studentInfo := ReportCardStudentInfoResponse{
		AdmissionNumber: first.AdmissionNumber,
		FirstName:       first.FirstName,
		LastName:        first.LastName,
		ClassName:       first.ClassName,
		ClassArmName:    first.ClassArmName,
	}
	reportCard := ReportCardResponse{
		StudentInfo:     studentInfo,
		AcademicSession: first.SessionName,
		Subjects:        subjects,
		Average:         average,
		GeneratedAt:     time.Now(),
	}
	return &reportCard, nil
}

func calculateGrade(percentage decimal.Decimal) (string, string) {
	p, _ := percentage.Float64()

	switch {
	case p >= 70:
		return "A", "Excellent"
	case p >= 60:
		return "B", "Very Good"
	case p >= 50:
		return "C", "Good"
	case p >= 40:
		return "D", "Fair"
	default:
		return "F", "Fail"
	}
}
