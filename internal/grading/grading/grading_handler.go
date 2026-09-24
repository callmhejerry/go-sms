package grading

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
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

func (h *GradingHandler) ComputeResults(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input ComputeResultsRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	count, err := h.service.ComputeResults(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	response := ComputeResultResponse{
		Message:      "results computed successfully",
		ResultsCount: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GradingHandler) GetStudentResults(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("student_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	sessionIDStr := r.URL.Query().Get("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		apierror.WriteError(w, apierror.Validation("session_id query parameter is required"), h.logger)
		return
	}

	results, err := h.service.GetStudentResults(r.Context(), claims.TenantID, GetStudentResultRequest{
		StudentID:         studentID,
		AcademicSessionID: sessionID,
	})
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *GradingHandler) GetStudentReportCard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("student_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	sessionID, err := uuid.Parse(r.URL.Query().Get("session_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("session_id query parameter is required"), h.logger)
		return
	}

	response, err := h.service.GetStudentReportCard(r.Context(), claims.TenantID, studentID, sessionID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
