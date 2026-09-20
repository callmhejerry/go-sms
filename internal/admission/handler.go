package admission

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
	service   *Service
	logger    *slog.Logger
	validator *validation.Validator
}

func NewHandler(service *Service, logger *slog.Logger, validator *validation.Validator) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}

func (handler *Handler) CreateAdmission(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request CreateAdmissionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	admission, err := handler.service.CreateAdmission(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admission)
}

func (handler *Handler) ListAdmissions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}
	admissions, err := handler.service.ListAdmissions(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admissions)
}

func (handler *Handler) AcceptAdmission(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	admissionId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		apierror.Validation("Invalid admission id")
		return
	}

	var body struct {
		ClassArmID *uuid.UUID `json:"class_arm_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	admission, student, err := handler.service.AcceptAdmission(
		r.Context(),
		claims.TenantID,
		admissionId,
		claims.UserID,
		body.ClassArmID,
	)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	response := AcceptAdmissionResponse{
		Admission: *admission,
		Student:   *student,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
