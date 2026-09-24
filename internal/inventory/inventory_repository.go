package inventory

import (
	"context"
	"errors"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/constants"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	CreateInventoryCategory(
		ctx context.Context,
		tenantId uuid.UUID,
		name string,
	) (*store.InventoryCategory, *apierror.AppError)

	ListInventoryCategories(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.InventoryCategory, *apierror.AppError)

	CreateInventoryItem(
		ctx context.Context,
		tenantId uuid.UUID,
		request CreateInventoryItemRequest,
	) (*store.InventoryItem, *apierror.AppError)

	ListInventoryItems(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.ListInventoryItemsRow, *apierror.AppError)

	RecordMovement(
		ctx context.Context,
		tenantId uuid.UUID,
		performedBy uuid.UUID,
		request RecordStockMovementRequest,
	) (*store.StockMovement, *apierror.AppError)

	ListItemMovement(
		ctx context.Context,
		tenantId uuid.UUID,
		itemId uuid.UUID,
	) ([]store.StockMovement, *apierror.AppError)

	CreateIssuance(
		ctx context.Context,
		tenantId, issuedById uuid.UUID,
		request CreateIssuanceRequest,
	) (*store.InventoryIssuance, *apierror.AppError)

	ListInventoryIssuances(
		ctx context.Context,
		tenantId, inventoryItemId uuid.UUID,
	) ([]store.InventoryIssuance, *apierror.AppError)

	ListLowStockItems(
		ctx context.Context, tenantID uuid.UUID,
	) ([]store.ListLowStockItemsRow, error)
}

type inventoryRepositoryImpl struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func NewRepositoryImpl(queries *store.Queries, pool *pgxpool.Pool) InventoryRepository {
	return &inventoryRepositoryImpl{
		queries: queries,
		pool:    pool,
	}
}

func (repo *inventoryRepositoryImpl) CreateInventoryCategory(
	ctx context.Context,
	tenantId uuid.UUID,
	name string,
) (*store.InventoryCategory, *apierror.AppError) {
	inventory, err := repo.queries.CreateInventoryCategory(
		ctx, store.CreateInventoryCategoryParams{
			TenantID: tenantId,
			Name:     name,
		},
	)

	if err != nil {
		return nil, translateInventoryCategoryError(err)
	}

	return &inventory, nil
}

func (repo *inventoryRepositoryImpl) ListInventoryCategories(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.InventoryCategory, *apierror.AppError) {
	rows, err := repo.queries.ListInventoryCategories(ctx, tenantId)
	if err != nil {
		return nil, translateInventoryCategoryError(err)
	}
	return rows, nil
}

func (repo *inventoryRepositoryImpl) CreateInventoryItem(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateInventoryItemRequest,
) (*store.InventoryItem, *apierror.AppError) {
	inventoryItem, err := repo.queries.CreateInventoryItem(ctx, store.CreateInventoryItemParams{
		TenantID:        tenantId,
		CategoryID:      request.CategoryID,
		Name:            request.Name,
		Code:            request.Code,
		Description:     request.Description,
		Unit:            request.Unit,
		QuantityInStock: request.QuantityInStock,
		ReorderLevel:    request.ReorderLevel,
		UnitCostKobo:    request.UnitCostKobo,
	})

	if err != nil {
		return nil, translateInventoryItemError(err)
	}
	return &inventoryItem, nil
}

func (repo *inventoryRepositoryImpl) ListInventoryItems(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.ListInventoryItemsRow, *apierror.AppError) {
	rows, err := repo.queries.ListInventoryItems(ctx, tenantId)
	if err != nil {
		return nil, translateInventoryItemError(err)
	}
	return rows, nil
}

func (repo *inventoryRepositoryImpl) RecordMovement(
	ctx context.Context,
	tenantId uuid.UUID,
	performedBy uuid.UUID,
	request RecordStockMovementRequest,
) (*store.StockMovement, *apierror.AppError) {

	var stockMovement store.StockMovement

	err := database.WithTx(ctx, repo.pool, func(queries *store.Queries) error {
		inventory, err := queries.GetInventoryItemForUpdate(ctx, store.GetInventoryItemForUpdateParams{
			ID:       request.ItemID,
			TenantID: tenantId,
		})

		if err != nil {
			return translateInventoryItemError(err)
		}

		var newStockQuantity int32

		switch request.MovementType {
		case string(constants.In):
			newStockQuantity = inventory.QuantityInStock + request.Quantity
		case string(constants.Out):
			if inventory.QuantityInStock < request.Quantity {
				return ErrInventoryItemQuantityGreaterThanZero
			}
			newStockQuantity = inventory.QuantityInStock - request.Quantity
		case string(constants.Adjust):
			newStockQuantity = request.Quantity
		default:
			return ErrInvalidStockMovementType
		}

		_, err = queries.UpdateItemStock(ctx, store.UpdateItemStockParams{
			ID:              request.ItemID,
			TenantID:        tenantId,
			QuantityInStock: newStockQuantity,
		})
		if err != nil {
			return translateInventoryItemError(err)
		}

		movement, err := queries.CreateStockMovement(ctx, store.CreateStockMovementParams{
			TenantID:     tenantId,
			ItemID:       request.ItemID,
			MovementType: request.MovementType,
			Quantity:     request.Quantity,
			Reason:       request.Reason,
			Reference:    request.Reference,
			PerformedBy:  &performedBy,
			Notes:        request.Notes,
		})

		if err != nil {
			return translateStockMovementError(err)
		}

		stockMovement = movement

		return nil
	})

	if err != nil {
		var appErr apierror.AppError
		if errors.As(err, &appErr) {
			return nil, &appErr
		} else {
			return nil, apierror.Internal(err, "something went wrong")
		}
	}

	return &stockMovement, nil
}

func (repo *inventoryRepositoryImpl) ListItemMovement(
	ctx context.Context,
	tenantId uuid.UUID,
	itemId uuid.UUID,
) ([]store.StockMovement, *apierror.AppError) {
	rows, err := repo.queries.ListStockMovementsByItem(ctx, store.ListStockMovementsByItemParams{
		TenantID: tenantId,
		ItemID:   itemId,
	})

	if err != nil {
		return nil, translateStockMovementError(err)
	}
	return rows, nil
}

func (repo *inventoryRepositoryImpl) CreateIssuance(
	ctx context.Context,
	tenantId, issuedById uuid.UUID,
	request CreateIssuanceRequest,
) (*store.InventoryIssuance, *apierror.AppError) {
	var inventoryIssuance store.InventoryIssuance

	err := database.WithTx(ctx, repo.pool, func(queries *store.Queries) error {
		inventory, err := queries.GetInventoryItemForUpdate(ctx, store.GetInventoryItemForUpdateParams{
			ID:       request.ItemID,
			TenantID: tenantId,
		})
		if err != nil {
			return translateInventoryItemError(err)
		}

		if inventory.QuantityInStock < request.Quantity {
			return apierror.Validation("insufficient stock")
		}

		newStockQuantity := inventory.QuantityInStock - request.Quantity

		_, err = queries.UpdateItemStock(ctx, store.UpdateItemStockParams{
			ID:              request.ItemID,
			TenantID:        tenantId,
			QuantityInStock: newStockQuantity,
		})
		if err != nil {
			return translateInventoryItemError(err)
		}

		reason := "issuance"
		_, err = queries.CreateStockMovement(ctx, store.CreateStockMovementParams{
			TenantID:     tenantId,
			ItemID:       request.ItemID,
			MovementType: string(constants.Out),
			Quantity:     request.Quantity,
			Reason:       &reason,
			Reference:    nil,
			PerformedBy:  &issuedById,
			Notes:        request.Notes,
		})

		if err != nil {
			return translateStockMovementError(err)
		}
		return nil
	})

	if err != nil {
		var appErr apierror.AppError
		if errors.As(err, &appErr) {
			return nil, translateInventoryIssuanceError(err)
		}
		return nil, apierror.Internal(err, "Something went wrong")
	}
	return &inventoryIssuance, nil
}

func (repo *inventoryRepositoryImpl) ListInventoryIssuances(
	ctx context.Context,
	tenantId, inventoryItemId uuid.UUID,
) ([]store.InventoryIssuance, *apierror.AppError) {
	rows, err := repo.queries.ListIssuancesByItem(
		ctx, store.ListIssuancesByItemParams{
			TenantID: tenantId,
			ItemID:   inventoryItemId,
		},
	)
	if err != nil {
		return nil, translateInventoryIssuanceError(err)
	}
	return rows, nil
}

func (repo *inventoryRepositoryImpl) ListLowStockItems(
	ctx context.Context, tenantID uuid.UUID,
) ([]store.ListLowStockItemsRow, error) {
	items, err := repo.queries.ListLowStockItems(ctx, tenantID)
	if err != nil {
		return nil, translateInventoryItemError(err)
	}
	return items, nil
}
