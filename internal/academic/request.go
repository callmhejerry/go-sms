package academic

import (
	"github.com/google/uuid"
)

type CreateAcademicSessionRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required,date"`
	EndDate   string `json:"end_date" validate:"required,date"`
	IsCurrent bool   `json:"is_current" validate:"required"`
}

type CreateClassRequest struct {
	Name       string `json:"name" validate:"required"`
	LevelOrder int    `json:"level_order" validate:"required"`
}

type CreateClassArmRequest struct {
	ClassID uuid.UUID `json:"class_id" validate:"required,uuid"`
	Name    string    `json:"name" validate:"required"`
}

type CreateStudentRequest struct {
	AdmissionNumber  string          `json:"admission_number" validate:"required"`
	FirstName        string          `json:"first_name" validate:"required,min=3"`
	LastName         string          `json:"last_name" validate:"required,min=3"`
	MiddleName       *string         `json:"middle_name" validate:"omitempty,min=3"`
	Gender           string          `json:"gender" validate:"oneof=male female"`
	DateOfBirth      string          `json:"date_of_birth" validate:"required,date"`
	CurrentClassArm  *uuid.UUID      `json:"current_class_arm_id" validate:"omitempty,uuid"`
	AdmissionSession *uuid.UUID      `json:"admission_session_id" validate:"omitempty,uuid"`
	Parents          []ParentRequest `json:"parents"`
}

type ParentRequest struct {
	FirstName    string `json:"first_name" validate:"required,min=3"`
	LastName     string `json:"last_name" validate:"required,min=3"`
	Email        string `json:"email" validate:"required,email"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Relationship string `json:"relationship" validate:"required"`
	IsPrimary    bool   `json:"is_primary" validate:"required"`
}
