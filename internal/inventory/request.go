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
