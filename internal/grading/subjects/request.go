package subjects

import "github.com/google/uuid"

type CreateSubjectRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type AssignTeacherRequest struct {
	UserID            uuid.UUID  `json:"user_id" validate:"required,uuid"`
	SubjectID         uuid.UUID  `json:"subject_id" validate:"required,uuid"`
	AcademicSessionID uuid.UUID  `json:"academic_session_id" validate:"required,uuid"`
	ClassID           uuid.UUID  `json:"class_id" validate:"required,uuid"`
	ClassArmID        *uuid.UUID `json:"class_arm_id,omitempty"` // optional
}
