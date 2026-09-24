package inventory

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInventoryCategoryNotFound = apierror.NotFound("Inventory category not found")

	ErrInventoryCategoryNameAlreadyExists = apierror.New("inventory_category_name_already_exist", "Inventory category with the same name already exist", http.StatusConflict, nil, nil)

	ErrInventoryItemNotFound = apierror.NotFound("Inventory item not found")

	ErrInventoryItemNameAlreadyExists = apierror.New("inventory_item_name_already_exist", "Inventory item with the same name already exist", http.StatusConflict, nil, nil)

	ErrInventoryItemCodeAlreadyExists = apierror.New("inventory_item_code_already_exist", "Inventory item with the same code already exist", http.StatusConflict, nil, nil)

	ErrInventoryItemQuantityGreaterThanZero = apierror.New("inventory_item_quantity_must_be_greater_than_zero", "Inventory item quantity must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrInventoryReorderLevelMustBeGreaterThanZero = apierror.New("inventory_item_reorder_level_must_be_greater_than_zero", "Inventory reorder level must be greater than zero", http.StatusBadRequest, nil, nil)
)

func translateInventoryCategoryError(err error) *apierror.AppError {
	if err != nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInventoryCategoryNotFound
	}

	var pgErr pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			switch pgErr.ConstraintName {
			case "inventory_category_unique":
				return ErrInventoryCategoryNameAlreadyExists
			}
		}
	}
	return apierror.Internal(err, "Something went wrong")
}

func translateInventoryItemError(err error) *apierror.AppError {
	if err != nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInventoryItemNotFound
	}

	var pgErr pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			switch pgErr.ConstraintName {
			case "inventory_items_code_unique":
				return ErrInventoryItemCodeAlreadyExists
			case "inventory_items_name_unique":
				return ErrInventoryItemNameAlreadyExists
			}
		case apierror.CheckViolation:
			switch pgErr.ColumnName {
			case "quantity_in_stock":
				return ErrInventoryItemQuantityGreaterThanZero
			case "reorder_level":
				return ErrInventoryReorderLevelMustBeGreaterThanZero
			}
		case apierror.ForeignKeyViolation:
			if pgErr.ColumnName == "category_id" {
				return ErrInventoryCategoryNotFound
			}
		}
	}
	return apierror.Internal(err, "Something went wrong")
}
