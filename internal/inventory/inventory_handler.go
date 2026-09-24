package inventory

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service      *InventoryService
	logger       *slog.Logger
	appValidator *validation.AppValidator
}

func NewHandler(service *InventoryService, logger *slog.Logger, appValidator *validation.AppValidator) *InventoryHandler {
	return &InventoryHandler{
		service:      service,
		logger:       logger,
		appValidator: appValidator,
	}
}

func (h *InventoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input CreateInventoryCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}
	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	cat, err := h.service.CreateInventoryCategory(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (h *InventoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	cats, err := h.service.ListInventoryCategories(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func (h *InventoryHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input CreateInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	item, err := h.service.CreateInventoryItem(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	items, err := h.service.ListInventoryItems(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) RecordStockMovement(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input RecordStockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	movement, err := h.service.RecordStockMovement(r.Context(), claims.TenantID, claims.UserID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movement)
}

func (h *InventoryHandler) ListItemMovements(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("item_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid item id"), h.logger)
		return
	}

	movements, err := h.service.ListItemMovements(r.Context(), claims.TenantID, itemID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movements)
}

func (h *InventoryHandler) CreateIssuance(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input CreateIssuanceRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	issuance, err := h.service.CreateInventoryIssuance(r.Context(), claims.TenantID, claims.UserID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issuance)
}

func (h *InventoryHandler) ListItemIssuances(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("item_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid item id"), h.logger)
		return
	}

	issuances, err := h.service.ListInventoryIssuance(r.Context(), claims.TenantID, itemID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issuances)
}
