package inventory

import (
	"context"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
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
}

type inventoryRepositoryImpl struct {
	queries *store.Queries
}

func NewRepositoryImpl(queries *store.Queries) InventoryRepository {
	return &inventoryRepositoryImpl{
		queries: queries,
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
