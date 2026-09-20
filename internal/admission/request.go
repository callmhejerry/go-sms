package admission

import (
	"time"

	"github.com/google/uuid"
)

type CreateAdmissionRequest struct {
	AcademicSessionID uuid.UUID `json:"academic_session_id" validate:"required,uuid"`
	FirstName         string    `json:"first_name" validate:"required"`
	LastName          string    `json:"last_name" validate:"required"`
	MiddleName        *string   `json:"middle_name"`
	Gender            string    `json:"gender" validate:"required"`
	DateOfBirth       time.Time `json:"date_of_birth" validate:"required"`
	PreferredClassID  uuid.UUID `json:"preferred_class_id" validate:"required,uuid"`

	ParentFirstName   string `json:"parent_first_name" validate:"required"`
	ParentLastName    string `json:"parent_last_name" validate:"required"`
	ParentPhoneNumber string `json:"parent_phone_number" validate:"required"`
	ParentEmail       string `json:"parent_email" validate:"required,email"`
	Relationship      string `json:"relationship" validate:"required"`
}
