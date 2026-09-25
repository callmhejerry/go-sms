package classes

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
)

type Handler struct {
	logger       *slog.Logger
	service      *Service
	appValidator *validation.AppValidator
}

func NewHandler(service *Service, logger *slog.Logger, appValidator *validation.AppValidator) *Handler {
	return &Handler{
		logger:       logger,
		service:      service,
		appValidator: appValidator,
	}
}

func (handler *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	var request CreateClassRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	class, err := handler.service.CreateClass(r.Context(), tenantId, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

func (handler *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	classes, err := handler.service.ListClasses(r.Context(), tenantId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

func (handler *Handler) CreateClassArm(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	var request CreateClassArmRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	classArm, err := handler.service.CreateClassArm(r.Context(), tenantId, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(classArm)
}

func (handler *Handler) ListClassArms(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	classId, perr := uuid.Parse(r.PathValue("class_id"))
	if perr != nil {
		apierror.WriteError(w, apierror.Validation("Invalid class id"), handler.logger)
		return
	}
	classArms, err := handler.service.ListClassArms(r.Context(), tenantId, classId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classArms)
}
