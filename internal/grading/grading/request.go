package grading

import "github.com/google/uuid"

type CreateAssessmentTypeRequest struct {
	Name     string   `json:"name" validate:"required"`
	MaxScore float64  `json:"max_score" validate:"required,min=1"`
	Weight   *float64 `json:"weight"`
	IsExam   bool     `json:"is_exam"`
}

type RecordScoreRequest struct {
	StudentID         uuid.UUID `json:"student_id" validate:"required,uuid"`
	SubjectID         uuid.UUID `json:"subject_id" validate:"required,uuid"`
	AssessmentTypeID  uuid.UUID `json:"assessment_type_id" validate:"required,uuid"`
	ClassArmID        uuid.UUID `json:"class_arm_id" validate:"required,uuid"`
	AcademicSessionID uuid.UUID `json:"academic_session_id" validate:"required,uuid"`
	Score             float64   `json:"score" validate:"required,min=1"`
}
