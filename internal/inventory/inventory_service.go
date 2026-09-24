package inventory

import (
	"context"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type InventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (service *InventoryService) CreateInventoryCategory(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateInventoryCategoryRequest,
) (*store.InventoryCategory, *apierror.AppError) {
	return service.repo.CreateInventoryCategory(ctx, tenantId, request.Name)
}

func (service *InventoryService) CreateInventoryItem(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateInventoryItemRequest,
) (*store.InventoryItem, *apierror.AppError) {
	return service.repo.CreateInventoryItem(ctx, tenantId, request)
}

func (service *InventoryService) ListInventoryCategories(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.InventoryCategory, *apierror.AppError) {
	return service.repo.ListInventoryCategories(ctx, tenantId)
}

func (service *InventoryService) ListInventoryItems(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.ListInventoryItemsRow, *apierror.AppError) {
	return service.repo.ListInventoryItems(ctx, tenantId)
}

func (service *InventoryService) RecordStockMovement(
	ctx context.Context,
	tenantId, recordedBy uuid.UUID,
	request RecordStockMovementRequest,
) (*store.StockMovement, *apierror.AppError) {
	return service.repo.RecordMovement(ctx, tenantId, recordedBy, request)
}

func (service *InventoryService) ListItemMovements(
	ctx context.Context,
	tenantId, itemId uuid.UUID,
) ([]store.StockMovement, *apierror.AppError) {
	return service.repo.ListItemMovement(ctx, tenantId, itemId)
}

func (service *InventoryService) CreateInventoryIssuance(
	ctx context.Context,
	tenantId, issuedBy uuid.UUID,
	request CreateIssuanceRequest,
) (*store.InventoryIssuance, *apierror.AppError) {
	return service.repo.CreateIssuance(ctx, tenantId, issuedBy, request)
}

func (service *InventoryService) ListInventoryIssuance(
	ctx context.Context,
	tenantId, inventoryItemId uuid.UUID,
) ([]store.InventoryIssuance, *apierror.AppError) {
	return service.repo.ListInventoryIssuances(ctx, tenantId, inventoryItemId)
}
