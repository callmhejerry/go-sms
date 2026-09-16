package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var input CreateTenantInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenant, err := handler.service.CreateTenant(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(tenant)
}

func (handler *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := handler.service.ListTenants(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}

func (handler *Handler) GetTenant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Tenant ID", http.StatusBadRequest)
		return
	}

	tenant, err := handler.service.GetTenantById(r.Context(), id)

	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}
