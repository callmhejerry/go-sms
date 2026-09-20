package academicsession

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
