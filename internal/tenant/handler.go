package tenant

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (handler *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var input CreateTenantRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), handler.logger)
		return
	}

	tenant, err := handler.service.CreateTenant(r.Context(), input)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateTenantResponse{
		TenantID:  tenant.ID.String(),
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		CreatedAt: tenant.CreatedAt.Time,
	})
}

func (handler *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := handler.service.ListTenants(r.Context())
	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}

func (handler *Handler) GetTenant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid tenant id"), handler.logger)
		return
	}

	tenant, err := handler.service.GetTenantById(r.Context(), id)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}
