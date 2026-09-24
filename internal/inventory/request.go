package inventory

import "github.com/google/uuid"

type CreateInventoryCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateInventoryItemRequest struct {
	CategoryID      *uuid.UUID `json:"category_id"`
	Name            string     `json:"name" validate:"required"`
	Code            *string    `json:"code"`
	Description     *string    `json:"description"`
	Unit            string     `json:"unit"`
	QuantityInStock int32      `json:"quantity_in_stock"`
	ReorderLevel    int32      `json:"reorder_level"`
	UnitCostKobo    *int64     `json:"unit_cost_kobo"`
}

type RecordStockMovementRequest struct {
	ItemID       uuid.UUID `json:"item_id" validate:"required"`
	MovementType string    `json:"movement_type" validate:"required,oneof=in out adjust"`
	Quantity     int32     `json:"quantity" validate:"required,gt=0"`
	Reason       *string   `json:"reason"`
	Reference    *string   `json:"reference"`
	Notes        *string   `json:"notes"`
}

type CreateIssuanceRequest struct {
	ItemID            uuid.UUID  `json:"item_id" validate:"required"`
	Quantity          int32      `json:"quantity" validate:"required,gt=0"`
	IssuedToType      string     `json:"issued_to_type" validate:"required,oneof=student department staff"`
	IssuedToID        uuid.UUID  `json:"issued_to_id" validate:"required"`
	AcademicSessionID *uuid.UUID `json:"academic_session_id"`
	Notes             *string    `json:"notes"`
}
