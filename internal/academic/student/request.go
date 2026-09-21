package student

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

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

type SearchStudentRequest struct {
	Query      *string    `json:"query"` //name or admission number
	ClassArmId *uuid.UUID `json:"class_arm_id" validate:"omitempty,uuid"`
	Status     *string    `json:"status"`
}

func parseSearchStudentQuery(r *http.Request) SearchStudentRequest {
	query := r.URL.Query()

	searchQuery := SearchStudentRequest{
		Query:      nil,
		ClassArmId: nil,
		Status:     nil,
	}

	if value := query.Get("search"); value != "" {
		searchQuery.Query = &value
	}
	if value := query.Get("class_arm_id"); value != "" {
		if parsedValue, err := uuid.Parse(value); err == nil {
			searchQuery.ClassArmId = &parsedValue
		}
	}
	if value := query.Get("status"); value != "" {
		searchQuery.Status = &value
	}

	return searchQuery
}

type UpdateStudentRequest struct {
	FirstName         *string    `json:"first_name" validate:"omitempty,min=3"`
	LastName          *string    `json:"last_name" validate:"omitempty,min=3"`
	MiddleName        *string    `json:"middle_name" validate:"omitempty,min=3"`
	Gender            *string    `json:"gender" validate:"omitempty,oneof=male female"`
	DateOfBirth       *string    `json:"date_of_birth" validate:"omitempty,date_of_birth"`
	CurrentClassArmID *uuid.UUID `json:"current_class_arm_id" validate:"omitempty,uuid"`
	Status            *string    `json:"status"`
}

type ListStudentCursor struct {
	LastName  *string    `json:"last_name"`
	FirstName *string    `json:"first_name"`
	ID        *uuid.UUID `json:"id"`
}

func DecodeListStudentCursor(b64 *string) ListStudentCursor {
	cursor := ListStudentCursor{}

	if b64 == nil || *b64 == "" {
		return cursor
	}

	if decodedBytes, err := base64.RawURLEncoding.DecodeString(*b64); err == nil {
		decodedStr := string(decodedBytes)
		parts := strings.Split(decodedStr, "|")

		if len(parts) != 3 {
			return cursor
		}

		if parsedId, err := uuid.Parse(parts[2]); err == nil {
			cursor.LastName = &parts[0]
			cursor.FirstName = &parts[1]
			cursor.ID = &parsedId
		}
	}

	return cursor
}

func EncodeListStudentCursor(lastname, firstname string, id uuid.UUID) string {
	idString := id.String()
	strToEncode := fmt.Sprintf("%s|%s|%s", lastname, firstname, idString)

	return base64.RawURLEncoding.EncodeToString([]byte(strToEncode))
}
