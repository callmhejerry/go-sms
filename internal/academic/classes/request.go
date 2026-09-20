package classes

import "github.com/google/uuid"

type CreateClassRequest struct {
	Name       string `json:"name" validate:"required"`
	LevelOrder int    `json:"level_order" validate:"required"`
}

type CreateClassArmRequest struct {
	ClassID uuid.UUID `json:"class_id" validate:"required,uuid"`
	Name    string    `json:"name" validate:"required"`
}
