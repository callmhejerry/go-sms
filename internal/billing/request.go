package billing

import "github.com/google/uuid"

type CreateFeeTypeRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	IsOptional  bool    `json:"is_optional" validate:"required"`
}

type CreatFeeStructureRequest struct {
	AcademicSessionId uuid.UUID  `json:"academic_session_id" validate:"required,uuid"`
	FeeTypeId         uuid.UUID  `json:"fee_type_id" validate:"required,uuid"`
	ClassId           *uuid.UUID `json:"class_id" validate:"omitempty,uuid"`
	DueDate           *string    `json:"due_date" validate:"omitempty,date"`
	AmountInKobo      int        `json:"amount_in_kobo" validate:"required,min=1"`
}
