package academic

import (
	"time"

	"github.com/google/uuid"
)

type CreateAcademicSessionRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `json:"is_current"`
}

type CreateClassRequest struct {
	Name       string `json:"name"`
	LevelOrder int    `json:"level_order"`
}

type CreateClassArmRequest struct {
	ClassID uuid.UUID `json:"class_id"`
	Name    string    `json:"name"`
}

type CreateStudentRequest struct {
	AdmissionNumber  string          `json:"admission_number"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	MiddleName       *string         `json:"middle_name"`
	Gender           string          `json:"gender"`
	DateOfBirth      time.Time       `json:"date_of_birth"`
	CurrentClassArm  *uuid.UUID      `json:"current_class_arm_id"`
	AdmissionSession *uuid.UUID      `json:"admission_session_id"`
	Parents          []ParentRequest `json:"parents"`
}

type ParentRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	Relationship string `json:"relationship"`
	IsPrimary    bool   `json:"is_primary"`
}
