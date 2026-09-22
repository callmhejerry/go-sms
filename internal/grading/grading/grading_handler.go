package grading

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
)

type GradingHandler struct {
	logger       *slog.Logger
	service      *GradingService
	appValidator *validation.AppValidator
}

func NewGradingHandler(service *GradingService, logger *slog.Logger, appValidator *validation.AppValidator) *GradingHandler {
	return &GradingHandler{
		logger:       logger,
		service:      service,
		appValidator: appValidator,
	}
}

func (h *GradingHandler) CreateAssessmentType(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input CreateAssessmentTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	at, err := h.service.CreateAssessmentType(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(at)
}

func (h *GradingHandler) ListAssessmentTypes(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	types, err := h.service.ListAssessmentTypes(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}
