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
